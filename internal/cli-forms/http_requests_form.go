package cliforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	hc "apix/internal/http-client"
)

// HTTPResponse holds the response data and provides querying capabilities
type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    map[string]string
	Body       []byte
	IsJSON     bool
	ParsedJSON interface{}
}

func HandleHttpRequests() {
	var selectedOption string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("HTTP Requests:").
				Options(
					huh.NewOption("GET Request", "get"),
					huh.NewOption("POST Request", "post"),
					huh.NewOption("PUT Request", "put"),
					huh.NewOption("PATCH Request", "patch"),
					huh.NewOption("DELETE Request", "delete"),
					huh.NewOption("Back to Main Menu", "back"),
				).
				Value(&selectedOption),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error running HTTP requests form", err)
		return
	}

	handleHttpSelection(selectedOption)
}

func handleHttpSelection(selection string) {
	switch selection {
	case "get":
		handleGetRequest()
	case "post":
		handlePostRequest()
	case "put":
		handlePutRequest()
	case "patch":
		handlePatchRequest()
	case "delete":
		handleDeleteRequest()
	case "back":
		RunInteractiveMode()
	default:
		showMessage("Unknown HTTP request option", "error")
		HandleHttpRequests()
	}
}

// GET Request Handler
func handleGetRequest() {
	endpoint := getEndpointInput()
	if endpoint == "" {
		showMessage("No endpoint provided. Returning to menu.", "info")
		HandleHttpRequests()
		return
	}

	options := handleRequestOptions("GET")
	response, err := hc.NewClient(10*time.Second).Get(endpoint, options.Headers, options.QueryParams)
	if err != nil {
		showError("Error while calling endpoint", err)
		HandleHttpRequests()
		return
	}

	handleResponse(response)
}

// POST Request Handler
func handlePostRequest() {
	endpoint := getEndpointInput()
	if endpoint == "" {
		showMessage("No endpoint provided. Returning to menu.", "info")
		HandleHttpRequests()
		return
	}

	_, body := handleBodyTypeSelection()
	options := handleRequestOptions("POST")
	response, err := hc.NewClient(10*time.Second).Post(endpoint, options.Headers, body)
	if err != nil {
		showError("Error while calling endpoint", err)
		HandleHttpRequests()
		return
	}

	handleResponse(response)
}

// PUT Request Handler
func handlePutRequest() {
	endpoint := getEndpointInput()
	if endpoint == "" {
		showMessage("No endpoint provided. Returning to menu.", "info")
		HandleHttpRequests()
		return
	}

	_, body := handleBodyTypeSelection()
	options := handleRequestOptions("PUT")
	response, err := hc.NewClient(10*time.Second).Put(endpoint, options.Headers, body)
	if err != nil {
		showError("Error while calling endpoint", err)
		HandleHttpRequests()
		return
	}

	handleResponse(response)
}

// PATCH Request Handler
func handlePatchRequest() {
	endpoint := getEndpointInput()
	if endpoint == "" {
		showMessage("No endpoint provided. Returning to menu.", "info")
		HandleHttpRequests()
		return
	}

	_, body := handlePatchBodyTypeSelection()
	options := handleRequestOptions("PATCH")
	response, err := hc.NewClient(10*time.Second).Patch(endpoint, options.Headers, body)
	if err != nil {
		showError("Error while calling endpoint", err)
		HandleHttpRequests()
		return
	}

	handleResponse(response)
}

// DELETE Request Handler
func handleDeleteRequest() {
	endpoint := getEndpointInput()
	if endpoint == "" {
		showMessage("No endpoint provided. Returning to menu.", "info")
		HandleHttpRequests()
		return
	}

	options := handleRequestOptions("DELETE")

	// Confirmation prompt for DELETE
	var confirmDelete bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm DELETE Request").
				Description(fmt.Sprintf("Are you sure you want to send a DELETE request to: %s", endpoint)).
				Affirmative("Yes, delete").
				Negative("Cancel").
				Value(&confirmDelete),
		),
	)

	err := confirmForm.Run()
	if err != nil {
		showError("Error running confirmation form", err)
		return
	}

	if !confirmDelete {
		showMessage("DELETE request cancelled.", "info")
		HandleHttpRequests()
		return
	}

	// Fix: Use Delete method instead of Patch
	response, err := hc.NewClient(10*time.Second).Delete(endpoint, options.Headers, options)
	if err != nil {
		showError("Error while calling endpoint", err)
		HandleHttpRequests()
		return
	}

	handleResponse(response)
}

// Helper function to get endpoint input
func getEndpointInput() string {
	var endpoint string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter API endpoint:").
				Description("Will be appended to base URL if configured").
				Placeholder("/api/users or https://api.example.com/users").
				Value(&endpoint),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error getting endpoint input", err)
		return ""
	}

	return strings.TrimSpace(endpoint)
}

