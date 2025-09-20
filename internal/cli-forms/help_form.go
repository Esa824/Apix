package cliforms

import (
	"fmt"
	"strings"

	"github.com/Esa824/apix/internal/model"
	"github.com/Esa824/apix/internal/utils"
)

// Global help content
var HelpSections = map[string]*model.HelpSection{
	"quick-start": {
		Title:       "Quick Start Guide",
		Description: "Get started with Apix in minutes",
		Content: `Welcome to Apix! Here's how to get started:

1. Interactive Mode: Run 'apix --cli' for guided forms
2. Direct Commands: Use 'apix get <url>' for quick requests
3. Authentication: Set up auth profiles for secure APIs
4. Response Formatting: Built-in JSON formatting and error handling

Basic workflow:
Choose HTTP method (GET, POST, PUT, DELETE)
Enter your API endpoint URL
Add request body (for POST/PUT)
View formatted response`,
		Examples: []string{
			"apix --cli",
			"apix get https://api.github.com/users/octocat",
			"apix post https://api.example.com/users",
		},
	},
	"commands": {
		Title:       "Command Examples",
		Description: "Learn by example with common API patterns",
		Content: `Here are common usage patterns for different API scenarios:

GET Requests:
Fetch user data
List resources
Check API status

POST Requests:
Create new resources
Submit form data
Upload content

PUT Requests:
Update existing resources
Replace entire objects
Modify configurations

DELETE Requests:
Remove resources
Clean up data
Revoke permissions`,
		Examples: []string{
			"# Get a user profile",
			"apix get https://jsonplaceholder.typicode.com/users/1",
			"",
			"# Create a new post",
			"apix post https://jsonplaceholder.typicode.com/posts",
			"",
			"# Update a resource",
			"apix put https://jsonplaceholder.typicode.com/posts/1",
			"",
			"# Delete a resource",
			"apix delete https://jsonplaceholder.typicode.com/posts/1",
		},
	},
	"shortcuts": {
		Title:       "Keyboard Shortcuts",
		Description: "Navigate faster with keyboard shortcuts",
		Content: `Speed up your workflow with these keyboard shortcuts:

Interactive Mode Navigation:
Tab / Shift+Tab: Move between fields
Enter: Confirm selection
Esc: Cancel current operation
Space: Toggle checkboxes/confirmations
Arrow Keys: Navigate options

Form Controls:
Ctrl+C: Copy current input
Ctrl+V: Paste into input field
Ctrl+A: Select all text
Ctrl+U: Clear current line
Backspace: Delete character

Quick Actions:
Ctrl+D: Exit application
Ctrl+L: Clear screen (in supported terminals)`,
		Examples: []string{
			"Tab Move to next field",
			"Enter Submit form",
			"Esc Cancel operation",
			"↑→ Navigate menu options",
			"Space Toggle selection",
		},
	},
	"api-patterns": {
		Title:       "Common API Patterns",
		Description: "Real-world API usage examples",
		Content: `Learn common API integration patterns:

REST API Basics:
GET /users - List all users
GET /users/123 - Get specific user
POST /users - Create new user
PUT /users/123 - Update user
DELETE /users/123 - Delete user

Authentication Patterns:
Bearer Token: Authorization: Bearer <token>
API Key: X-API-Key: <key>
Basic Auth: username:password encoded

Common Headers:
Content-Type: application/json
Accept: application/json
User-Agent: Apix/1.0`,
		Examples: []string{
			"# GitHub API with token",
			"apix get https://api.github.com/user",
			"# Header: Authorization: Bearer ghp_xxxx",
			"",
			"# REST API with JSON",
			"apix post https://api.example.com/users",
			"'# Body: {'name': 'John', 'email': 'john@example.com'}'",
			"",
			"# API with custom headers",
			"apix get https://api.service.com/data",
			"# Header: X-API-Key: your-api-key",
		},
	},
	"troubleshooting": {
		Title:       "Troubleshooting Guide",
		Description: "Common issues and solutions",
		Content: `Having trouble? Here are solutions to common issues:

Connection Issues:
Check your internet connection
Verify the API endpoint URL
Ensure the API server is running
Check firewall/proxy settings

Authentication Errors:
Verify your API key/token is correct
Check if the token has expired
Ensure proper header format
Confirm API permissions

Response Issues:
Check API documentation for expected format
Verify request method (GET/POST/PUT/DELETE)
Ensure required fields are provided
Check Content-Type headers

General Tips:
Start with simple GET requests
Use public APIs for testing
Check API rate limits
Review error messages carefully`,
		Examples: []string{
			"# Test with a simple public API",
			"apix get https://httpbin.org/get",
			"",
			"# Check what you're sending",
			"apix post https://httpbin.org/post",
			"# This endpoint echoes your request",
			"",
			"# Verify SSL issues",
			"apix get https://httpbin.org/status/200",
		},
	},
}

