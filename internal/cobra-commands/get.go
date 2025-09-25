package cobracommands

import (
	"fmt"
	"time"

	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/utils"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get [URL]",
	Short: "Make a GET request to the specified URL",
	Long:  `Make a GET request to the specified URL with optional headers and parameters.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  handleGetRequest,
}

func handleGetRequest(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Argument not provided")
	}
	opts, err := GetRequestOptionsFromFlags(cmd)
	if err != nil {
		return err
	}
	opts.Method = "GET"
	opts.URL = utils.NormalizeURL(args[0])
	response, err := hc.NewClient(10*time.Second).Do(*opts, true)
	fmt.Print(response)
	return err
}
