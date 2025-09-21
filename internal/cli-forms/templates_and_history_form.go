package cliforms

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/model"
	"github.com/Esa824/apix/internal/utils"
)

func HandleTemplatesAndHistory() {
	var selectedOption string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Templates & History:").Options(
				huh.NewOption("Saved Templates", "saved-templates"),
				huh.NewOption("Request History", "request-history"),
				huh.NewOption("Back to Main Menu", "back"),
			).
				Value(&selectedOption),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running templates & history form: %v\n", err)
		return
	}

	handleTemplatesHistorySelection(selectedOption)
}

func handleTemplatesHistorySelection(selection string) {
	switch selection {
	case "saved-templates":
		handleSavedTemplates()
	case "request-history":
		handleRequestHistory()
	case "back":
		RunInteractiveMode()
	default:
		fmt.Println("Unknown templates & history option")
	}
}

func handleSavedTemplates() {
	var selectedOption string

	templates, err := hc.GetTemplates()

	options := []huh.Option[string]{}
	for _, template := range templates {
		options = append(options, huh.NewOption(fmt.Sprintf("%s (%s)", template.Name, template.Method), template.Name))
	}

	// Add management options
	options = append(options,
		huh.NewOption("Create Templates From Swagger File", "create-templates-from-swagger-file"),
		huh.NewOption("Back", "back"),
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Saved Templates:").
				Description("Select a template to execute, or manage your templates").
				Options(options...).
				Value(&selectedOption),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	handleTemplateSelection(selectedOption)
}

func handleTemplateSelection(selection string) {
	switch selection {
	case "create-templates-from-swagger-file":
		handleCreateTemplatesFromSwaggerFile()
	case "back":
		HandleTemplatesAndHistory()
	default:
		template, err := hc.GetTemplateByName(selection)
		if template != nil || err != nil {
			handleTemplateActions(template)
		} else {
			fmt.Println("Template not found")
			askContinueOrReturnTemplates()
		}
	}
}

func handleTemplateActions(template *model.Template) {
	var selectedAction string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Template: %s", template.Name)).
				Description(fmt.Sprintf("%s %s", template.Method, template.URL)).
				Options(
					huh.NewOption("Execute Template", "execute"),
					huh.NewOption("Edit Template", "edit"),
					huh.NewOption("Delete Template", "delete"),
					huh.NewOption("Back to Templates", "back"),
				).
				Value(&selectedAction),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch selectedAction {
	case "execute":
		executeTemplate(template)
	case "edit":
		editTemplate(template)
	case "delete":
		deleteTemplate(template)
	case "back":
		handleSavedTemplates()
	}
}

func handleCreateTemplatesFromSwaggerFile() {
	var swaggerFilePath, baseURL string
	var selectedEndpoints []string

	// Step 1: Get Swagger file path
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Swagger File Path:").
				Description("Path to your Swagger/OpenAPI JSON or YAML file").
				Placeholder("/path/to/swagger.json or https://api.example.com/swagger.json").
				Value(&swaggerFilePath),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if swaggerFilePath == "" {
		fmt.Println("No file path provided. Returning to templates menu.")
		askContinueOrReturnTemplates()
		return
	}

	// Step 2: Read and parse Swagger file
	var swaggerData map[string]interface{}

	if strings.HasPrefix(swaggerFilePath, "http") {
		// Handle URL
		resp, err := http.Get(swaggerFilePath)
		if err != nil {
			fmt.Printf("Error fetching Swagger file: %v\n", err)
			askContinueOrReturnTemplates()
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			askContinueOrReturnTemplates()
			return
		}

		if strings.Contains(swaggerFilePath, ".yaml") || strings.Contains(swaggerFilePath, ".yml") {
			err = yaml.Unmarshal(body, &swaggerData)
		} else {
			err = json.Unmarshal(body, &swaggerData)
		}
	} else {
		// Handle local file
		fileData, err := os.ReadFile(swaggerFilePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			askContinueOrReturnTemplates()
			return
		}

		if strings.HasSuffix(swaggerFilePath, ".yaml") || strings.HasSuffix(swaggerFilePath, ".yml") {
			err = yaml.Unmarshal(fileData, &swaggerData)
		} else {
			err = json.Unmarshal(fileData, &swaggerData)
		}
	}

	if err != nil {
		fmt.Printf("Error parsing Swagger file: %v\n", err)
		askContinueOrReturnTemplates()
		return
	}

	// Step 3: Extract base URL from Swagger or ask user
	swaggerBaseURL := ""
	if schemes, ok := swaggerData["schemes"].([]interface{}); ok && len(schemes) > 0 {
		if host, ok := swaggerData["host"].(string); ok {
			if basePath, ok := swaggerData["basePath"].(string); ok {
				swaggerBaseURL = fmt.Sprintf("%s://%s%s", schemes[0], host, basePath)
			} else {
				swaggerBaseURL = fmt.Sprintf("%s://%s", schemes[0], host)
			}
		}
	}

	// OpenAPI 3.0 format
	if servers, ok := swaggerData["servers"].([]interface{}); ok && len(servers) > 0 {
		if server, ok := servers[0].(map[string]interface{}); ok {
			if url, ok := server["url"].(string); ok {
				swaggerBaseURL = url
			}
		}
	}

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Base URL:").
				Description("Base URL for API requests").
				Placeholder("https://api.example.com").
				Value(&baseURL).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("base URL is required")
					}
					return nil
				}),
		),
	)

	if swaggerBaseURL != "" {
		baseURL = swaggerBaseURL
	}

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Step 4: Parse paths and create endpoint options
	paths, ok := swaggerData["paths"].(map[string]interface{})
	if !ok {
		fmt.Println("No paths found in Swagger file")
		askContinueOrReturnTemplates()
		return
	}

	var endpointOptions []huh.Option[string]
	var endpointMap = make(map[string]map[string]interface{})

	for path, pathData := range paths {
		if pathMethods, ok := pathData.(map[string]interface{}); ok {
			for method, methodData := range pathMethods {
				method = strings.ToUpper(method)
				if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
					if methodInfo, ok := methodData.(map[string]interface{}); ok {
						summary := ""
						if s, ok := methodInfo["summary"].(string); ok {
							summary = s
						}

						endpointKey := fmt.Sprintf("%s:%s", method, path)
						label := fmt.Sprintf("%s %s", method, path)
						if summary != "" {
							label += fmt.Sprintf(" - %s", summary)
						}

						endpointOptions = append(endpointOptions, huh.NewOption(label, endpointKey))
						endpointMap[endpointKey] = methodInfo
					}
				}
			}
		}
	}

	if len(endpointOptions) == 0 {
		fmt.Println("No valid endpoints found in Swagger file")
		askContinueOrReturnTemplates()
		return
	}

	// Step 5: Let user select endpoints
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Endpoints to Create Templates:").
				Description("Choose which endpoints you want to create templates for").
				Options(endpointOptions...).
				Value(&selectedEndpoints),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(selectedEndpoints) == 0 {
		fmt.Println("No endpoints selected")
		askContinueOrReturnTemplates()
		return
	}

	// Step 6: Parse security definitions
	var globalAuth *string
	var securityDefs map[string]interface{}

	// Swagger 2.0 security definitions
	if swaggerSecDefs, ok := swaggerData["securityDefinitions"].(map[string]interface{}); ok {
		securityDefs = swaggerSecDefs
	}

	// OpenAPI 3.0 security schemes
	if components, ok := swaggerData["components"].(map[string]interface{}); ok {
		if secSchemes, ok := components["securitySchemes"].(map[string]interface{}); ok {
			securityDefs = secSchemes
		}
	}

	// Check for global security requirements
	if security, ok := swaggerData["security"].([]interface{}); ok && len(security) > 0 {
		if secReq, ok := security[0].(map[string]interface{}); ok {
			for secName := range secReq {
				if secDef, ok := securityDefs[secName].(map[string]interface{}); ok {
					authType := parseAuthType(secDef)
					if authType != "" {
						globalAuth = &authType
						break
					}
				}
			}
		}
	}

	// Step 7: Create templates
	createdCount := 0
	for _, endpointKey := range selectedEndpoints {
		parts := strings.SplitN(endpointKey, ":", 2)
		if len(parts) != 2 {
			continue
		}

		method := parts[0]
		path := parts[1]
		methodInfo := endpointMap[endpointKey]

		// Generate template name
		templateName := fmt.Sprintf("%s %s", method, path)
		if summary, ok := methodInfo["summary"].(string); ok && summary != "" {
			templateName = summary
		}

		// Build headers
		headers := make(map[string]string)

		// Check if endpoint consumes JSON
		if consumes, ok := methodInfo["consumes"].([]interface{}); ok {
			for _, consume := range consumes {
				if consumeStr, ok := consume.(string); ok && strings.Contains(consumeStr, "application/json") {
					headers["Content-Type"] = "application/json"
					break
				}
			}
		}

		// OpenAPI 3.0 requestBody check
		if requestBody, ok := methodInfo["requestBody"].(map[string]interface{}); ok {
			if content, ok := requestBody["content"].(map[string]interface{}); ok {
				if _, hasJson := content["application/json"]; hasJson {
					headers["Content-Type"] = "application/json"
				}
			}
		}

		// Determine auth type for this endpoint
		var endpointAuth *string

		// Check endpoint-specific security (standard security definitions)
		if security, ok := methodInfo["security"].([]interface{}); ok && len(security) > 0 {
			if secReq, ok := security[0].(map[string]interface{}); ok {
				for secName := range secReq {
					if secDef, ok := securityDefs[secName].(map[string]interface{}); ok {
						authType := parseAuthType(secDef)
						if authType != "" {
							endpointAuth = &authType
							break
						}
					}
				}
			}
		} else if globalAuth != nil {
			// Use global auth if no endpoint-specific auth
			endpointAuth = globalAuth
		}

		// Check for JWT token in parameters (common pattern in Swagger files without proper security definitions)
		if endpointAuth == nil {
			if parameters, ok := methodInfo["parameters"].([]interface{}); ok {
				for _, param := range parameters {
					if paramMap, ok := param.(map[string]interface{}); ok {
						name, nameOk := paramMap["name"].(string)
						in, inOk := paramMap["in"].(string)

						if nameOk && inOk && in == "header" && strings.ToLower(name) == "authorization" {
							if desc, ok := paramMap["description"].(string); ok {
								if strings.Contains(strings.ToLower(desc), "jwt") ||
									strings.Contains(strings.ToLower(desc), "bearer") {
									bearerAuth := "bearer"
									endpointAuth = &bearerAuth
									break
								}
							}
							// Default to bearer token for Authorization header
							bearerAuth := "bearer"
							endpointAuth = &bearerAuth
							break
						}
					}
				}
			}
		}

		// Generate sample request body for POST/PUT/PATCH
		var sampleBody string
		if method == "POST" || method == "PUT" || method == "PATCH" {
			// Swagger 2.0 parameters
			if parameters, ok := methodInfo["parameters"].([]interface{}); ok {
				for _, param := range parameters {
					if paramMap, ok := param.(map[string]interface{}); ok {
						if in, ok := paramMap["in"].(string); ok && in == "body" {
							if schema, ok := paramMap["schema"].(map[string]interface{}); ok {
								sampleBody = generateSampleBody(schema, swaggerData)
							}
						}
					}
				}
			}

			// OpenAPI 3.0 requestBody
			if requestBody, ok := methodInfo["requestBody"].(map[string]interface{}); ok {
				if content, ok := requestBody["content"].(map[string]interface{}); ok {
					if jsonContent, ok := content["application/json"].(map[string]interface{}); ok {
						if schema, ok := jsonContent["schema"].(map[string]interface{}); ok {
							sampleBody = generateSampleBody(schema, swaggerData)
						}
					}
				}
			}

			// If no body schema found, provide a basic JSON template
			if sampleBody == "" && len(headers) > 0 && headers["Content-Type"] == "application/json" {
				sampleBody = "{}"
			}
		}

		// Create template
		template := model.Template{
			Name:        templateName,
			Method:      method,
			URL:         baseURL + path,
			Headers:     headers,
			QueryParams: make(map[string]string),
			Body:        sampleBody,
		}

		if endpointAuth != nil {
			template.Auth = &model.Auth{Type: *endpointAuth}
		}

		err := hc.SaveTemplate(template)
		if err != nil {
			fmt.Printf("Error saving template '%s': %v\n", templateName, err)
		} else {
			createdCount++
		}
	}

	fmt.Printf("\n✅ Successfully created %d templates from Swagger file!\n", createdCount)
	fmt.Printf("Templates saved and ready to use.\n")

	askContinueOrReturnTemplates()
}

