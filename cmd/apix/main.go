package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cf "github.com/Esa824/apix/internal/cli-forms"
	cc "github.com/Esa824/apix/internal/cobra-commands"
	hc "github.com/Esa824/apix/internal/http-client"
	"github.com/Esa824/apix/internal/model"
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
	rootCmd.Flags().BoolVar(&cliMode, "cli", false, "Enable interactive CLI mode")
	rootCmd.AddCommand(cc.GetCmd)
	rootCmd.AddCommand(cc.PostCmd)
	rootCmd.AddCommand(cc.PutCmd)
	rootCmd.AddCommand(cc.DeleteCmd)
	AddHTTPRequestFlags(cc.GetCmd)
	AddHTTPRequestFlags(cc.PostCmd)
	AddHTTPRequestFlags(cc.PutCmd)
	AddHTTPRequestFlags(cc.DeleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// AddHTTPRequestFlags adds comprehensive HTTP request flags to a cobra command
// This function organizes flags into logical groups for better UX
func AddHTTPRequestFlags(cmd *cobra.Command) *hc.RequestOptions {
	opts := &hc.RequestOptions{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Files:       make(map[string]string),
		Cookies:     make(map[string]string),
		Auth:        &model.Auth{},
		Context:     context.Background(),
		Time:        time.Now(),
	}

	// === BASIC REQUEST FLAGS ===
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name for this request (for templates/saving)")
	cmd.Flags().BoolVar(&opts.IsTemplate, "template", false, "Save this request as a template")

	// === HEADERS GROUP ===
	var headerStrings []string
	cmd.Flags().StringSliceVarP(&headerStrings, "header", "H", nil, "HTTP headers (format: 'Key:Value', repeatable)")
	cmd.Flags().String("content-type", "", "Content-Type header (shorthand)")
	cmd.Flags().String("accept", "", "Accept header (shorthand)")
	cmd.Flags().String("user-agent", "", "User-Agent header (shorthand)")

	// === QUERY PARAMETERS GROUP ===
	var queryStrings []string
	cmd.Flags().StringSliceVarP(&queryStrings, "query", "q", nil, "Query parameters (format: 'key=value', repeatable)")
	cmd.Flags().StringSliceVarP(&queryStrings, "param", "p", nil, "Alias for --query (format: 'key=value', repeatable)")

	// === BODY/DATA GROUP ===
	var bodyString string
	cmd.Flags().StringVarP(&bodyString, "data", "d", "", "Request body data (string, JSON, or @filename)")
	cmd.Flags().String("data-raw", "", "Request body data (raw string, no file processing)")
	cmd.Flags().String("json", "", "Send JSON data (sets Content-Type automatically)")
	cmd.Flags().String("form", "", "Send form data (application/x-www-form-urlencoded)")

	// === FILE UPLOADS GROUP ===
	var fileStrings []string
	cmd.Flags().StringSliceVarP(&fileStrings, "file", "F", nil, "File uploads (format: 'field=@/path/to/file', repeatable)")
	cmd.Flags().StringSlice("form-file", nil, "Form file uploads (format: 'field=@/path/to/file', repeatable)")

	// === COOKIES GROUP ===
	var cookieStrings []string
	cmd.Flags().StringSliceVarP(&cookieStrings, "cookie", "b", nil, "Cookies (format: 'name=value', repeatable)")
	cmd.Flags().String("cookie-jar", "", "Cookie jar file path")

	// === AUTHENTICATION GROUP ===
	cmd.Flags().StringVar(&opts.Auth.Type, "auth-type", "", "Authentication type (bearer, apikey, basic)")
	cmd.Flags().StringVar(&opts.Auth.Primary, "auth-token", "", "Authentication token/key/username")
	cmd.Flags().StringVar(&opts.Auth.Secondary, "auth-password", "", "Authentication password (for basic auth)")

	// Auth shortcuts
	cmd.Flags().String("bearer", "", "Bearer token (sets auth-type=bearer)")
	cmd.Flags().String("basic", "", "Basic auth (format: 'username:password')")
	cmd.Flags().String("apikey", "", "API key (sets auth-type=apikey)")

	// === TIMEOUT & BEHAVIOR GROUP ===
	var timeoutDuration time.Duration
	cmd.Flags().DurationVarP(&timeoutDuration, "timeout", "t", 30*time.Second, "Request timeout")
	cmd.Flags().Int("max-redirects", 10, "Maximum number of redirects to follow")
	cmd.Flags().Bool("no-redirect", false, "Don't follow redirects")
	cmd.Flags().Bool("insecure", false, "Skip SSL certificate verification")
	cmd.Flags().Bool("verbose", false, "Enable verbose output")

	// === OUTPUT & FORMATTING GROUP ===
	cmd.Flags().String("output", "", "Output file path")
	cmd.Flags().StringP("format", "f", "auto", "Output format (json, xml, text, auto)")
	cmd.Flags().Bool("include-headers", false, "Include response headers in output")
	cmd.Flags().Bool("only-headers", false, "Show only response headers")
	cmd.Flags().Bool("silent", false, "Silent mode (no output)")

	// === PROXY & NETWORK GROUP ===
	cmd.Flags().String("proxy", "", "Proxy URL (http://proxy:port)")
	cmd.Flags().String("interface", "", "Network interface to use")
	cmd.Flags().String("dns-servers", "", "Custom DNS servers (comma-separated)")

	// Set up flag groups for better help organization
	cmd.Flags().SortFlags = false // Maintain our custom order

	// Add validation in PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Parse headers
		if len(headerStrings) > 0 {
			for _, header := range headerStrings {
				if err := parseKeyValue(header, opts.Headers, ":"); err != nil {
					return fmt.Errorf("invalid header format '%s': %w", header, err)
				}
			}
		}

		// Parse query parameters
		if len(queryStrings) > 0 {
			for _, query := range queryStrings {
				if err := parseKeyValue(query, opts.QueryParams, "="); err != nil {
					return fmt.Errorf("invalid query parameter format '%s': %w", query, err)
				}
			}
		}

		// Parse cookies
		if len(cookieStrings) > 0 {
			for _, cookie := range cookieStrings {
				if err := parseKeyValue(cookie, opts.Cookies, "="); err != nil {
					return fmt.Errorf("invalid cookie format '%s': %w", cookie, err)
				}
			}
		}

		// Parse files
		if len(fileStrings) > 0 {
			for _, file := range fileStrings {
				if err := parseKeyValue(file, opts.Files, "="); err != nil {
					return fmt.Errorf("invalid file format '%s': %w", file, err)
				}
			}
		}

		// Handle shorthand headers
		if contentType := cmd.Flag("content-type").Value.String(); contentType != "" {
			opts.Headers["Content-Type"] = contentType
		}
		if accept := cmd.Flag("accept").Value.String(); accept != "" {
			opts.Headers["Accept"] = accept
		}
		if userAgent := cmd.Flag("user-agent").Value.String(); userAgent != "" {
			opts.Headers["User-Agent"] = userAgent
		}

		// Handle body data
		if bodyString != "" {
			opts.Body = bodyString
		}
		if jsonData := cmd.Flag("json").Value.String(); jsonData != "" {
			opts.Body = jsonData
			opts.Headers["Content-Type"] = "application/json"
		}
		if formData := cmd.Flag("form").Value.String(); formData != "" {
			opts.Body = formData
			opts.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}

		// Handle auth shortcuts
		if bearer := cmd.Flag("bearer").Value.String(); bearer != "" {
			opts.Auth.Type = "bearer"
			opts.Auth.Primary = bearer
		}
		if basic := cmd.Flag("basic").Value.String(); basic != "" {
			parts := strings.SplitN(basic, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("basic auth must be in format 'username:password'")
			}
			opts.Auth.Type = "basic"
			opts.Auth.Primary = parts[0]
			opts.Auth.Secondary = parts[1]
		}
		if apikey := cmd.Flag("apikey").Value.String(); apikey != "" {
			opts.Auth.Type = "apikey"
			opts.Auth.Primary = apikey
		}

		// Set timeout in context
		if timeoutDuration > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
			_ = cancel // We'll need to handle this properly in the actual implementation
			opts.Context = ctx
		}

		return nil
	}

	return opts
}

// parseKeyValue parses key-value pairs with a specified separator
func parseKeyValue(input string, target map[string]string, separator string) error {
	parts := strings.SplitN(input, separator, 2)
	if len(parts) != 2 {
		return fmt.Errorf("must be in format 'key%svalue'", separator)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	target[key] = value
	return nil
}
