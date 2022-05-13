package osf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	defaultBaseURL     = "https://api.osf.io/v2/"
	defaultBaseTestURL = "https://api.test.osf.io/v2/"

	userAgent = "go-osf"
)

type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	BaseURL *url.URL

	UserAgent string

	rateMu sync.Mutex

	common service

	Citations *CitationsService
	Preprints *PreprintsService
	Files     *FilesService
}

type service struct {
	client *Client
}

func (c *Client) Client() *http.Client {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	clientCopy := *c.client
	return &clientCopy
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Citations = (*CitationsService)(&c.common)
	c.Preprints = (*PreprintsService)(&c.common)
	c.Files = (*FilesService)(&c.common)
	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	// u, err := c.BaseURL.Parse(urlStr)
	// if err != nil {
	// 	return nil, err
	// }
	u := c.BaseURL.String() + urlStr

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// TODO: why the /v2 gone
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "*/*")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

type ListOptions struct {
	Page    int `url:"page[number],omitempty"`
	PerPage int `url:"page[size],omitempty"`
}

type PaginationMeta struct {
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
}

type PaginationLinks struct {
	First *string         `json:"first"`
	Last  *string         `json:"last"`
	Prev  *string         `json:"prev"`
	Next  *string         `json:"next"`
	Meta  *PaginationMeta `json:"meta"`
}

type Data[T any] struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Attributes    T                      `json:"attributes"`
	Links         map[string]interface{} `json:"links"`
	Relationships map[string]interface{} `json:"relationships"`
}

type SingleResponse[T any] struct {
	Data Data[T] `json:"data"`
}

func (r *SingleResponse[T]) GetData() T {
	return r.Data.Attributes
}

type ManyResponse[T any] struct {
	Data            []Data[T]        `json:"data"`
	PaginationLinks *PaginationLinks `json:"links"`

	// TODO: Include more user-friendly attributes regarding paginations.
}

func (r *ManyResponse[T]) GetData() []T {
	res := make([]T, 0, len(r.Data))
	for _, obj := range r.Data {
		res = append(res, obj.Attributes)
	}
	return res
}

// do performs logic for doSingle and doMany via generic a generic method.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func do[T any](c *Client, ctx context.Context, req *http.Request) (*T, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http error")
	}

	defer resp.Body.Close()

	data := new(T)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling payload")
	}

	return data, nil
}

func getIDFieldIndex(obj interface{}) int {
	v := reflect.ValueOf(obj)

	// For this to work, T needs to be a pointer.
	if v.Type().Kind() != reflect.Pointer {
		return -1
	}

	v = v.Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get("json"), ",")
		if len(tags) > 0 && tags[0] == "id" {
			return i
		}
	}

	return -1
}

// doSingle performs a request for a single payload.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func doSingle[T any](c *Client, ctx context.Context, req *http.Request) (*SingleResponse[T], error) {
	res, err := do[SingleResponse[T]](c, ctx, req)
	if err != nil {
		return nil, err
	}

	// Inject ID into T, if it exists.
	idFieldIndex := getIDFieldIndex(res.Data.Attributes)
	if idFieldIndex != -1 {
		reflect.ValueOf(res.Data.Attributes).Elem().Field(idFieldIndex).Set(reflect.ValueOf(res.Data.ID))
	}

	return res, err
}

// doMany performs a request for a paginated payload.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func doMany[T any](c *Client, ctx context.Context, req *http.Request) (*ManyResponse[T], error) {
	res, err := do[ManyResponse[T]](c, ctx, req)
	if err != nil {
		return nil, err
	}

	// Inject ID into T, if it exists.
	if len(res.Data) > 0 {
		idFieldIndex := getIDFieldIndex(res.Data[0].Attributes)

		if idFieldIndex != -1 {
			for _, obj := range res.Data {
				reflect.ValueOf(obj.Attributes).Elem().Field(idFieldIndex).Set(reflect.ValueOf(obj.ID))
			}
		}
	}

	return res, err
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
