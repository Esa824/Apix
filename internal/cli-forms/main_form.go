package cliforms

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	cc "apix/internal/cobra-commands"
)

func RunInteractiveMode() {
	var selectedOption string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Home area:").
				Options(
					huh.NewOption("Configuration", "configuration"),
					huh.NewOption("HTTP Requests", "http-requests"),
					huh.NewOption("Templates and History", "templates-and-history"),
					huh.NewOption("Authentication Management", "authentication-management"),
					huh.NewOption("Utilities", "utilities"),
					huh.NewOption("Settings", "settings"),
					huh.NewOption("Help", "help"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&selectedOption),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return
	}

	handleSelection(selectedOption)
}

func handleSelection(selection string) {
	switch selection {
	case "configuration":
		HandleConfiguration()
	case "http-requests":
		HandleHttpRequests()
	case "templates-and-history":
		HandleTemplatesAndHistory()
	case "delete":
		cc.HandleDeleteRequest()
	case "headers":
		handleHeaders()
	case "auth":
		handleAuthentication()
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Unknown selection")
	}
}
