// File: client.go
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
)

var ConfigPath = "./testconfigs/config1/"

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
	Id          int
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        any
	Cookies     map[string]string
	Auth        *BasicAuth
	Context     context.Context
	Time        time.Time
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

	response, err := req.Execute(opts.Method, opts.URL)
	if err == nil {
		UpdateHistory(opts)
	}
	return response, err
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

// GetHistory reads and returns all request history from the history file
func GetHistory() ([]RequestOptions, error) {
	filepath := filepath.Join(ConfigPath, "history")

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// File doesn't exist, return empty array
		return []RequestOptions{}, nil
	}

	// Read file contents
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []RequestOptions{}, nil
	}

	// Unmarshal into array
	var history []RequestOptions
	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal history data: %w", err)
	}

	return history, nil
}

// UpdateHistory appends a new request to the history and saves it
func UpdateHistory(request RequestOptions) error {
	// Get existing history
	history, err := GetHistory()
	if err != nil {
		return fmt.Errorf("failed to get existing history: %w", err)
	}

	// Append new request to history
	history = append(history, request)

	if len(history) > 0 {
		history[len(history)-1].Id = len(history) - 1
	}

	// Marshal updated history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	// Write to file
	filepath := filepath.Join(ConfigPath, "history")
	err = os.WriteFile(filepath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// DeleteHistory deletes the history file
func DeleteHistory() error {
	filepath := filepath.Join(ConfigPath, "history")
	return fmt.Errorf("failed to delete history: %w", os.Remove(filepath))
}
