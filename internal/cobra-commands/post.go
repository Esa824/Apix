package cobracommands

import (
	"fmt"
	"time"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/utils"

	"github.com/spf13/cobra"
)

var PostCmd = &cobra.Command{
	Use:   "post [URL]",
	Short: "Make a POST request to the specified URL",
	Long:  `Make a POST request to the specified URL with optional body, headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  handlePostRequest,
}

func handlePostRequest(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Argument not provided")
	}
	opts, err := GetRequestOptionsFromFlags(cmd)
	if err != nil {
		return err
	}
	opts.Method = "POST"
	opts.URL = utils.NormalizeURL(args[0])
	response, err := hc.NewClient(10*time.Second).Do(*opts, true)
	fmt.Print(response)
	return err
}