// Enhanced response handler with JSON querying
func handleResponse(response interface{}) {
	// Convert response to our HTTPResponse type
	httpResp := parseResponse(response)

	// Display the response
	displayResponse(httpResp)

	// If it's JSON, offer querying options
	if httpResp.IsJSON {
		handleJSONQuerying(httpResp)
	} else {
		askContinueOrReturnHttpRequests()
	}
}

func parseResponse(response interface{}) *HTTPResponse {
	// This is a simplified version - you'll need to adapt based on your actual response type
	// Assuming your response has methods like Body(), Status(), etc.

	var body []byte
	var status string

	// Use reflection or type assertion to extract data from response
	// This is a placeholder - adjust based on your actual response structure
	if resp, ok := response.(interface{ Body() []byte }); ok {
		body = resp.Body()
	}

	if resp, ok := response.(interface{ Status() string }); ok {
		status = resp.Status()
	}

	formatted, isJSON := FormatJSON(body)

	httpResp := &HTTPResponse{
		Status:  status,
		Body:    formatted,
		IsJSON:  isJSON,
		Headers: make(map[string]string),
	}

	if isJSON {
		json.Unmarshal(formatted, &httpResp.ParsedJSON)
	}

	return httpResp
}

func displayResponse(response *HTTPResponse) {
	responseText := fmt.Sprintf("Status: %s\n\nBody:\n%s",
		response.Status,
		string(response.Body))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("HTTP Response").
				Description(responseText),
		),
	)

	form.Run()
}

func handleJSONQuerying(response *HTTPResponse) {
	for {
		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("JSON Response Options").
					Options(
						huh.NewOption("Query JSON", "query"),
						huh.NewOption("View Raw Response", "raw"),
						huh.NewOption("Continue", "continue"),
					).
					Value(&choice),
			),
		)

		err := form.Run()
		if err != nil {
			showError("Error in JSON options", err)
			break
		}

		switch choice {
		case "query":
			queryJSON(response)
		case "raw":
			displayRawResponse(response)
		case "continue":
			askContinueOrReturnHttpRequests()
			return
		}
	}
}

func queryJSON(response *HTTPResponse) {
	var query string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter JSON Query").
				Description("Examples: .name, .users[0].email, .data.items").
				Placeholder(".field.subfield").
				Value(&query),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error getting query input", err)
		return
	}

	if strings.TrimSpace(query) == "" {
		showMessage("No query provided", "info")
		return
	}

	result := executeJSONQuery(response.ParsedJSON, query)
	displayQueryResult(query, result)
}

func executeJSONQuery(data interface{}, query string) interface{} {
	if data == nil {
		return "null"
	}

	// Remove leading dot if present
	query = strings.TrimPrefix(query, ".")

	if query == "" {
		return data
	}

	parts := strings.Split(query, ".")
	current := data

	for _, part := range parts {
		if current == nil {
			return "null"
		}

		// Handle array indexing
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			current = handleArrayAccess(current, part)
		} else {
			current = handleFieldAccess(current, part)
		}

		if current == nil {
			return "Field not found"
		}
	}

	return current
}

func handleArrayAccess(data interface{}, accessor string) interface{} {
	parts := strings.Split(accessor, "[")
	fieldName := parts[0]
	indexStr := strings.TrimSuffix(parts[1], "]")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "Invalid array index"
	}

	// Get the field first if it exists
	if fieldName != "" {
		data = handleFieldAccess(data, fieldName)
		if data == nil {
			return nil
		}
	}

	// Handle array access
	switch arr := data.(type) {
	case []interface{}:
		if index < 0 || index >= len(arr) {
			return "Array index out of bounds"
		}
		return arr[index]
	default:
		return "Not an array"
	}
}

func handleFieldAccess(data interface{}, field string) interface{} {
	switch obj := data.(type) {
	case map[string]interface{}:
		return obj[field]
	case map[interface{}]interface{}:
		return obj[field]
	default:
		return nil
	}
}

func displayQueryResult(query string, result interface{}) {
	var resultStr string

	if result == nil {
		resultStr = "null"
	} else {
		switch v := result.(type) {
		case string:
			resultStr = v
		case map[string]interface{}, []interface{}:
			if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
				resultStr = string(jsonBytes)
			} else {
				resultStr = fmt.Sprintf("%v", v)
			}
		default:
			resultStr = fmt.Sprintf("%v", v)
		}
	}

	displayText := fmt.Sprintf("Query: %s\n\nResult:\n%s", query, resultStr)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Query Result").
				Description(displayText),
		),
	)

	form.Run()
}