// Helper function to parse auth type from security definition
func parseAuthType(secDef map[string]interface{}) string {
	if authType, ok := secDef["type"].(string); ok {
		switch authType {
		case "apiKey":
			if in, ok := secDef["in"].(string); ok {
				switch in {
				case "header":
					return "apikey"
				case "query":
					return "apikey"
				}
			}
			return "apikey" // default
		case "http":
			if scheme, ok := secDef["scheme"].(string); ok {
				switch scheme {
				case "basic":
					return "basic"
				case "bearer":
					return "bearer"
				}
			}
			return "bearer" // default for http
		case "oauth2":
			return "oauth2"
		}
	}

	// Swagger 2.0 specific
	if authType, ok := secDef["type"].(string); ok {
		switch authType {
		case "basic":
			return "basic"
		case "oauth2":
			return "oauth2"
		}
	}

	return ""
}

// Helper function to generate sample request body from schema
func generateSampleBody(schema map[string]interface{}, swaggerData map[string]interface{}) string {
	sampleData := generateSampleFromSchema(schema, swaggerData, make(map[string]bool))

	if sampleData != nil {
		jsonBytes, err := json.MarshalIndent(sampleData, "", "  ")
		if err == nil {
			return string(jsonBytes)
		}
	}

	return "{}"
}

