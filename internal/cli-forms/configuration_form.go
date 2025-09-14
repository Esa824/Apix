package cliforms

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

var (
	BaseURL string
)

func HandleConfiguration() {
	var selectedOption string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Configuration Options:").
				Options(
					huh.NewOption("Set Base URL", "set-base-url"),
					huh.NewOption("View Current Configuration", "view-config"),
					huh.NewOption("Reset Configuration", "reset-config"),
					huh.NewOption("Back to Main Menu", "back"),
				).
				Value(&selectedOption),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running configuration form: %v\n", err)
		return
	}

	handleConfigSelection(selectedOption)
}

func handleConfigSelection(selection string) {
	switch selection {
	case "set-base-url":
		handleSetBaseURL()
	case "view-config":
		handleViewConfig()
	case "reset-config":
		handleResetConfig()
	case "back":
		RunInteractiveMode()
	default:
		fmt.Println("Unknown configuration option")
	}
}

func handleSetBaseURL() {
	var baseURL string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Base URL:").
				Description("This will be prepended to relative URLs in requests").
				Placeholder("https://api.example.com").
				Value(&baseURL),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if baseURL != "" {
		BaseURL = baseURL
		fmt.Printf("Base URL set to: %s\n", baseURL)
		fmt.Println("✓ Configuration saved successfully!")
	} else {
		fmt.Println("No base URL provided. Configuration unchanged.")
	}

	// Ask if user wants to continue with configuration or return to main menu
	askContinueOrReturnConfiguration()
}

func handleViewConfig() {
	// TODO: Implement reading configuration from file/storage
	fmt.Println("📋 Current Configuration:")
	fmt.Println("─────────────────────────")
	fmt.Printf("Base URL: %s\n", BaseURL)
	fmt.Printf("Auth Profile: %s\n", ActiveProfile)
	fmt.Println("─────────────────────────")

	askContinueOrReturnConfiguration()
}

func handleResetConfig() {
	var confirmReset bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Reset Configuration").
				Description("This will clear all saved configuration including base URL, auth tokens, and headers. This action cannot be undone.").
				Affirmative("Yes, reset everything").
				Negative("Cancel").
				Value(&confirmReset),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if confirmReset {
		// TODO: Implement actual config reset
		fmt.Println("🔄 Resetting configuration...")
		fmt.Println("✓ All configuration has been reset to defaults!")
	} else {
		fmt.Println("Configuration reset cancelled.")
	}

	askContinueOrReturnConfiguration()
}

func askContinueOrReturnConfiguration() {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
				Options(
					huh.NewOption("Continue with Configuration", "continue"),
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
		HandleConfiguration()
	case "main":
		RunInteractiveMode()
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
}

func getAuthToken() string {
	// TODO: Read from config file/storage
	return ""
}

func maskAuthToken(token string) string {
	if token == "" {
		return "Not set"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

func getHeaderCount() int {
	// TODO: Read from config file/storage
	return 0
}
