package cliforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	hc "apix/internal/http-client"
	"apix/internal/utils"
)

// HTTPResponse holds the response data and provides querying capabilities
type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    map[string]string
	Body       []byte
	IsJSON     bool
	ParsedJSON any
}

type FileUpload struct {
	Path      string
	FieldName string
}

func HandleHttpRequests() {
	choice, err := utils.AskSelection("HTTP Requests:", []utils.SelectionOption{
		{"GET Request", "get"},
		{"POST Request", "post"},
		{"PUT Request", "put"},
		{"PATCH Request", "patch"},
		{"DELETE Request", "delete"},
		{"Back to Main Menu", "back"},
	})

	if err != nil {
		utils.ShowError("Error running HTTP requests form", err)
		return
	}

	handleHttpSelection(choice)
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
		utils.ShowWarning("Unknown HTTP request option")
		HandleHttpRequests()
	}
}

// GET Request Handler
func handleGetRequest() {
	endpoint, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter API endpoint:",
		Description: "Will be appended to base URL if configured",
		Placeholder: "/api/users or https://api.example.com/users",
		Required:    true,
	})

	if err != nil || endpoint == "" {
		utils.ShowMessage("No endpoint provided. Returning to menu.")
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	endpoint = fmt.Sprintf("%s%s", BaseURL, strings.TrimSpace(endpoint))
	options := handleRequestOptions("GET", endpoint, "")
	response, err := hc.NewClient(10 * time.Second).Do(options)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	handleResponse(response)
}

// POST Request Handler
func handlePostRequest() {
	endpoint, body := getEndpointAndBody("POST")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("POST", endpoint, body)
	response, err := hc.NewClient(10 * time.Second).Do(options)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	handleResponse(response)
}

// PUT Request Handler
func handlePutRequest() {
	endpoint, body := getEndpointAndBody("PUT")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("PUT", endpoint, body)
	response, err := hc.NewClient(10 * time.Second).Do(options)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	handleResponse(response)
}

// PATCH Request Handler
func handlePatchRequest() {
	endpoint, body := getEndpointAndBody("PATCH")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("PATCH", endpoint, body)
	response, err := hc.NewClient(10 * time.Second).Do(options)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	handleResponse(response)
}

// DELETE Request Handler
func handleDeleteRequest() {
	endpoint, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter API endpoint:",
		Description: "Will be appended to base URL if configured",
		Placeholder: "/api/users/1 or https://api.example.com/users/1",
		Required:    true,
	})

	if err != nil || endpoint == "" {
		utils.ShowMessage("No endpoint provided. Returning to menu.")
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	endpoint = fmt.Sprintf("%s%s", BaseURL, strings.TrimSpace(endpoint))

	// Confirmation prompt for DELETE
	confirmed, err := utils.AskDangerousConfirmation(
		"Delete Resource",
		"Are you sure you want to send a DELETE request to",
		endpoint,
	)

	if err != nil {
		utils.ShowError("Error running confirmation", err)
		return
	}

	if !confirmed {
		utils.ShowMessage("DELETE request cancelled.")
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	options := handleRequestOptions("DELETE", endpoint, "")
	response, err := hc.NewClient(10 * time.Second).Do(options)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	handleResponse(response)
}

// Helper function to get endpoint and body for POST/PUT/PATCH
func getEndpointAndBody(method string) (string, string) {
	inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
		{
			Title:       "Enter API endpoint:",
			Description: "Will be appended to base URL if configured",
			Placeholder: "/api/users or https://api.example.com/users",
			Required:    true,
		},
	})

	if err != nil || len(inputs) == 0 || inputs[0] == "" {
		utils.ShowMessage("No endpoint provided. Returning to menu.")
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return "", ""
	}

	fullEndpoint := fmt.Sprintf("%s%s", BaseURL, strings.TrimSpace(inputs[0]))
	_, body := handleBodyTypeSelection(method)

	return fullEndpoint, body
}

// Enhanced response handler with JSON querying
func handleResponse(response any) {
	httpResp := parseResponse(response)
	displayResponse(httpResp)

	if httpResp.IsJSON {
		handleJSONQuerying(httpResp)
	} else {
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Another Request", "Main Menu")
	}
}

func parseResponse(response any) *HTTPResponse {
	var body []byte
	var status string

	// Use reflection or type assertion to extract data from response
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
	if string(response.Body) == "" {
		response.Body = []byte("Not set")
	}
	responseText := fmt.Sprintf("Status: %s\n\nBody:\n%s",
		response.Status,
		string(response.Body))

	utils.DisplayFormattedText("🌐 HTTP Response", responseText)
}

