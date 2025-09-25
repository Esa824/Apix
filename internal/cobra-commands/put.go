package cobracommands

import (
	"fmt"
	"time"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/utils"

	"github.com/spf13/cobra"
)

var PutCmd = &cobra.Command{
	Use:   "put [URL]",
	Short: "Make a PUT request to the specified URL",
	Long:  `Make a PUT request to the specified URL with optional body, headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  handlePutRequest,
}

func handlePutRequest(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Argument not provided")
	}
	opts, err := GetRequestOptionsFromFlags(cmd)
	if err != nil {
		return err
	}
	opts.Method = "PUT"
	opts.URL = utils.NormalizeURL(args[0])
	response, err := hc.NewClient(10*time.Second).Do(*opts, true)
	fmt.Print(response)
	return err
}
