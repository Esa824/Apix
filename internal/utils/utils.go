package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/Esa824/apix/internal/model"
)

// =============================================================================
// MESSAGING UTILITIES - Consistent UI messaging across all forms
// =============================================================================

// ShowMessage displays an informational message using huh forms
func ShowMessage(message string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Information").
				Description(message),
		),
	)
	form.Run()
}

// ShowSuccess displays a success message with checkmark
func ShowSuccess(message string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("✅ Success").
				Description(message),
		),
	)
	form.Run()
}

// ShowError displays an error message with X icon
func ShowError(title string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("❌ " + title).
				Description(err.Error()),
		),
	)
	form.Run()
}

// ShowWarning displays a warning message with warning icon
func ShowWarning(message string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("⚠️ Warning").
				Description(message),
		),
	)
	form.Run()
}

// ShowCustomMessage displays a message with custom title and icon
func ShowCustomMessage(title, message, icon string) {
	displayTitle := title
	if icon != "" {
		displayTitle = icon + " " + title
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(displayTitle).
				Description(message),
		),
	)
	form.Run()
}

// =============================================================================
// NAVIGATION UTILITIES - Common navigation patterns
// =============================================================================

// NavigationChoice represents a navigation option
type NavigationChoice struct {
	Label string
	Value string
}

// AskNavigation presents navigation choices to the user
func AskNavigation(title string, choices []NavigationChoice) (string, error) {
	var choice string

	options := make([]huh.Option[string], len(choices))
	for i, c := range choices {
		options[i] = huh.NewOption(c.Label, c.Value)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Value(&choice),
		),
	)

	err := form.Run()
	return choice, err
}

// AskContinueOrReturn provides standard continue/return/exit navigation
func AskContinueOrReturn(continueAction, returnAction func(), continueLabel, returnLabel string) {
	if continueLabel == "" {
		continueLabel = "Continue"
	}
	if returnLabel == "" {
		returnLabel = "Return to Main Menu"
	}

	choices := []NavigationChoice{
		{continueLabel, "continue"},
		{returnLabel, "return"},
		{"Exit", "exit"},
	}

	choice, err := AskNavigation("What would you like to do next?", choices)
	if err != nil {
		ShowError("Navigation error", err)
		return
	}

	switch choice {
	case "continue":
		if continueAction != nil {
			continueAction()
		}
	case "return":
		if returnAction != nil {
			returnAction()
		}
	case "exit":
		ShowMessage("Goodbye!")
		os.Exit(0)
	}
}

// =============================================================================
// CONFIRMATION UTILITIES - Reusable confirmation dialogs
// =============================================================================

// AskConfirmation presents a yes/no confirmation dialog
func AskConfirmation(title, description, affirmative, negative string) (bool, error) {
	if affirmative == "" {
		affirmative = "Yes"
	}
	if negative == "" {
		negative = "No"
	}

	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Affirmative(affirmative).
				Negative(negative).
				Value(&confirmed),
		),
	)

	err := form.Run()
	return confirmed, err
}

// AskDangerousConfirmation presents a confirmation for destructive actions
func AskDangerousConfirmation(title, description, itemName string) (bool, error) {
	affirmative := fmt.Sprintf("Yes, %s", strings.ToLower(title))
	return AskConfirmation(
		title,
		fmt.Sprintf("%s '%s'? This action cannot be undone.", description, itemName),
		affirmative,
		"Cancel",
	)
}

// =============================================================================
// INPUT UTILITIES - Common input patterns
// =============================================================================

// InputConfig configures an input field
type InputConfig struct {
	Title       string
	Description string
	Placeholder string
	Required    bool
	Password    bool
	Multiline   bool
}

