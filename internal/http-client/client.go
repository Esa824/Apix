// File: client.go
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Esa824/apix/internal/model"
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
	Files       map[string]string
	Body        any
	Cookies     map[string]string
	Auth        *BasicAuth
	Context     context.Context
	Time        time.Time
	IsTemplate  bool
	Name        string
}

type BasicAuth struct {
	Username string
	Password string
}

func (c *Client) Do(opts RequestOptions, saveToHistory bool) (*resty.Response, error) {
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

	if opts.Files != nil {
		req = req.SetFiles(opts.Files)
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
		if saveToHistory {
			UpdateHistory(opts)
		}

		if opts.IsTemplate {
			SaveTemplate(model.Template{
				Name:    opts.Name,
				Method:  opts.Method,
				URL:     opts.URL,
				Headers: opts.Headers,
				Body:    opts.Body,
			})
		}
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
	}, true)
}

func (c *Client) Post(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "POST",
		URL:     url,
		Headers: headers,
		Body:    body,
	}, true)
}

func (c *Client) Put(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "PUT",
		URL:     url,
		Headers: headers,
		Body:    body,
	}, true)
}

func (c *Client) Patch(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "PATCH",
		URL:     url,
		Headers: headers,
		Body:    body,
	}, true)
}

func (c *Client) Delete(url string, headers map[string]string, body any) (*resty.Response, error) {
	return c.Do(RequestOptions{
		Method:  "DELETE",
		URL:     url,
		Headers: headers,
		Body:    body,
	}, true)
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

// initTemplatesDir ensures the templates directory exists
func initTemplatesDir() error {
	templatesDir := filepath.Join(ConfigPath, "templates")
	return os.MkdirAll(templatesDir, 0755)
}

// saveTemplate saves a single template to a JSON file
func saveTemplate(template *model.Template) error {
	if err := initTemplatesDir(); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	templatesDir := filepath.Join(ConfigPath, "templates")
	filename := fmt.Sprintf("%s.json", template.Name)
	filepath := filepath.Join(templatesDir, filename)

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	return os.WriteFile(filepath, data, 0600)
}

// deleteTemplateFile deletes the JSON file for a template
func deleteTemplateFile(templateName string) error {
	templatesDir := filepath.Join(ConfigPath, "templates")
	filename := fmt.Sprintf("%s.json", templateName)
	filepath := filepath.Join(templatesDir, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete template file: %w", err)
	}

	return nil
}

// loadTemplates loads all templates from the templates directory
func loadTemplates() (map[string]*model.Template, error) {
	templatesDir := filepath.Join(ConfigPath, "templates")

	// Check if directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return make(map[string]*model.Template), nil // No templates directory yet, return empty map
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	templates := make(map[string]*model.Template)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filepath := filepath.Join(templatesDir, entry.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			// You might want to use your utils.ShowWarning here instead
			fmt.Printf("Warning: Failed to read template file %s: %v\n", entry.Name(), err)
			continue
		}

		var template model.Template
		if err := json.Unmarshal(data, &template); err != nil {
			fmt.Printf("Warning: Failed to parse template file %s: %v\n", entry.Name(), err)
			continue
		}

		templates[template.Name] = &template
	}

	return templates, nil
}

// GetTemplates returns all saved templates as a slice (maintains compatibility with existing API)
func GetTemplates() ([]model.Template, error) {
	templatesMap, err := loadTemplates()
	if err != nil {
		return nil, err
	}

	templates := make([]model.Template, 0, len(templatesMap))
	for _, template := range templatesMap {
		templates = append(templates, *template)
	}

	return templates, nil
}

// SaveTemplate saves a new template to its own file
func SaveTemplate(template model.Template) error {
	// Remove the ID assignment since we're not using sequential IDs anymore
	// The template name serves as the unique identifier
	return saveTemplate(&template)
}

// UpdateTemplate updates an existing template by name
func UpdateTemplate(template model.Template) error {
	// Check if template exists first
	templatesMap, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load existing templates: %w", err)
	}

	if _, exists := templatesMap[template.Name]; !exists {
		return fmt.Errorf("template with name '%s' not found", template.Name)
	}

	return saveTemplate(&template)
}

// DeleteTemplate removes a template by name
func DeleteTemplate(name string) error {
	// Check if template exists first
	templatesMap, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load existing templates: %w", err)
	}

	if _, exists := templatesMap[name]; !exists {
		return fmt.Errorf("template with name '%s' not found", name)
	}

	return deleteTemplateFile(name)
}

// DeleteAllTemplates deletes the entire templates directory
func DeleteAllTemplates() error {
	templatesDir := filepath.Join(ConfigPath, "templates")
	if err := os.RemoveAll(templatesDir); err != nil {
		return fmt.Errorf("failed to delete templates directory: %w", err)
	}
	return nil
}

// GetTemplateByID returns a specific template by its ID (deprecated - use GetTemplateByName instead)
// This function is kept for backward compatibility but will need the Template model to have an ID field
func GetTemplateByID(id int) (*model.Template, error) {
	templates, err := GetTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}

	for _, template := range templates {
		if template.Id == id {
			return &template, nil
		}
	}

	return nil, fmt.Errorf("template with ID %d not found", id)
}

// GetTemplateByName returns a specific template by its name
func GetTemplateByName(name string) (*model.Template, error) {
	templatesMap, err := loadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	if template, exists := templatesMap[name]; exists {
		return template, nil
	}

	return nil, fmt.Errorf("template with name '%s' not found", name)
}
