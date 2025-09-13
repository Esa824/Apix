package cobracommands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var PutCmd = &cobra.Command{
	Use:   "put [URL]",
	Short: "Make a PUT request to the specified URL",
	Long:  `Make a PUT request to the specified URL with optional body, headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Making PUT request to: %s\n", args[0])
		// TODO: Implement PUT request logic
	},
}

func HandlePutRequest() {
	var url, body string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter URL for PUT request:").
				Placeholder("https://api.example.com/users/1").
				Value(&url),
			huh.NewText().
				Title("Enter request body (JSON):").
				Placeholder(`{"name": "John Updated", "email": "john.updated@example.com"}`).
				Value(&body),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if url != "" {
		fmt.Printf("Making PUT request to: %s\n", url)
		if body != "" {
			fmt.Printf("With body: %s\n", body)
		}
		// TODO: Implement actual PUT request
	}
}