// Recursive function to generate sample data from schema with reference resolution
func generateSampleFromSchema(schema map[string]interface{}, swaggerData map[string]interface{}, visited map[string]bool) interface{} {
	// Handle $ref
	if ref, ok := schema["$ref"].(string); ok {
		// Prevent infinite recursion
		if visited[ref] {
			return map[string]interface{}{}
		}
		visited[ref] = true

		// Resolve reference
		if refSchema := resolveReference(ref, swaggerData); refSchema != nil {
			result := generateSampleFromSchema(refSchema, swaggerData, visited)
			delete(visited, ref)
			return result
		}
		return map[string]interface{}{}
	}

	// Handle array type
	if schemaType, ok := schema["type"].(string); ok && schemaType == "array" {
		if items, ok := schema["items"].(map[string]interface{}); ok {
			sampleItem := generateSampleFromSchema(items, swaggerData, visited)
			return []interface{}{sampleItem}
		}
		return []interface{}{}
	}

	// Handle object type
	if schemaType, ok := schema["type"].(string); ok && schemaType == "object" {
		sampleData := make(map[string]interface{})

		if properties, ok := schema["properties"].(map[string]interface{}); ok {
			for propName, propSchema := range properties {
				if propMap, ok := propSchema.(map[string]interface{}); ok {
					sampleData[propName] = generateSampleFromSchema(propMap, swaggerData, visited)
				}
			}
		}

		return sampleData
	}

	// Handle primitive types
	if schemaType, ok := schema["type"].(string); ok {
		switch schemaType {
		case "string":
			if example, ok := schema["example"].(string); ok {
				return example
			}
			if format, ok := schema["format"].(string); ok {
				switch format {
				case "date":
					return "2023-01-01"
				case "date-time":
					return "2023-01-01T00:00:00Z"
				case "email":
					return "example@email.com"
				default:
					return "string"
				}
			}
			return "string"
		case "integer":
			if example, ok := schema["example"].(float64); ok {
				return int(example)
			}
			return 0
		case "number":
			if example, ok := schema["example"].(float64); ok {
				return example
			}
			return 0.0
		case "boolean":
			if example, ok := schema["example"].(bool); ok {
				return example
			}
			return true
		}
	}

	return map[string]interface{}{}
}

