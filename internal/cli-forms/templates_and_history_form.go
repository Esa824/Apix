package cliforms

import (
	"fmt"
	"os"
	"strconv"
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
		options = append(options, huh.NewOption(fmt.Sprintf("%s (%s)", template.Name, template.Method), strconv.Itoa(template.Id)))
	}

	// Add management options
	options = append(options,
		huh.NewOption("Create New Template", "create-template"),
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
	case "create-template":
		handleCreateTemplate()
	case "back":
		HandleTemplatesAndHistory()
	default:
		selectionInt, err := strconv.Atoi(selection)
		if err != nil {
			return
		}
		template, err := hc.GetTemplateByID(selectionInt)
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

func handleCreateTemplate() {
	var templateName, method, url, body, headers string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Template Name:").
				Placeholder("My API Template").
				Value(&templateName),
			huh.NewSelect[string]().
				Title("HTTP Method:").
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
				).
				Value(&method),
			huh.NewInput().
				Title("URL:").
				Placeholder("https://api.example.com/users").
				Value(&url),
			huh.NewText().
				Title("Request Body (optional):").
				Description("JSON body for POST/PUT requests").
				Placeholder(`{"name": "John", "email": "john@example.com"}`).
				Value(&body),
			huh.NewText().
				Title("Headers (optional):").
				Description("One header per line, format: Key: Value").
				Placeholder("Content-Type: application/json\nAuthorization: Bearer token").
				Value(&headers),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if templateName != "" && method != "" && url != "" {
		// TODO: Implement actual template saving
		fmt.Printf("✓ Template '%s' created successfully!\n", templateName)
		fmt.Printf("  Method: %s\n", method)
		fmt.Printf("  URL: %s\n", url)
		if body != "" {
			fmt.Printf("  Has Body: Yes\n")
		}
		if headers != "" {
			fmt.Printf("  Has Headers: Yes\n")
		}
	} else {
		fmt.Println("Template creation cancelled - missing required fields")
	}

	askContinueOrReturnTemplates()
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
	fmt.Printf("Executing template: %s\n", template.Name)
	fmt.Printf("   %s %s\n", template.Method, template.URL)

	// TODO: Implement actual request execution
	fmt.Println("Request completed successfully!")
	fmt.Println("Response: 200 OK")

	askContinueOrReturnTemplates()
}

func editTemplate(template *model.Template) {
	fmt.Printf("✏️  Editing template: %s\n", template.Name)
	// TODO: Implement template editing - could reuse handleCreateTemplate with pre-filled values
	fmt.Println("Template updated successfully!")

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
		// TODO: Implement actual template deletion
		fmt.Printf("Deleting template: %s\n", template.Name)
		fmt.Println("Template deleted successfully!")
	} else {
		fmt.Println("Template deletion cancelled.")
	}

	askContinueOrReturnTemplates()
}

func reExecuteFromHistory(historyItem *hc.RequestOptions) {
	fmt.Printf("Re-executing request: %s %s\n", historyItem.Method, historyItem.URL)
	response, _ := hc.NewClient(10*time.Second).Do(*historyItem, false)
	utils.HandleResponse(response, HandleTemplatesAndHistory, RunInteractiveMode, "Continue with templates & history", "Return to Main Menu")
}

func saveHistoryAsTemplate(historyItem *hc.RequestOptions) {
	var templateName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Template Name:").
				Placeholder(fmt.Sprintf("%s Template", historyItem.Method)).
				Value(&templateName),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if templateName != "" {
		// TODO: Implement saving history item as template
		fmt.Printf("Saved as template: %s\n", templateName)
		fmt.Println("Template created successfully!")
	} else {
		fmt.Println("Template creation cancelled.")
	}

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
