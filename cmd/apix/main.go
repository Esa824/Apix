package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cf "apix/internal/cli-forms"
	cc "apix/internal/cobra-commands"
)

var (
	cliMode bool
)

var rootCmd = &cobra.Command{
	Use:   "apix",
	Short: "A CLI application that simplifies API testing and interaction",
	Long: `Apix is a CLI application that simplifies API testing and interaction. 
Unlike curl's verbose syntax and complex flag management, Apix provides an 
intuitive interface for making HTTP requests, managing authentication, 
and handling responses with built-in formatting and error handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cliMode {
			cf.RunInteractiveMode()
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&cliMode, "cli", false, "Enable interactive CLI mode")
	rootCmd.AddCommand(cc.GetCmd)
	rootCmd.AddCommand(cc.PostCmd)
	rootCmd.AddCommand(cc.PutCmd)
	rootCmd.AddCommand(cc.DeleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