// Helper function to resolve $ref references in Swagger schema
func resolveReference(ref string, swaggerData map[string]interface{}) map[string]interface{} {
	// Handle #/definitions/... references
	if strings.HasPrefix(ref, "#/definitions/") {
		definitionName := strings.TrimPrefix(ref, "#/definitions/")
		if definitions, ok := swaggerData["definitions"].(map[string]interface{}); ok {
			if definition, ok := definitions[definitionName].(map[string]interface{}); ok {
				return definition
			}
		}
	}

	// Handle #/components/schemas/... references (OpenAPI 3.0)
	if strings.HasPrefix(ref, "#/components/schemas/") {
		schemaName := strings.TrimPrefix(ref, "#/components/schemas/")
		if components, ok := swaggerData["components"].(map[string]interface{}); ok {
			if schemas, ok := components["schemas"].(map[string]interface{}); ok {
				if schema, ok := schemas[schemaName].(map[string]interface{}); ok {
					return schema
				}
			}
		}
	}

	return nil
}

func handleRequestHistory() {
	var selectedOption string

	// TODO: Replace with actual request history from storage
	history, err := hc.GetHistory()

	options := []huh.Option[string]{}
	for _, request := range history {
		label := fmt.Sprintf("%s %s - %s", request.Method, request.URL, utils.FormatTime(request.Time))
		options = append(options, huh.NewOption(label, strconv.Itoa(request.Id)))
	}

	// Add management options
	options = append(options,
		huh.NewOption("Clear History", "clear-history"),
		huh.NewOption("Back", "back"),
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Request History:").
				Description("Select a request to re-execute or save as template").
				Options(options...).
				Value(&selectedOption),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	handleHistorySelection(selectedOption)
}

func handleHistorySelection(selection string) {
	switch selection {
	case "clear-history":
		handleClearHistory()
	case "back":
		HandleTemplatesAndHistory()
	default:
		id, err := strconv.Atoi(selection)
		if err != nil {
			return
		}
		history, err := hc.GetHistory()
		if err != nil {
			return
		}
		handleHistoryActions(&history[id])
		return
	}
}

func handleHistoryActions(historyItem *hc.RequestOptions) {
	var selectedAction string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Request: %s %s", historyItem.Method, historyItem.URL)).
				Description(fmt.Sprintf("Executed: %s", utils.FormatTime(historyItem.Time))).
				Options(
					huh.NewOption("Re-execute Request", "reexecute"),
					huh.NewOption("Save as Template", "save-template"),
					huh.NewOption("View Details", "view-details"),
					huh.NewOption("Back to History", "back"),
				).
				Value(&selectedAction),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch selectedAction {
	case "reexecute":
		reExecuteFromHistory(historyItem)
	case "save-template":
		saveHistoryAsTemplate(historyItem)
	case "view-details":
		viewHistoryDetails(historyItem)
	case "back":
		handleRequestHistory()
	}
}

func handleClearHistory() {
	var confirmClear bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Clear Request History").
				Description("This will permanently delete all request history. This action cannot be undone.").
				Affirmative("Yes, clear history").
				Negative("Cancel").
				Value(&confirmClear),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if confirmClear {
		hc.DeleteHistory()
		fmt.Println("Clearing request history...")
		fmt.Println("Request history cleared successfully!")
	} else {
		fmt.Println("History clearing cancelled.")
	}

	askContinueOrReturnTemplates()
}

func executeTemplate(template *model.Template) {
	// Get endpoint input
	input, err := utils.AskInput(utils.InputConfig{
		Title:       "Enter API endpoint:",
		Description: "Will be appended to base URL if configured",
		Placeholder: "/api/users or https://api.example.com/users",
		Value:       template.URL,
		Required:    true,
	})
	if err != nil || input == "" {
		utils.ShowMessage("No endpoint provided. Returning to menu.")
		utils.AskContinueOrReturn(HandleHttpRequests, RunInteractiveMode, "Try Again", "Main Menu")
		return
	}

	var body string
	if template.Body != nil && template.Body != "" {
		_, body = handleBodyTypeSelection(template.Method, template.Body)
	}
	headers := make(map[string]string)
	if template.Headers != nil {
		headers = utils.CollectKeyValuePairs("Header", "Content-Type", "application/json", template.Headers)
	}
	queryParams := make(map[string]string)
	if template.QueryParams != nil {
		queryParams = utils.CollectKeyValuePairs("Parameter", "page", "1", template.QueryParams)
	}

	// Setup request options
	options := hc.RequestOptions{
		Method:      template.Method,
		URL:         input,
		Body:        body,
		Headers:     headers,
		QueryParams: queryParams,
		Time:        time.Now(),
	}

	// Handle authentication
	if template.Auth != nil {
		if authType, authValue := handleAuthentication(*template.Auth); authType != "" && authValue != "" {
			applyAuthentication(&options, authType, authValue)
		}
	}

	// Handle file uploads
	if template.Files != nil {
		options.Files = handleFileUploads(template.Files)
	}

	// Execute request and handle response
	response, _ := hc.NewClient(10*time.Second).Do(options, false)
	utils.HandleResponse(response, HandleTemplatesAndHistory, RunInteractiveMode, "Continue with templates & history", "Return to Main Menu")
}

func editTemplate(template *model.Template) {
	// Create a copy of the template to edit
	editedTemplate := *template

	// Step 1: Edit basic details
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Template Name:").
				Description("Enter a descriptive name for your template").
				Placeholder("My API Template").
				Value(&editedTemplate.Name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("template name is required")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("HTTP Method:").
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
					huh.NewOption("PATCH", "PATCH"),
				).
				Value(&editedTemplate.Method),
			huh.NewInput().
				Title("URL:").
				Description("Full URL or endpoint path").
				Placeholder("https://api.example.com/users or /api/users").
				Value(&editedTemplate.URL).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("URL is required")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Step 2: Edit headers
	var editHeaders bool
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Edit Headers?").
				Description("Would you like to modify the headers for this template?").
				Affirmative("Yes").
				Negative("No").
				Value(&editHeaders),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if editHeaders {
		editedTemplate.Headers = utils.CollectKeyValuePairs("Header", "Content-Type", "application/json", editedTemplate.Headers)
	}

	// Step 3: Edit query parameters
	var editParams bool
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Edit Query Parameters?").
				Description("Would you like to modify the query parameters for this template?").
				Affirmative("Yes").
				Negative("No").
				Value(&editParams),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if editParams {
		editedTemplate.QueryParams = utils.CollectKeyValuePairs("Parameter", "page", "1", editedTemplate.QueryParams)
	}

	// Step 4: Edit body (for POST, PUT, PATCH)
	if editedTemplate.Method == "POST" || editedTemplate.Method == "PUT" || editedTemplate.Method == "PATCH" {
		var editBody bool
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Edit Request Body?").
					Description("Would you like to modify the request body for this template?").
					Affirmative("Yes").
					Negative("No").
					Value(&editBody),
			),
		)

		err = form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if editBody {
			currentBody := ""
			if editedTemplate.Body != nil {
				currentBody = editedTemplate.Body.(string)
			}
			_, newBody := handleBodyTypeSelection(editedTemplate.Method, currentBody)
			editedTemplate.Body = newBody
		}
	}

	// Step 5: Edit authentication
	var editAuth bool
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Edit Authentication?").
				Description("Would you like to modify the authentication settings for this template?").
				Affirmative("Yes").
				Negative("No").
				Value(&editAuth),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if editAuth {
		var authType string
		authOptions := []huh.Option[string]{
			huh.NewOption("No Authentication", "none"),
			huh.NewOption("Bearer Token", "bearer"),
			huh.NewOption("API Key", "apikey"),
			huh.NewOption("Basic Authentication", "basic"),
		}

		// Set current auth type as default
		currentAuthType := "none"
		if editedTemplate.Auth != nil {
			currentAuthType = editedTemplate.Auth.Type
		}

		form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Authentication Type:").
					Options(authOptions...).
					Value(&authType),
			),
		)

		// Pre-select current auth type
		authType = currentAuthType

		err = form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if authType == "none" {
			editedTemplate.Auth = nil
		} else {
			// Get current auth details if they exist
			currentAuth := model.Auth{Type: authType}
			if editedTemplate.Auth != nil && editedTemplate.Auth.Type == authType {
				currentAuth = *editedTemplate.Auth
			}

			// Use existing authentication handler with current values
			if newAuthType, newAuthValue := handleAuthentication(currentAuth); newAuthType != "" {
				editedTemplate.Auth = &model.Auth{
					Type:    newAuthType,
					Primary: newAuthValue,
				}
			}
		}
	}

	// Step 6: Confirm changes
	var confirmSave bool
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Save Changes?").
				Description(fmt.Sprintf("Save changes to template '%s'?", editedTemplate.Name)).
				Affirmative("Yes, save changes").
				Negative("Cancel").
				Value(&confirmSave),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if confirmSave {
		err = hc.UpdateTemplate(editedTemplate, template.Name)
		if err != nil {
			fmt.Printf("Error updating template: %v\n", err)
		} else {
			fmt.Printf("✅ Template '%s' updated successfully!\n", editedTemplate.Name)
		}
	} else {
		fmt.Println("Changes cancelled.")
	}

	askContinueOrReturnTemplates()
}