func handleJSONQuerying(response *HTTPResponse) {
	for {
		choice, err := utils.AskSelection("JSON Response Options", []utils.SelectionOption{
			{"Query JSON", "query"},
			{"View Response", "response"},
			{"Continue", "continue"},
		})

		if err != nil {
			utils.ShowError("Error in JSON options", err)
			break
		}

		switch choice {
		case "query":
			queryJSON(response)
		case "response":
			displayResponse(response)
		case "continue":
			utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Another Request", "Main Menu")
			return
		}
	}
}

func queryJSON(response *HTTPResponse) {
	query, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter JSON Query",
		Description: "Examples: .name, .users[0].email, .data.items",
		Placeholder: ".field.subfield",
	})

	if err != nil {
		utils.ShowError("Error getting query input", err)
		return
	}

	if strings.TrimSpace(query) == "" {
		utils.ShowMessage("No query provided")
		return
	}

	result := executeJSONQuery(response.ParsedJSON, query)
	displayQueryResult(query, result)
}

func executeJSONQuery(data any, query string) any {
	if data == nil {
		return "null"
	}

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

func handleArrayAccess(data any, accessor string) any {
	parts := strings.Split(accessor, "[")
	fieldName := parts[0]
	indexStr := strings.TrimSuffix(parts[1], "]")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "Invalid array index"
	}

	if fieldName != "" {
		data = handleFieldAccess(data, fieldName)
		if data == nil {
			return nil
		}
	}

	switch arr := data.(type) {
	case []any:
		if index < 0 || index >= len(arr) {
			return "Array index out of bounds"
		}
		return arr[index]
	default:
		return "Not an array"
	}
}

func handleFieldAccess(data any, field string) any {
	switch obj := data.(type) {
	case map[string]any:
		return obj[field]
	case map[any]any:
		return obj[field]
	default:
		return nil
	}
}

func displayQueryResult(query string, result any) {
	var resultStr string

	if result == nil {
		resultStr = "null"
	} else {
		switch v := result.(type) {
		case string:
			resultStr = v
		case map[string]any, []any:
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
	utils.DisplayFormattedText("🔍 Query Result", displayText)
}

// Body Type Selection
func handleBodyTypeSelection(method string) (string, string) {
	var options []utils.SelectionOption

	if method == "PATCH" {
		options = []utils.SelectionOption{
			{"JSON", "json"},
			{"Form Data", "form"},
			{"Raw Text", "raw"},
			{"No Body", "none"},
		}
	} else {
		options = []utils.SelectionOption{
			{"JSON", "json"},
			{"Form Data", "form"},
			{"Multipart Form", "multipart"},
			{"Raw Text", "raw"},
			{"File Upload", "file"},
			{"No Body", "none"},
		}
	}

	bodyType, err := utils.AskSelection("Select Body Type:", options)
	if err != nil {
		utils.ShowError("Error selecting body type", err)
		return "none", ""
	}

	return bodyType, handleBodyInput(bodyType)
}

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
	jsonBody, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter JSON Body:",
		Placeholder: `{"name": "John Doe", "email": "john@example.com"}`,
		Multiline:   true,
	})

	if err != nil {
		utils.ShowError("Error getting JSON input", err)
		return ""
	}

	// Validate JSON
	if strings.TrimSpace(jsonBody) != "" {
		var temp any
		if err := json.Unmarshal([]byte(jsonBody), &temp); err != nil {
			utils.ShowError("Invalid JSON format", err)
			return handleJSONInput() // Retry
		}
	}

	return jsonBody
}

func handleFormDataInput() string {
	formData, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter Form Data:",
		Description: "Format: key1=value1&key2=value2",
		Placeholder: "name=John Doe&email=john@example.com",
		Multiline:   true,
	})

	if err != nil {
		utils.ShowError("Error getting form data", err)
		return ""
	}

	return formData
}

func handleMultipartFormInput() string {
	multipartData, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter Multipart Form Data:",
		Description: "Format: key1=value1&key2=value2",
		Placeholder: "name=John Doe&email=john@example.com&file=@/path/to/file",
		Multiline:   true,
	})

	if err != nil {
		utils.ShowError("Error getting multipart data", err)
		return ""
	}

	return multipartData
}

func handleRawTextInput() string {
	rawText, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter Raw Text:",
		Placeholder: "Enter your raw text content here...",
		Multiline:   true,
	})

	if err != nil {
		utils.ShowError("Error getting raw text", err)
		return ""
	}

	return rawText
}