// AskInput presents a single input field with configuration
func AskInput(config InputConfig) (string, error) {
	var value string
	var input huh.Field

	if config.Multiline {
		textInput := huh.NewText().
			Title(config.Title).
			Placeholder(config.Placeholder).
			Value(&value)

		if config.Description != "" {
			textInput = textInput.Description(config.Description)
		}
		input = textInput
	} else {
		regularInput := huh.NewInput().
			Title(config.Title).
			Placeholder(config.Placeholder).
			Password(config.Password).
			Value(&value)

		if config.Description != "" {
			regularInput = regularInput.Description(config.Description)
		}
		input = regularInput
	}

	form := huh.NewForm(huh.NewGroup(input))
	err := form.Run()

	if err != nil {
		return "", err
	}

	if config.Required && value == "" {
		return "", fmt.Errorf("this field is required")
	}

	return value, nil
}

// AskMultipleInputs presents multiple input fields at once
func AskMultipleInputs(configs []InputConfig) ([]string, error) {
	values := make([]string, len(configs))
	fields := make([]huh.Field, len(configs))

	for i, config := range configs {
		if config.Multiline {
			textInput := huh.NewText().
				Title(config.Title).
				Placeholder(config.Placeholder).
				Value(&values[i])

			if config.Description != "" {
				textInput = textInput.Description(config.Description)
			}
			fields[i] = textInput
		} else {
			regularInput := huh.NewInput().
				Title(config.Title).
				Placeholder(config.Placeholder).
				Password(config.Password).
				Value(&values[i])

			if config.Description != "" {
				regularInput = regularInput.Description(config.Description)
			}
			fields[i] = regularInput
		}
	}

	form := huh.NewForm(huh.NewGroup(fields...))
	err := form.Run()

	if err != nil {
		return nil, err
	}

	// Check required fields
	for i, config := range configs {
		if config.Required && values[i] == "" {
			return nil, fmt.Errorf("'%s' is required", config.Title)
		}
	}

	return values, nil
}

// =============================================================================
// SELECTION UTILITIES - Common selection patterns
// =============================================================================

// SelectionOption represents an option in a selection menu
type SelectionOption struct {
	Label string
	Value string
}

// AskSelection presents a selection menu
func AskSelection(title string, options []SelectionOption) (string, error) {
	var selected string

	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Label, opt.Value)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(huhOptions...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// AskMultiSelection presents a multi-selection menu
func AskMultiSelection(title string, options []SelectionOption) ([]string, error) {
	var selected []string

	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Label, opt.Value)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Options(huhOptions...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// =============================================================================
// VALIDATION UTILITIES - Input validation helpers
// =============================================================================

// ValidateURL checks if a string is a valid URL format
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	return nil
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateNotEmpty checks if a string is not empty
func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateLength checks if a string meets length requirements
func ValidateLength(value string, minLength, maxLength int, fieldName string) error {
	length := len(value)
	if length < minLength {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, minLength)
	}
	if maxLength > 0 && length > maxLength {
		return fmt.Errorf("%s cannot be longer than %d characters", fieldName, maxLength)
	}
	return nil
}

// =============================================================================
// DISPLAY UTILITIES - Common display helpers
// =============================================================================

// MaskSensitive masks sensitive information (tokens, passwords, etc.)
func MaskSensitive(value string) string {
	if value == "" {
		return "Not set"
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}

// FormatKeyValue formats key-value pairs for display
func FormatKeyValue(key, value string, maskValue bool) string {
	if maskValue {
		value = MaskSensitive(value)
	}
	return fmt.Sprintf("   %s: %s", key, value)
}

// BuildDisplayText creates formatted text for display
func BuildDisplayText(title string, items map[string]string, maskSensitive map[string]bool) string {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("📋 %s:\n", title))
	text.WriteString("─────────────────────────────────\n\n")

	for key, value := range items {
		masked := maskSensitive != nil && maskSensitive[key]
		text.WriteString(FormatKeyValue(key, value, masked))
		text.WriteString("\n")
	}

	text.WriteString("─────────────────────────────────")
	return text.String()
}

// DisplayFormattedText shows formatted text using huh Note
func DisplayFormattedText(title, content string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(title).
				Description(content),
		),
	)
	form.Run()
}

// =============================================================================
// TIME UTILITIES - Time input and formatting helpers
// =============================================================================

// ParseTimeInput parses various time input formats
func ParseTimeInput(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}

	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02",
		"15:04",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, input); err == nil {
			return &parsed, nil
		}
	}

	return nil, fmt.Errorf("invalid time format. Use YYYY-MM-DD HH:MM or similar")
}