func displayRawResponse(response *HTTPResponse) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Raw Response").
				Description(string(response.Body)),
		),
	)

	form.Run()
}

// Body Type Selection for POST/PUT
func handleBodyTypeSelection() (string, string) {
	var bodyType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Body Type:").
				Options(
					huh.NewOption("JSON", "json"),
					huh.NewOption("Form Data", "form"),
					huh.NewOption("Multipart Form", "multipart"),
					huh.NewOption("Raw Text", "raw"),
					huh.NewOption("File Upload", "file"),
					huh.NewOption("No Body", "none"),
				).
				Value(&bodyType),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error selecting body type", err)
		return "none", ""
	}

	return bodyType, handleBodyInput(bodyType)
}

// Body Type Selection for PATCH (limited options)
func handlePatchBodyTypeSelection() (string, string) {
	var bodyType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Body Type:").
				Options(
					huh.NewOption("JSON", "json"),
					huh.NewOption("Form Data", "form"),
					huh.NewOption("Raw Text", "raw"),
					huh.NewOption("No Body", "none"),
				).
				Value(&bodyType),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error selecting body type", err)
		return "none", ""
	}

	return bodyType, handleBodyInput(bodyType)
}

// Handle Body Input based on type
func handleBodyInput(bodyType string) string {
	switch bodyType {
	case "json":
		return handleJSONInput()
	case "form":
		return handleFormDataInput()
	case "multipart":
		return handleMultipartFormInput()
	case "raw":
		return handleRawTextInput()
	case "file":
		return handleFileUploadInput()
	default:
		return ""
	}
}

func handleJSONInput() string {
	var jsonBody string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter JSON Body:").
				Placeholder(`{"name": "John Doe", "email": "john@example.com"}`).
				Value(&jsonBody),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error getting JSON input", err)
		return ""
	}

	// Validate JSON
	if strings.TrimSpace(jsonBody) != "" {
		var temp interface{}
		if err := json.Unmarshal([]byte(jsonBody), &temp); err != nil {
			showError("Invalid JSON format", err)
			return handleJSONInput() // Retry
		}
	}

	return jsonBody
}

func handleFormDataInput() string {
	var formData string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter Form Data:").
				Description("Format: key1=value1&key2=value2").
				Placeholder("name=John Doe&email=john@example.com").
				Value(&formData),
		),
	)

	form.Run()
	return formData
}

func handleMultipartFormInput() string {
	var multipartData string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter Multipart Form Data:").
				Description("Format: key1=value1&key2=value2").
				Placeholder("name=John Doe&email=john@example.com&file=@/path/to/file").
				Value(&multipartData),
		),
	)

	form.Run()
	return multipartData
}

func handleRawTextInput() string {
	var rawText string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter Raw Text:").
				Placeholder("Enter your raw text content here...").
				Value(&rawText),
		),
	)

	form.Run()
	return rawText
}

func handleFileUploadInput() string {
	var filePath, fieldName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter file path:").
				Placeholder("/path/to/your/file.jpg").
				Value(&filePath),
			huh.NewInput().
				Title("Enter field name:").
				Placeholder("image").
				Value(&fieldName),
		),
	)

	form.Run()
	return fmt.Sprintf("file:%s;field:%s", filePath, fieldName)
}

// Request Options Structure
type RequestOptions struct {
	Headers        map[string]string
	QueryParams    map[string]string
	AuthType       string
	AuthValue      string
	Files          []FileUpload
	SaveAsTemplate bool
}

type FileUpload struct {
	Path      string
	FieldName string
}

// Handle Request Options
func handleRequestOptions(method string) RequestOptions {
	options := RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	// Add Headers?
	var addHeaders bool
	headerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add Headers?").
				Value(&addHeaders),
		),
	)
	headerForm.Run()

	if addHeaders {
		options.Headers = handleHeaders()
	}

	// Add Query Parameters?
	var addParams bool
	paramForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add Query Parameters?").
				Value(&addParams),
		),
	)
	paramForm.Run()

	if addParams {
		options.QueryParams = handleQueryParams()
	}

	// Authentication?
	var addAuth bool
	authForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add Authentication?").
				Value(&addAuth),
		),
	)
	authForm.Run()

	if addAuth {
		options.AuthType, options.AuthValue = handleAuthentication()
	}

	// Add Files/Images? (for POST/PUT only)
	if method == "POST" || method == "PUT" {
		var addFiles bool
		fileForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add Files/Images?").
					Value(&addFiles),
			),
		)
		fileForm.Run()

		if addFiles {
			options.Files = handleFileUploads()
		}
	}

	// Save as Template?
	var saveTemplate bool
	templateForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Save as Template?").
				Value(&saveTemplate),
		),
	)
	templateForm.Run()

	options.SaveAsTemplate = saveTemplate

	return options
}

