package osf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

// func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return resp, err
// 	}
// 	defer resp.Body.Close()

// 	switch v := v.(type) {
// 	case nil:
// 	case io.Writer:
// 		_, err = io.Copy(v, resp.Body)
// 	default:
// 		decErr := json.NewDecoder(resp.Body).Decode(v)
// 		if decErr == io.EOF {
// 			decErr = nil // ignore EOF errors caused by empty response body
// 		}
// 		if decErr != nil {
// 			err = decErr
// 		}
// 	}
// 	return resp, err
// }

// type Response struct {
// }

// func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
// 	if ctx == nil {
// 		return nil, errNonNilContext
// 	}

// 	req = withContext(ctx, req)

// 	rateLimitCategory := category(req.URL.Path)

// 	if bypass := ctx.Value(bypassRateLimitCheck); bypass == nil {
// 		// If we've hit rate limit, don't make further requests before Reset time.
// 		if err := c.checkRateLimitBeforeDo(req, rateLimitCategory); err != nil {
// 			return &Response{
// 				Response: err.Response,
// 				Rate:     err.Rate,
// 			}, err
// 		}
// 	}

// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		// If we got an error, and the context has been canceled,
// 		// the context's error is probably more useful.
// 		select {
// 		case <-ctx.Done():
// 			return nil, ctx.Err()
// 		default:
// 		}

// 		// If the error type is *url.Error, sanitize its URL before returning.
// 		if e, ok := err.(*url.Error); ok {
// 			if url, err := url.Parse(e.URL); err == nil {
// 				e.URL = sanitizeURL(url).String()
// 				return nil, e
// 			}
// 		}

// 		return nil, err
// 	}

// 	response := newResponse(resp)

// 	// Don't update the rate limits if this was a cached response.
// 	// X-From-Cache is set by https://github.com/gregjones/httpcache
// 	if response.Header.Get("X-From-Cache") == "" {
// 		c.rateMu.Lock()
// 		c.rateLimits[rateLimitCategory] = response.Rate
// 		c.rateMu.Unlock()
// 	}

// 	err = CheckResponse(resp)
// 	if err != nil {
// 		defer resp.Body.Close()
// 		// Special case for AcceptedErrors. If an AcceptedError
// 		// has been encountered, the response's payload will be
// 		// added to the AcceptedError and returned.
// 		//
// 		// Issue #1022
// 		aerr, ok := err.(*AcceptedError)
// 		if ok {
// 			b, readErr := ioutil.ReadAll(resp.Body)
// 			if readErr != nil {
// 				return response, readErr
// 			}

// 			aerr.Raw = b
// 			err = aerr
// 		}
// 	}
// 	return response, err
// }