func handleFileUploadInput() string {
	inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
		{
			Title:       "Enter file path:",
			Placeholder: "/path/to/your/file.jpg",
			Required:    true,
		},
		{
			Title:       "Enter field name:",
			Placeholder: "image",
			Required:    true,
		},
	})

	if err != nil || len(inputs) < 2 {
		utils.ShowError("Error getting file upload info", err)
		return ""
	}

	return fmt.Sprintf("file:%s;field:%s", inputs[0], inputs[1])
}

// Handle Request Options
func handleRequestOptions(method string, endpoint string, body string) hc.RequestOptions {
	options := hc.RequestOptions{
		Method:      method,
		URL:         endpoint,
		Body:        body,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	// Add Headers?
	if addHeaders, _ := utils.AskConfirmation("Add Headers?", "", "", ""); addHeaders {
		options.Headers = collectKeyValuePairs("Header", "Content-Type", "application/json")
	}

	// Add Query Parameters?
	if addParams, _ := utils.AskConfirmation("Add Query Parameters?", "", "", ""); addParams {
		options.QueryParams = collectKeyValuePairs("Parameter", "page", "1")
	}

	// Authentication?
	if addAuth, _ := utils.AskConfirmation("Add Authentication?", "", "", ""); addAuth {
		authType, authValue := handleAuthentication()
		if authType != "" && authValue != "" {
			applyAuthentication(&options, authType, authValue)
		}
	}

	// Add Files/Images? (for POST/PUT only)
	if method == "POST" || method == "PUT" {
		if addFiles, _ := utils.AskConfirmation("Add Files/Images?", "", "", ""); addFiles {
			//options.Files = handleFileUploads()
		}
	}

	// Save as Template?
	//	options.SaveAsTemplate, _ = utils.AskConfirmation("Save as Template?", "", "", "")

	return options
}

// applyAuthentication applies the authentication to the request options
func applyAuthentication(options *hc.RequestOptions, authType, authValue string) {
	switch authType {
	case "bearer":
		// Add Authorization header with Bearer token
		options.Headers["Authorization"] = fmt.Sprintf(" Bearer %s", authValue)

	case "apikey":
		// Parse the header:key format
		parts := strings.SplitN(authValue, ":", 2)
		if len(parts) == 2 {
			headerName := parts[0]
			apiKey := parts[1]
			options.Headers[headerName] = apiKey
		} else {
			// Fallback to default header if parsing fails
			options.Headers["X-API-Key"] = authValue
		}

	case "basic":
		// Parse the username:password format
		parts := strings.SplitN(authValue, ":", 2)
		if len(parts) == 2 {
			username := parts[0]
			password := parts[1]

			// Option 1: Use the BasicAuth struct if your http client supports it
			options.Auth = &hc.BasicAuth{
				Username: username,
				Password: password,
			}

			// Option 2: Or add Authorization header with base64 encoded credentials
			// credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			// options.Headers["Authorization"] = fmt.Sprintf("Basic %s", credentials)
		}

	default:
		utils.ShowWarning(fmt.Sprintf("Unknown authentication type: %s", authType))
	}
}

func collectKeyValuePairs(itemType, keyPlaceholder, valuePlaceholder string) map[string]string {
	items := make(map[string]string)

	for {
		inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
			{
				Title:       fmt.Sprintf("%s Key:", itemType),
				Placeholder: keyPlaceholder,
				Required:    true,
			},
			{
				Title:       fmt.Sprintf("%s Value:", itemType),
				Placeholder: valuePlaceholder,
				Required:    true,
			},
		})

		if err != nil || len(inputs) < 2 || inputs[0] == "" {
			break
		}

		items[strings.TrimSpace(inputs[0])] = strings.TrimSpace(inputs[1])

		addMore, _ := utils.AskConfirmation(fmt.Sprintf("Add another %s?", strings.ToLower(itemType)), "", "", "")
		if !addMore {
			break
		}
	}

	return items
}

func handleAuthentication() (string, string) {
	selectOptions := []utils.SelectionOption{
		{"Bearer Token", "bearer"},
		{"API Key", "apikey"},
		{"Basic Auth", "basic"},
	}
	if ActiveProfile != "" {
		selectOptions = append(selectOptions, utils.SelectionOption{"Use Auth Profile", "profile"})
	}
	authType, err := utils.AskSelection("Select Authentication Type:", selectOptions)

	if err != nil {
		utils.ShowError("Error selecting auth type", err)
		return "", ""
	}

	switch authType {
	case "bearer":
		return handleBearerToken()
	case "apikey":
		return handleAPIKey()
	case "basic":
		return handleBasicAuth()
	case "profile":
		return handleAuthProfile()
	default:
		return "", ""
	}
}

