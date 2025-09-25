package cobracommands

import (
	"context"
	"fmt"
	"time"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/model"
	"github.com/Esa824/apix/internal/utils"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [URL]",
	Short: "Make a DELETE request to the specified URL",
	Long:  `Make a DELETE request to the specified URL with optional headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  handleDeleteRequest,
}

func handleDeleteRequest(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Argument not provided")
	}
	opts, err := GetRequestOptionsFromFlags(cmd)
	if err != nil {
		return err
	}
	opts.Method = "DELETE"
	opts.URL = utils.NormalizeURL(args[0])
	response, err := hc.NewClient(10*time.Second).Do(*opts, true)
	fmt.Print(response)
	return err
}

// GetRequestOptionsFromFlags extracts and returns the configured RequestOptions
func GetRequestOptionsFromFlags(cmd *cobra.Command) (*hc.RequestOptions, error) {
	// This would be called in the command's RunE function
	// The actual opts would be populated by the PreRunE validation

	// For demonstration, showing how you'd access flags:
	timeout, _ := cmd.Flags().GetDuration("timeout")
	verbose, _ := cmd.Flags().GetBool("verbose")

	opts := &hc.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Files:       make(map[string]string),
		Cookies:     make(map[string]string),
		Auth:        &model.Auth{},
		Time:        time.Now(),
	}

	// Set context with timeout
	if timeout > 0 {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		opts.Context = ctx
	} else {
		opts.Context = context.Background()
	}

	// The PreRunE function would have already populated the maps
	// This is just showing the structure

	if verbose {
		fmt.Printf("Request Options: URL=%s, Method=%s\n", opts.URL, opts.Method)
	}

	return opts, nil
}