func deleteTemplate(template *model.Template) {
	var confirmDelete bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Delete Template: %s", template.Name)).
				Description(fmt.Sprintf("This will permanently delete the template '%s' (%s %s). This action cannot be undone.", template.Name, template.Method, template.URL)).
				Affirmative("Yes, delete template").
				Negative("Cancel").
				Value(&confirmDelete),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if confirmDelete {
		hc.DeleteTemplate(template.Name)
	}

	askContinueOrReturnTemplates()
}

func reExecuteFromHistory(historyItem *hc.RequestOptions) {
	fmt.Printf("Re-executing request: %s %s\n", historyItem.Method, historyItem.URL)
	response, _ := hc.NewClient(10*time.Second).Do(*historyItem, false)
	utils.HandleResponse(response, HandleTemplatesAndHistory, RunInteractiveMode, "Continue with templates & history", "Return to Main Menu")
}

func saveHistoryAsTemplate(opts *hc.RequestOptions) {
	templateName, err := utils.AskInput(utils.InputConfig{
		Title:       "Name",
		Description: "Enter a name for your template",
		Placeholder: "Get products",
		Required:    true,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	hc.SaveTemplate(model.Template{
		Name:        templateName,
		Method:      opts.Method,
		URL:         opts.URL,
		Headers:     opts.Headers,
		QueryParams: opts.QueryParams,
		Files:       opts.Files,
		Auth:        opts.Auth,
		Body:        opts.Body,
	})
	askContinueOrReturnTemplates()
}

func viewHistoryDetails(historyItem *hc.RequestOptions) {
	fmt.Println("📋 Request Details:")
	fmt.Println("─────────────────────────")
	fmt.Printf("Method: %s\n", historyItem.Method)
	fmt.Printf("URL: %s\n", historyItem.URL)
	fmt.Printf("Timestamp: %s\n", "")
	fmt.Printf("Status: %s\n", "")
	if historyItem.Body != "" {
		fmt.Printf("Body: %s\n", historyItem.Body)
	}
	if len(historyItem.Headers) == 0 {
		fmt.Printf("Headers: %s\n", historyItem.Headers)
	}
	fmt.Println("─────────────────────────")

	askContinueOrReturnTemplates()
}

func askContinueOrReturnTemplates() {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
				Options(
					huh.NewOption("Continue with Templates & History", "continue"),
					huh.NewOption("Return to Main Menu", "main"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&choice),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch choice {
	case "continue":
		HandleTemplatesAndHistory()
	case "main":
		RunInteractiveMode()
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
}
