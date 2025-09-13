package cobracommands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get [URL]",
	Short: "Make a GET request to the specified URL",
	Long:  `Make a GET request to the specified URL with optional headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Making GET request to: %s\n", args[0])
		// TODO: Implement GET request logic
	},
}

func HandleGetRequest() {
	var url string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter URL for GET request:").
				Placeholder("https://api.example.com/users").
				Value(&url),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if url != "" {
		fmt.Printf("Making GET request to: %s\n", url)
		// TODO: Implement actual GET request
	}
}
