package cliforms

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/utils"
)

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
	response, err := hc.NewClient(10*time.Second).Do(options, true)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	utils.HandleResponse(response, HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
}

// POST Request Handler
func handlePostRequest() {
	endpoint, body := getEndpointAndBody("POST")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("POST", endpoint, body)
	response, err := hc.NewClient(10*time.Second).Do(options, true)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	utils.HandleResponse(response, HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
}

// PUT Request Handler
func handlePutRequest() {
	endpoint, body := getEndpointAndBody("PUT")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("PUT", endpoint, body)
	response, err := hc.NewClient(10*time.Second).Do(options, true)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	utils.HandleResponse(response, HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
}

// PATCH Request Handler
func handlePatchRequest() {
	endpoint, body := getEndpointAndBody("PATCH")
	if endpoint == "" {
		return
	}

	options := handleRequestOptions("PATCH", endpoint, body)
	response, err := hc.NewClient(10*time.Second).Do(options, true)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	utils.HandleResponse(response, HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
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
	response, err := hc.NewClient(10*time.Second).Do(options, true)
	if err != nil {
		utils.ShowError("Error while calling endpoint", err)
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	utils.HandleResponse(response, HandleHttpRequests, RunInteractiveMode, "Another Request", "Main Menu")
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

// Body Type Selection with optional existing body
func handleBodyTypeSelection(method string, existingBody ...any) (string, string) {
	// If existing body is provided, determine its type and handle it
	if len(existingBody) > 0 && existingBody[0] != nil {
		bodyType := determineBodyType(existingBody[0])
		return bodyType, handleBodyInput(bodyType, existingBody[0])
	}

	// Original logic for new body creation
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
	return bodyType, handleBodyInput(bodyType, nil)
}

// Determine the type of existing body
func determineBodyType(body interface{}) string {
	switch v := body.(type) {
	case string:
		// Try to parse as JSON first
		var jsonTest interface{}
		if json.Unmarshal([]byte(v), &jsonTest) == nil {
			return "json"
		}
		// Check if it's URL encoded form data
		if _, err := url.ParseQuery(v); err == nil && strings.Contains(v, "=") {
			return "form"
		}
		// Check if it's multipart form (contains boundary or file references)
		if strings.Contains(v, "boundary=") || strings.Contains(v, "@/") {
			return "multipart"
		}
		// Check if it looks like file upload format
		if strings.HasPrefix(v, "file:") && strings.Contains(v, ";field:") {
			return "file"
		}
		// Default to raw text
		return "raw"
	case map[string]interface{}, []interface{}:
		return "json"
	case url.Values:
		return "form"
	case *multipart.Form:
		return "multipart"
	default:
		// Try to marshal to JSON to see if it's JSON-serializable
		if _, err := json.Marshal(v); err == nil {
			return "json"
		}
		return "raw"
	}
}

// Enhanced body input handler
func handleBodyInput(bodyType string, existingBody interface{}) string {
	switch bodyType {
	case "json":
		return handleJSONInput(existingBody)
	case "form":
		return handleFormDataInput(existingBody)
	case "multipart":
		return handleMultipartFormInput(existingBody)
	case "raw":
		return handleRawTextInput(existingBody)
	case "file":
		return handleFileUploadInput(existingBody)
	default:
		return ""
	}
}

// Enhanced JSON input handler
func handleJSONInput(existingBody interface{}) string {
	var currentJSON string
	var parsedData map[string]interface{}

	// Parse existing body if provided
	if existingBody != nil {
		switch v := existingBody.(type) {
		case string:
			if err := json.Unmarshal([]byte(v), &parsedData); err != nil {
				utils.ShowError("Error parsing existing JSON", err)
				return handleJSONInputFromScratch()
			}
			currentJSON = v
		case map[string]interface{}:
			parsedData = v
			jsonBytes, _ := json.MarshalIndent(v, "", "  ")
			currentJSON = string(jsonBytes)
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return handleJSONInputFromScratch()
			}
			if err := json.Unmarshal(jsonBytes, &parsedData); err != nil {
				return handleJSONInputFromScratch()
			}
			currentJSON = string(jsonBytes)
		}

		// Show current JSON and allow field-by-field editing
		fmt.Printf("Current JSON:\n%s\n\n", currentJSON)

		// Ask if user wants to edit individual fields
		editChoice, err := utils.AskSelection("How would you like to edit?", []utils.SelectionOption{
			{"Edit individual fields", "fields"},
			{"Replace entire JSON", "replace"},
			{"Keep as is", "keep"},
		})
		if err != nil {
			return currentJSON
		}

		switch editChoice {
		case "fields":
			return editJSONFields(parsedData)
		case "replace":
			return handleJSONInputFromScratch()
		case "keep":
			return currentJSON
		}
	}

	return handleJSONInputFromScratch()
}

// Edit JSON fields individually
func editJSONFields(data map[string]interface{}) string {
	editedData := make(map[string]interface{})

	// Copy original data
	for k, v := range data {
		editedData[k] = v
	}

	for key, value := range data {
		currentValueStr := fmt.Sprintf("%v", value)

		newValue, err := utils.AskInput(utils.InputConfig{
			Title:       fmt.Sprintf("Edit field '%s':", key),
			Description: fmt.Sprintf("Current value: %s", currentValueStr),
			Placeholder: currentValueStr,
			Value:       currentValueStr,
		})
		if err != nil {
			continue // Skip this field if error
		}

		// Try to preserve the original type
		editedData[key] = convertToOriginalType(value, newValue)
	}

	// Marshal back to JSON
	result, err := json.MarshalIndent(editedData, "", "  ")
	if err != nil {
		utils.ShowError("Error creating edited JSON", err)
		return ""
	}

	return string(result)
}

// Convert new value to match original type
func convertToOriginalType(original interface{}, newValue string) interface{} {
	if newValue == "" {
		return original
	}

	switch original.(type) {
	case bool:
		if strings.ToLower(newValue) == "true" {
			return true
		} else if strings.ToLower(newValue) == "false" {
			return false
		}
		return original
	case float64:
		if f, err := parseFloat(newValue); err == nil {
			return f
		}
		return original
	case int:
		if i, err := parseInt(newValue); err == nil {
			return i
		}
		return original
	default:
		return newValue
	}
}

// Original JSON input from scratch
func handleJSONInputFromScratch() string {
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
			return handleJSONInputFromScratch() // Retry
		}
	}
	return jsonBody
}

// Enhanced Form Data input handler
func handleFormDataInput(existingBody interface{}) string {
	if existingBody != nil {
		var formData string

		switch v := existingBody.(type) {
		case string:
			formData = v
		case url.Values:
			formData = v.Encode()
		default:
			return handleFormDataInputFromScratch()
		}

		// Parse existing form data
		values, err := url.ParseQuery(formData)
		if err != nil {
			return handleFormDataInputFromScratch()
		}

		fmt.Printf("Current Form Data: %s\n\n", formData)

		// Ask if user wants to edit individual fields
		editChoice, err := utils.AskSelection("How would you like to edit?", []utils.SelectionOption{
			{"Edit individual fields", "fields"},
			{"Replace entire form data", "replace"},
			{"Keep as is", "keep"},
		})
		if err != nil {
			return formData
		}

		switch editChoice {
		case "fields":
			return editFormDataFields(values)
		case "replace":
			return handleFormDataInputFromScratch()
		case "keep":
			return formData
		}
	}

	return handleFormDataInputFromScratch()
}

// Edit form data fields individually
func editFormDataFields(values url.Values) string {
	editedValues := make(url.Values)

	for key, valueList := range values {
		currentValue := ""
		if len(valueList) > 0 {
			currentValue = valueList[0] // Take first value for simplicity
		}

		newValue, err := utils.AskInput(utils.InputConfig{
			Title:       fmt.Sprintf("Edit field '%s':", key),
			Description: fmt.Sprintf("Current value: %s", currentValue),
			Placeholder: currentValue,
			Value:       currentValue,
		})
		if err != nil {
			editedValues[key] = valueList // Keep original if error
			continue
		}

		editedValues[key] = []string{newValue}
	}

	return editedValues.Encode()
}

// Original form data input from scratch
func handleFormDataInputFromScratch() string {
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

// Enhanced Multipart Form input handler
func handleMultipartFormInput(existingBody interface{}) string {
	if existingBody != nil {
		var multipartData string

		switch v := existingBody.(type) {
		case string:
			multipartData = v
		default:
			return handleMultipartFormInputFromScratch()
		}

		fmt.Printf("Current Multipart Data: %s\n\n", multipartData)

		// Parse multipart data (assuming key=value&key2=value2 format)
		values, err := url.ParseQuery(multipartData)
		if err != nil {
			return handleMultipartFormInputFromScratch()
		}

		// Ask if user wants to edit individual fields
		editChoice, err := utils.AskSelection("How would you like to edit?", []utils.SelectionOption{
			{"Edit individual fields", "fields"},
			{"Replace entire multipart data", "replace"},
			{"Keep as is", "keep"},
		})
		if err != nil {
			return multipartData
		}

		switch editChoice {
		case "fields":
			return editMultipartFields(values)
		case "replace":
			return handleMultipartFormInputFromScratch()
		case "keep":
			return multipartData
		}
	}

	return handleMultipartFormInputFromScratch()
}

// Edit multipart fields individually
func editMultipartFields(values url.Values) string {
	editedValues := make(url.Values)

	for key, valueList := range values {
		currentValue := ""
		if len(valueList) > 0 {
			currentValue = valueList[0]
		}

		// Check if it's a file upload field
		if strings.HasPrefix(currentValue, "@/") {
			newValue, err := utils.AskInput(utils.InputConfig{
				Title:       fmt.Sprintf("Edit file path for field '%s':", key),
				Description: fmt.Sprintf("Current path: %s", currentValue),
				Placeholder: currentValue,
				Value:       currentValue,
			})
			if err != nil {
				editedValues[key] = valueList
				continue
			}
			editedValues[key] = []string{newValue}
		} else {
			newValue, err := utils.AskInput(utils.InputConfig{
				Title:       fmt.Sprintf("Edit field '%s':", key),
				Description: fmt.Sprintf("Current value: %s", currentValue),
				Placeholder: currentValue,
				Value:       currentValue,
			})
			if err != nil {
				editedValues[key] = valueList
				continue
			}
			editedValues[key] = []string{newValue}
		}
	}

	return editedValues.Encode()
}

// Original multipart form input from scratch
func handleMultipartFormInputFromScratch() string {
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

// Enhanced Raw Text input handler
func handleRawTextInput(existingBody interface{}) string {
	if existingBody != nil {
		var currentText string

		switch v := existingBody.(type) {
		case string:
			currentText = v
		default:
			currentText = fmt.Sprintf("%v", v)
		}

		fmt.Printf("Current Raw Text:\n%s\n\n", currentText)

		// Ask if user wants to edit or keep
		editChoice, err := utils.AskSelection("How would you like to edit?", []utils.SelectionOption{
			{"Edit text", "edit"},
			{"Replace entirely", "replace"},
			{"Keep as is", "keep"},
		})
		if err != nil {
			return currentText
		}

		switch editChoice {
		case "edit", "replace":
			newText, err := utils.AskInput(utils.InputConfig{
				Title:       "Edit Raw Text:",
				Placeholder: "Enter your raw text content here...",
				Value:       currentText,
				Multiline:   true,
			})
			if err != nil {
				return currentText
			}
			return newText
		case "keep":
			return currentText
		}
	}

	return handleRawTextInputFromScratch()
}

// Original raw text input from scratch
func handleRawTextInputFromScratch() string {
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

// Enhanced File Upload input handler
func handleFileUploadInput(existingBody interface{}) string {
	if existingBody != nil {
		var currentFileData string

		switch v := existingBody.(type) {
		case string:
			currentFileData = v
		default:
			return handleFileUploadInputFromScratch()
		}

		// Parse existing file upload data (format: "file:path;field:name")
		if strings.HasPrefix(currentFileData, "file:") && strings.Contains(currentFileData, ";field:") {
			parts := strings.Split(currentFileData, ";field:")
			if len(parts) == 2 {
				currentPath := strings.TrimPrefix(parts[0], "file:")
				currentField := parts[1]

				fmt.Printf("Current File Upload:\nPath: %s\nField: %s\n\n", currentPath, currentField)

				// Ask if user wants to edit individual components
				editChoice, err := utils.AskSelection("How would you like to edit?", []utils.SelectionOption{
					{"Edit file path", "path"},
					{"Edit field name", "field"},
					{"Edit both", "both"},
					{"Replace entirely", "replace"},
					{"Keep as is", "keep"},
				})
				if err != nil {
					return currentFileData
				}

				switch editChoice {
				case "path":
					newPath, err := utils.AskInput(utils.InputConfig{
						Title:       "Edit file path:",
						Placeholder: "/path/to/your/file.jpg",
						Value:       currentPath,
						Required:    true,
					})
					if err != nil {
						return currentFileData
					}
					return fmt.Sprintf("file:%s;field:%s", newPath, currentField)

				case "field":
					newField, err := utils.AskInput(utils.InputConfig{
						Title:       "Edit field name:",
						Placeholder: "image",
						Value:       currentField,
						Required:    true,
					})
					if err != nil {
						return currentFileData
					}
					return fmt.Sprintf("file:%s;field:%s", currentPath, newField)

				case "both":
					inputs, err := utils.AskMultipleInputs([]utils.InputConfig{
						{
							Title:       "Edit file path:",
							Placeholder: "/path/to/your/file.jpg",
							Value:       currentPath,
							Required:    true,
						},
						{
							Title:       "Edit field name:",
							Placeholder: "image",
							Value:       currentField,
							Required:    true,
						},
					})
					if err != nil || len(inputs) < 2 {
						return currentFileData
					}
					return fmt.Sprintf("file:%s;field:%s", inputs[0], inputs[1])

				case "replace":
					return handleFileUploadInputFromScratch()

				case "keep":
					return currentFileData
				}
			}
		}
	}

	return handleFileUploadInputFromScratch()
}

// Original file upload input from scratch
func handleFileUploadInputFromScratch() string {
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

// Helper functions for type conversion
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// Handle Request Options
func handleRequestOptions(method string, endpoint string, body string) hc.RequestOptions {
	options := hc.RequestOptions{
		Method:      method,
		URL:         endpoint,
		Body:        body,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Time:        time.Now(),
	}

	// Add Headers?
	if addHeaders, _ := utils.AskConfirmation("Add Headers?", "", "", ""); addHeaders {
		options.Headers = utils.CollectKeyValuePairs("Header", "Content-Type", "application/json")
	}

	// Add Query Parameters?
	if addParams, _ := utils.AskConfirmation("Add Query Parameters?", "", "", ""); addParams {
		options.QueryParams = utils.CollectKeyValuePairs("Parameter", "page", "1")
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
			options.Files = handleFileUploads()
		}
	}

	options.IsTemplate, _ = utils.AskConfirmation("Save as Template?", "", "", "")
	if options.IsTemplate {
		options.Name, _ = utils.AskInput(utils.InputConfig{
			Title:       "Name",
			Description: "Enter a name for your template",
			Placeholder: "Get products",
			Required:    true,
		})
	}

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

func handleFileUploads() map[string]string {
	var files map[string]string

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

		files[strings.TrimSpace(inputs[1])] = strings.TrimSpace(inputs[0])

		addMore, _ := utils.AskConfirmation("Add another file?", "", "", "")
		if !addMore {
			break
		}
	}

	return files
}
