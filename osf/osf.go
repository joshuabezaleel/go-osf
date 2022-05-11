package osf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/google/jsonapi"
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

	log.Println(c.BaseURL.String())
	log.Println(urlStr)
	log.Println(u)

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

	// log.Println(u.String())
	log.Println("JOBEL")

	// TODO: why the /v2 gone
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	log.Println(req.URL)

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

type Response struct {
	*http.Response
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, errors.Wrap(err, "http error")
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Ptr:
			s := t.Elem()
			switch s.Kind() {
			case reflect.Struct:
				err := jsonapi.UnmarshalPayload(resp.Body, v)
				if err != nil {
					return resp, errors.Wrap(err, "failed to unmarshal payload")
				}
				return resp, nil

			case reflect.Slice:
				sliceType := s.Elem()
				if sliceType.Kind() != reflect.Ptr {
					return resp, errors.New("v should be a slice of pointers, not a slice of structs")
				}

				objsIface, err := jsonapi.UnmarshalManyPayload(resp.Body, sliceType)
				if err != nil {
					return resp, errors.Wrap(err, "failed to unmarshal payload")
				}

				results := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, len(objsIface))
				for _, obj := range objsIface {
					v := reflect.ValueOf(obj)
					if v.Type() != sliceType {
						return resp, errors.New("failed to unmarshal payload")
					}
					results = reflect.Append(results, reflect.ValueOf(obj))
				}

				reflect.ValueOf(v).Elem().Set(results)
				return resp, nil

			default:
				// Do nothing.
				return resp, nil
			}

		default:
			// Do nothing.
			return resp, nil
		}

		// decErr := json.NewDecoder(resp.Body).Decode(v)
		// if decErr == io.EOF {
		// 	decErr = nil // ignore EOF errors caused by empty response body
		// }
		// if decErr != nil {
		// 	err = decErr
		// }
	}
	return resp, nil
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