func HandleHelpAndDocumentation() {
	options := []utils.SelectionOption{
		{"Quick Start Guide", "quick-start"},
		{"Command Examples", "commands"},
		{"Keyboard Shortcuts", "shortcuts"},
		{"Common API Patterns", "api-patterns"},
		{"Troubleshooting Guide", "troubleshooting"},
		{"View All Documentation", "view-all"},
		{"Back to Main Menu", "back"},
	}

	selectedOption, err := utils.AskSelection("Help & Documentation:", options)
	if err != nil {
		utils.ShowError("Error running help menu", err)
		return
	}

	handleHelpSelection(selectedOption)
}

func handleHelpSelection(selection string) {
	switch selection {
	case "quick-start":
		showHelpSection("quick-start")
	case "commands":
		showHelpSection("commands")
	case "shortcuts":
		showHelpSection("shortcuts")
	case "api-patterns":
		showHelpSection("api-patterns")
	case "troubleshooting":
		showHelpSection("troubleshooting")
	case "view-all":
		showAllDocumentation()
	case "back":
		RunInteractiveMode()
	default:
		utils.ShowMessage("Unknown help option")
	}
}

func showHelpSection(sectionKey string) {
	section, exists := HelpSections[sectionKey]
	if !exists {
		utils.ShowError("Help section not found", fmt.Errorf("section '%s' does not exist", sectionKey))
		askContinueOrReturnHelp()
		return
	}

	var content strings.Builder

	// Add the main content
	content.WriteString(section.Content)
	content.WriteString("\n\n")

	// Add examples if they exist
	if len(section.Examples) > 0 {
		content.WriteString("Examples:\n")
		content.WriteString("─────────────────────────────────\n")
		for _, example := range section.Examples {
			if strings.HasPrefix(example, "#") {
				// Comment line
				content.WriteString(fmt.Sprintf("  %s\n", example))
			} else if example == "" {
				// Empty line for spacing
				content.WriteString("\n")
			} else {
				// Command example
				content.WriteString(fmt.Sprintf("  $ %s\n", example))
			}
		}
	}

	utils.DisplayFormattedText(fmt.Sprintf("%s", section.Title), content.String())
	askContinueOrReturnHelp()
}

func showAllDocumentation() {
	var allContent strings.Builder

	allContent.WriteString("Complete Apix Documentation\n")
	allContent.WriteString("═══════════════════════════════════════════\n\n")

	// Order the sections for logical flow
	orderedSections := []string{"quick-start", "commands", "shortcuts", "api-patterns", "troubleshooting"}

	for i, sectionKey := range orderedSections {
		if i > 0 {
			allContent.WriteString("\n\n")
		}

		section := HelpSections[sectionKey]
		allContent.WriteString(fmt.Sprintf("%s\n", section.Title))
		allContent.WriteString("─────────────────────────────────\n")
		allContent.WriteString(section.Description)
		allContent.WriteString("\n\n")
		allContent.WriteString(section.Content)

		if len(section.Examples) > 0 {
			allContent.WriteString("\n\nExamples:\n")
			for _, example := range section.Examples {
				if strings.HasPrefix(example, "#") {
					allContent.WriteString(fmt.Sprintf("  %s\n", example))
				} else if example == "" {
					allContent.WriteString("\n")
				} else {
					allContent.WriteString(fmt.Sprintf("  $ %s\n", example))
				}
			}
		}
	}

	allContent.WriteString("\n\n═══════════════════════════════════════════")
	allContent.WriteString("\nNeed more help? Check the project README or open an issue!")

	utils.DisplayFormattedText("Complete Documentation", allContent.String())
	askContinueOrReturnHelp()
}

func askContinueOrReturnHelp() {
	utils.AskContinueOrReturn(
		HandleHelpAndDocumentation,
		RunInteractiveMode,
		"Browse More Help Topics",
		"Return to Main Menu",
	)
}

// GetHelpSection returns a specific help section
func GetHelpSection(sectionKey string) *model.HelpSection {
	return HelpSections[sectionKey]
}

// GetAllHelpSections returns all help sections
func GetAllHelpSections() map[string]*model.HelpSection {
	return HelpSections
}

// AddCustomHelpSection allows adding custom help content
func AddCustomHelpSection(key string, section *model.HelpSection) {
	HelpSections[key] = section
}

// SearchHelpContent searches through all help content for a term
func SearchHelpContent(searchTerm string) []string {
	var results []string
	searchLower := strings.ToLower(searchTerm)

	for _, section := range HelpSections {
		// Search in title
		if strings.Contains(strings.ToLower(section.Title), searchLower) {
			results = append(results, fmt.Sprintf("%s - %s", section.Title, section.Description))
		}

		// Search in content
		if strings.Contains(strings.ToLower(section.Content), searchLower) {
			results = append(results, fmt.Sprintf("Found in %s", section.Title))
		}

		// Search in examples
		for _, example := range section.Examples {
			if strings.Contains(strings.ToLower(example), searchLower) {
				results = append(results, fmt.Sprintf("Example in %s: %s", section.Title, example))
				break // Only show one example match per section
			}
		}
	}

	return results
}
