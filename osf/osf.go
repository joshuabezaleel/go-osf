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
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	defaultBaseURL     = "https://api.osf.io/v2/"
	defaultBaseTestURL = "https://api.test.osf.io/v2/"

	userAgent = "go-osf"

	TypePreprints         = "preprints"
	TypeProviders         = "providers"
	TypePreprintProviders = "preprint_providers"
	TypeFiles             = "files"
)

type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	BaseURL *url.URL

	UserAgent string

	rateMu sync.Mutex

	common service

	Preprints         *PreprintsService
	PreprintProviders *PreprintProvidersService
	Files             *FilesService
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
	c.Preprints = (*PreprintsService)(&c.common)
	c.PreprintProviders = (*PreprintProvidersService)(&c.common)
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

type Meta map[string]interface{}

// https://jsonapi.org/format/#document-links
type Link struct {
	Href string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
}

type Links struct {
	Self    *Link `json:"self,omitempty"`
	Related *Link `json:"related,omitempty"`
}

// https://jsonapi.org/format/#document-resource-object-relationships
type Relationship struct {
	Links *Links                          `json:"links,omitempty"`
	Data  *Data[interface{}, interface{}] `json:"data,omitempty"`
	Meta  Meta                            `json:"meta,omitempty"`
}

type Relationships map[string]Relationship

type Data[T any, U any] struct {
	Type string  `json:"type"`
	ID   *string `json:"id,omitempty"`

	Attributes    T             `json:"attributes,omitempty"`
	Links         U             `json:"links,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

type ErrorSource struct {
	Pointer string `json:"pointer"`
}

type Error struct {
	Source *ErrorSource `json:"source,omitempty"`
	Detail string       `json:"detail"`
}

func (e *Error) Error() string {
	msg := e.Detail
	if e.Source != nil {
		msg += " (" + e.Source.Pointer + ")"
	}
	return msg
}

type Errors []*Error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}

	msgs := make([]string, 0)
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return "multiple errors: " + strings.Join(msgs, ", ")
}

type SinglePayload[T any, U any] struct {
	RawData *Data[T, U] `json:"data,omitempty"`
	Errors  Errors      `json:"errors,omitempty"`

	Data T `json:"-"`
}

type PaginationMeta struct {
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
	Page    int `json:"-"`
}

type PaginationLinks struct {
	First *string         `json:"first"`
	Last  *string         `json:"last"`
	Prev  *string         `json:"prev"`
	Next  *string         `json:"next"`
	Meta  *PaginationMeta `json:"meta"`
}

type ManyPayload[T any, U any] struct {
	RawData         []*Data[T, U]    `json:"data"`
	PaginationLinks *PaginationLinks `json:"links"`
	Errors          Errors           `json:"errors,omitempty"`

	Data           []T             `json:"-"`
	PaginationMeta *PaginationMeta `json:"-"`
}

// do performs logic for doSingle and doMany via generic a generic method.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func do[T any](c *Client, ctx context.Context, req *http.Request) (*T, error) {
	// spew.Dump(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http error")
	}

	defer resp.Body.Close()

	data := new(T)

	// TODO: Remove this.
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// ioutil.WriteFile("dump.json", body, 0644)
	// if err := json.NewDecoder(bytes.NewReader(body)).Decode(data); err != nil {
	// 	return nil, errors.Wrap(err, "error unmarshaling payload")
	// }

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

type BuildDataFn[T any, U any] func(obj *Data[T, U]) (T, error)

// doSingle performs a request for a single payload.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func doSingle[T any, U any](c *Client, ctx context.Context, req *http.Request, build ...BuildDataFn[T, U]) (*SinglePayload[T, U], error) {
	res, err := do[SinglePayload[T, U]](c, ctx, req)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		return res, res.Errors
	}

	// Inject ID into Attributes, if it exists.
	if res.RawData.ID != nil {
		idFieldIndex := getIDFieldIndex(res.RawData.Attributes)
		if idFieldIndex != -1 {
			reflect.ValueOf(res.RawData.Attributes).Elem().Field(idFieldIndex).Set(reflect.ValueOf(res.RawData.ID).Elem())
		}
	}

	if len(build) > 0 {
		res.Data, err = build[0](res.RawData)
		if err != nil {
			return nil, err
		}
	} else {
		res.Data = res.RawData.Attributes
	}

	return res, err
}

// doMany performs a request for a paginated payload.
// HACK: since Go has not supported generics for struct methods (yet), we need to make this standalone.
func doMany[T any, U any](c *Client, ctx context.Context, req *http.Request, build ...BuildDataFn[T, U]) (*ManyPayload[T, U], error) {
	res, err := do[ManyPayload[T, U]](c, ctx, req)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		return res, res.Errors
	}

	// Inject ID into T, if it exists.
	if len(res.Data) > 0 {
		idFieldIndex := getIDFieldIndex(res.RawData[0].Attributes)

		if idFieldIndex != -1 {
			for _, obj := range res.RawData {
				if obj.ID != nil {
					reflect.ValueOf(obj.Attributes).Elem().Field(idFieldIndex).Set(reflect.ValueOf(obj.ID).Elem())
				}
			}
		}
	}

	res.Data = make([]T, 0, len(res.RawData))
	if len(build) > 0 {
		for _, raw := range res.RawData {
			obj, err := build[0](raw)
			if err != nil {
				return nil, err
			}
			res.Data = append(res.Data, obj)
		}
	} else {
		for _, raw := range res.RawData {
			res.Data = append(res.Data, raw.Attributes)
		}
	}

	// Inject PaginationMeta into payload.
	if res.PaginationLinks != nil && res.PaginationLinks.Meta != nil {
		res.PaginationMeta = res.PaginationLinks.Meta

		// Get pagination info from query string, since PaginationLinks are not always reliable.

		page := 1
		if req.URL.Query().Has("page[number]") {
			if num, err := strconv.Atoi(req.URL.Query().Get("page[number]")); err == nil {
				page = num
			}
		}
		res.PaginationMeta.Page = page

		perPage := res.PaginationMeta.PerPage
		if perPage == 0 {
			perPage = 10
			if req.URL.Query().Has("page[size]") {
				if num, err := strconv.Atoi(req.URL.Query().Get("page[size]")); err == nil {
					perPage = num
				}
			}
		}
		res.PaginationMeta.PerPage = perPage
	}

	return res, err
}

type ListOptions struct {
	Page    int `url:"page[number],omitempty"`
	PerPage int `url:"page[size],omitempty"`

	Filter map[string]string `url:"-"`
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}, additionalQueries ...map[string]string) (string, error) {
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

// addOptions adds the parameters in opts as URL query parameters to s. opts
// along with queries for filtering on list endpoint.
// must be a struct whose fields may contain "url" tags.
func addOptionsWithFilter(s string, opts interface{}, additionalQueries ...map[string]string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	var qs url.Values
	v := reflect.ValueOf(opts)
	if !(v.Kind() == reflect.Ptr && v.IsNil()) {
		qs, err = query.Values(opts)
		if err != nil {
			return s, err
		}
	}

	for _, q := range additionalQueries {
		for k, v := range q {
			qs.Add("filter["+k+"]", v)
		}
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