func handleBearerToken() (string, string) {
	token, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter Bearer Token:",
		Placeholder: "your-jwt-token-here",
		Password:    true,
		Required:    true,
	})

	if err != nil {
		utils.ShowError("Error getting bearer token", err)
		return "", ""
	}

	return "bearer", token
}

func handleAPIKey() (string, string) {
	inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
		{
			Title:    "Enter API Key:",
			Password: true,
			Required: true,
		},
		{
			Title:       "Header Name (optional):",
			Placeholder: "X-API-Key",
		},
	})

	if err != nil || len(inputs) < 2 {
		utils.ShowError("Error getting API key", err)
		return "", ""
	}

	headerName := inputs[1]
	if headerName == "" {
		headerName = "X-API-Key"
	}

	return "apikey", fmt.Sprintf("%s:%s", headerName, inputs[0])
}

func handleBasicAuth() (string, string) {
	inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
		{
			Title:    "Username:",
			Required: true,
		},
		{
			Title:    "Password:",
			Password: true,
			Required: true,
		},
	})

	if err != nil || len(inputs) < 2 {
		utils.ShowError("Error getting basic auth", err)
		return "", ""
	}

	return "basic", fmt.Sprintf("%s:%s", inputs[0], inputs[1])
}

func handleAuthProfile() (string, string) {
	if ActiveProfile == "" {
		utils.ShowError("No active profile", fmt.Errorf("no authentication profile is currently active"))
		return "", ""
	}

	authProfile, exists := AuthProfiles[ActiveProfile]
	if !exists {
		utils.ShowError("Profile not found", fmt.Errorf("authentication profile '%s' not found", ActiveProfile))
		return "", ""
	}

	// Check if profile is active
	if !authProfile.Active {
		utils.ShowError("Profile inactive", fmt.Errorf("authentication profile '%s' is not active", ActiveProfile))
		return "", ""
	}

	// Check if profile has expired
	if authProfile.Expiry != nil && time.Now().After(*authProfile.Expiry) {
		utils.ShowError("Profile expired", fmt.Errorf("authentication profile '%s' has expired", ActiveProfile))
		return "", ""
	}

	switch authProfile.Type {
	case "bearer":
		if authProfile.Token == "" {
			utils.ShowError("Invalid profile", fmt.Errorf("bearer token is empty in profile '%s'", ActiveProfile))
			return "", ""
		}
		return "bearer", authProfile.Token

	case "apikey":
		if authProfile.APIKey == "" {
			utils.ShowError("Invalid profile", fmt.Errorf("API key is empty in profile '%s'", ActiveProfile))
			return "", ""
		}
		headerName := authProfile.Header
		if headerName == "" {
			headerName = "X-API-Key" // Default header name
		}
		return "apikey", fmt.Sprintf("%s:%s", headerName, authProfile.APIKey)

	case "basic":
		if authProfile.Username == "" || authProfile.Password == "" {
			utils.ShowError("Invalid profile", fmt.Errorf("username or password is empty in profile '%s'", ActiveProfile))
			return "", ""
		}
		return "basic", fmt.Sprintf("%s:%s", authProfile.Username, authProfile.Password)

	case "oauth":
		if authProfile.Token == "" {
			utils.ShowError("Invalid profile", fmt.Errorf("OAuth token is empty in profile '%s'", ActiveProfile))
			return "", ""
		}
		// OAuth tokens are typically used as bearer tokens
		return "bearer", authProfile.Token

	default:
		utils.ShowError("Unsupported auth type", fmt.Errorf("authentication type '%s' in profile '%s' is not supported", authProfile.Type, ActiveProfile))
		return "", ""
	}
}

func handleFileUploads() []FileUpload {
	var files []FileUpload

	for {
		inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
			{
				Title:       "File Path:",
				Placeholder: "/path/to/file.jpg",
				Required:    true,
			},
			{
				Title:       "Field Name:",
				Placeholder: "image",
				Required:    true,
			},
		})

		if err != nil || len(inputs) < 2 || inputs[0] == "" {
			break
		}

		files = append(files, FileUpload{
			Path:      strings.TrimSpace(inputs[0]),
			FieldName: strings.TrimSpace(inputs[1]),
		})

		addMore, _ := utils.AskConfirmation("Add another file?", "", "", "")
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

	var raw any
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
