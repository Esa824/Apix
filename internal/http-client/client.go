// File: client.go
package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

func NewClient(timeout time.Duration) *Client {
	client := resty.New().
		SetTimeout(timeout).
		SetHeader("User-Agent", "GoRestyClient/1.0")

	return &Client{resty: client}
}

type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        any
	Cookies     map[string]string
	Auth        *BasicAuth
	Context     context.Context
}

type BasicAuth struct {
	Username string
	Password string
}

func (c *Client) Do(opts RequestOptions) (*resty.Response, error) {
	req := c.resty.R()

	if opts.Context != nil {
		req = req.SetContext(opts.Context)
	}

	if opts.Headers != nil {
		req = req.SetHeaders(opts.Headers)
	}

	if opts.QueryParams != nil {
		req = req.SetQueryParams(opts.QueryParams)
	}

	if opts.Cookies != nil {
		for k, v := range opts.Cookies {
			req = req.SetCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}

	if opts.Body != nil {
		req = req.SetBody(opts.Body)
	}

	if opts.Auth != nil {
		req = req.SetBasicAuth(opts.Auth.Username, opts.Auth.Password)
	}

	return req.Execute(opts.Method, opts.URL)
}

// Convenience methods for common HTTP methods:

func (c *Client) Get(url string, headers map[string]string, query map[string]string) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:      "GET",
		URL:         url,
		Headers:     headers,
		QueryParams: query,
	})
}

func (c *Client) Post(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "POST",
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c *Client) Put(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "PUT",
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c *Client) Patch(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "PATCH",
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

func (c *Client) Delete(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "DELETE",
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}