// FormatTimeForDisplay formats time for user-friendly display
func FormatTimeForDisplay(t *time.Time) string {
	if t == nil {
		return "Not set"
	}
	return t.Format("2006-01-02 15:04")
}

// =============================================================================
// GENERIC MAP UTILITIES - Working with key-value data
// =============================================================================

// CreateOptionsFromMap creates selection options from a map
func CreateOptionsFromMap[T comparable](items map[string]T, labelFormatter func(key string, value T) string) []SelectionOption {
	options := make([]SelectionOption, 0, len(items))

	for key, value := range items {
		label := key
		if labelFormatter != nil {
			label = labelFormatter(key, value)
		}

		options = append(options, SelectionOption{
			Label: label,
			Value: key,
		})
	}

	return options
}

// FilterMapByStatus filters items by an active/inactive status
func FilterMapByStatus[T interface{ GetActive() bool }](items map[string]T, activeOnly bool) []string {
	var filtered []string

	for key, item := range items {
		if !activeOnly || item.GetActive() {
			filtered = append(filtered, key)
		}
	}

	return filtered
}

func FormatTime(t time.Time) string {
	return t.Format("Jan 2, 2006 3:04 PM")
}

// JSON Utilities
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

// JSON Query Utilities
func ExecuteJSONQuery(data any, query string) any {
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

// Response Parsing Utilities
func ParseResponse(response any) *model.HTTPResponse {
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

	httpResp := &model.HTTPResponse{
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

// Display Utilities
func DisplayResponse(response *model.HTTPResponse) {
	if string(response.Body) == "" {
		response.Body = []byte("Not set")
	}
	responseText := fmt.Sprintf("Status: %s\n\nBody:\n%s",
		response.Status,
		string(response.Body))

	DisplayFormattedText("🌐 HTTP Response", responseText)
}

func DisplayQueryResult(query string, result any) {
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
	DisplayFormattedText("🔍 Query Result", displayText)
}

// Key-Value Collection Utilities
func CollectKeyValuePairs(itemType, keyPlaceholder, valuePlaceholder string) map[string]string {
	items := make(map[string]string)

	for {
		inputs, err := AskMultipleInputs([]InputConfig{
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

		addMore, _ := AskConfirmation(fmt.Sprintf("Add another %s?", strings.ToLower(itemType)), "", "", "")
		if !addMore {
			break
		}
	}

	return items
}

// Enhanced response handler with JSON querying
func HandleResponse(response any, continueAction, returnAction func(), continueLabel, returnLabel string) {
	httpResp := ParseResponse(response)
	DisplayResponse(httpResp)

	if httpResp.IsJSON {
		handleJSONQuerying(httpResp, continueAction, returnAction, continueLabel, returnLabel)
	} else {
		AskContinueOrReturn(continueAction, returnAction, continueLabel, returnLabel)
	}
}

func handleJSONQuerying(response *model.HTTPResponse, continueAction, returnAction func(), continueLabel, returnLabel string) {
	for {
		choice, err := AskSelection("JSON Response Options", []SelectionOption{
			{"Query JSON", "query"},
			{"View Response", "response"},
			{"Continue", "continue"},
		})

		if err != nil {
			ShowError("Error in JSON options", err)
			break
		}

		switch choice {
		case "query":
			queryJSON(response)
		case "response":
			DisplayResponse(response)
		case "continue":
			AskContinueOrReturn(continueAction, returnAction, continueLabel, returnLabel)
			return
		}
	}
}

func queryJSON(response *model.HTTPResponse) {
	query, err := AskInput(InputConfig{
		Title:       "Enter JSON Query",
		Description: "Examples: .name, .users[0].email, .data.items",
		Placeholder: ".field.subfield",
	})

	if err != nil {
		ShowError("Error getting query input", err)
		return
	}

	if strings.TrimSpace(query) == "" {
		ShowMessage("No query provided")
		return
	}

	result := ExecuteJSONQuery(response.ParsedJSON, query)
	DisplayQueryResult(query, result)
}