func handleHeaders() map[string]string {
	headers := make(map[string]string)

	for {
		var key, value string
		var addMore bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Header Key:").
					Placeholder("Content-Type").
					Value(&key),
				huh.NewInput().
					Title("Header Value:").
					Placeholder("application/json").
					Value(&value),
				huh.NewConfirm().
					Title("Add another header?").
					Value(&addMore),
			),
		)

		err := form.Run()
		if err != nil || strings.TrimSpace(key) == "" {
			break
		}

		headers[strings.TrimSpace(key)] = strings.TrimSpace(value)

		if !addMore {
			break
		}
	}

	return headers
}

func handleQueryParams() map[string]string {
	params := make(map[string]string)

	for {
		var key, value string
		var addMore bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Parameter Key:").
					Placeholder("page").
					Value(&key),
				huh.NewInput().
					Title("Parameter Value:").
					Placeholder("1").
					Value(&value),
				huh.NewConfirm().
					Title("Add another parameter?").
					Value(&addMore),
			),
		)

		err := form.Run()
		if err != nil || strings.TrimSpace(key) == "" {
			break
		}

		params[strings.TrimSpace(key)] = strings.TrimSpace(value)

		if !addMore {
			break
		}
	}

	return params
}

func handleAuthentication() (string, string) {
	var authType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Authentication Type:").
				Options(
					huh.NewOption("Bearer Token", "bearer"),
					huh.NewOption("API Key", "apikey"),
					huh.NewOption("Basic Auth", "basic"),
				).
				Value(&authType),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error selecting auth type", err)
		return "", ""
	}

	switch authType {
	case "bearer":
		return handleBearerToken()
	case "apikey":
		return handleAPIKey()
	case "basic":
		return handleBasicAuth()
	default:
		return "", ""
	}
}

func handleBearerToken() (string, string) {
	var token string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Bearer Token:").
				Placeholder("your-jwt-token-here").
				Password(true).
				Value(&token),
		),
	)

	form.Run()
	return "bearer", token
}

func handleAPIKey() (string, string) {
	var apiKey, headerName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter API Key:").
				Password(true).
				Value(&apiKey),
			huh.NewInput().
				Title("Header Name (optional):").
				Placeholder("X-API-Key").
				Value(&headerName),
		),
	)

	form.Run()
	if headerName == "" {
		headerName = "X-API-Key"
	}
	return "apikey", fmt.Sprintf("%s:%s", headerName, apiKey)
}

func handleBasicAuth() (string, string) {
	var username, password string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username:").
				Value(&username),
			huh.NewInput().
				Title("Password:").
				Password(true).
				Value(&password),
		),
	)

	form.Run()
	return "basic", fmt.Sprintf("%s:%s", username, password)
}

func handleFileUploads() []FileUpload {
	var files []FileUpload

	for {
		var filePath, fieldName string
		var addMore bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("File Path:").
					Placeholder("/path/to/file.jpg").
					Value(&filePath),
				huh.NewInput().
					Title("Field Name:").
					Placeholder("image").
					Value(&fieldName),
				huh.NewConfirm().
					Title("Add another file?").
					Value(&addMore),
			),
		)

		err := form.Run()
		if err != nil || strings.TrimSpace(filePath) == "" {
			break
		}

		files = append(files, FileUpload{
			Path:      strings.TrimSpace(filePath),
			FieldName: strings.TrimSpace(fieldName),
		})

		if !addMore {
			break
		}
	}

	return files
}

func FormatJSON(data []byte) ([]byte, bool) {
	if len(data) == 0 {
		return data, false
	}

	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return data, false
	}

	formatted, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return data, false
	}

	formatted = bytes.TrimSuffix(formatted, []byte("\n"))
	return formatted, true
}

func askContinueOrReturnHttpRequests() {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
				Options(
					huh.NewOption("Make Another Request", "continue"),
					huh.NewOption("Return to Main Menu", "main"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&choice),
		),
	)

	err := form.Run()
	if err != nil {
		showError("Error in continuation menu", err)
		return
	}

	switch choice {
	case "continue":
		HandleHttpRequests()
	case "main":
		RunInteractiveMode()
	case "exit":
		showMessage("Goodbye!", "info")
		os.Exit(0)
	}
}

// Helper functions for consistent UI messaging
func showMessage(message, msgType string) {
	var title string
	switch msgType {
	case "error":
		title = "Error"
	case "success":
		title = "Success"
	case "info":
		title = "Information"
	default:
		title = "Message"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(title).
				Description(message),
		),
	)

	form.Run()
}

func showError(message string, err error) {
	errorText := fmt.Sprintf("%s: %v", message, err)
	showMessage(errorText, "error")
}

