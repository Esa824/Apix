package cobracommands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [URL]",
	Short: "Make a DELETE request to the specified URL",
	Long:  `Make a DELETE request to the specified URL with optional headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Making DELETE request to: %s\n", args[0])
		// TODO: Implement DELETE request logic
	},
}

func HandleDeleteRequest() {
	var url string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter URL for DELETE request:").
				Placeholder("https://api.example.com/users/1").
				Value(&url),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if url != "" {
		fmt.Printf("Making DELETE request to: %s\n", url)
		// TODO: Implement actual DELETE request
	}
}
