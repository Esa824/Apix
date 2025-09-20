package cliforms

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Esa824/apix/internal/model"
	"github.com/Esa824/apix/internal/utils"
)

// Default settings
var AppSettings = &model.GlobalSettings{
	Display: model.DisplaySettings{
		ResponseFormat:  "pretty-json",
		ColorOutput:     true,
		ShowTiming:      true,
		ShowHeaders:     false,
		ShowStatusCode:  true,
		MaxResponseSize: 1024, // 1MB
		IndentSize:      2,
		LineNumbers:     false,
		SyntaxHighlight: true,
	},
	Behavior: model.BehaviorSettings{
		AutoSaveRequests:       false,
		ConfirmDeleteRequests:  true,
		ConfirmDestructive:     true,
		RequestTimeout:         30,
		MaxRetries:             3,
		RetryDelay:             1,
		FollowRedirects:        true,
		MaxRedirects:           5,
		ValidateSSL:            true,
		CacheResponses:         false,
		CacheDuration:          10,
		ShowProgressBar:        true,
		VerboseMode:            false,
		SaveFailedRequests:     true,
		AutoAddHeaders:         true,
		DefaultContentType:     "application/json",
		PreserveSessionCookies: true,
	},
	Network: model.NetworkSettings{
		DefaultTimeout:     30,
		ConnectTimeout:     10,
		ReadTimeout:        30,
		WriteTimeout:       30,
		MaxConnections:     10,
		UserAgent:          "Apix/1.0",
		ProxyEnabled:       false,
		KeepAlive:          true,
		CompressionEnabled: true,
	},
	Logging: model.LoggingSettings{
		EnableLogging:  false,
		LogLevel:       "info",
		LogFile:        "apix.log",
		LogRequests:    true,
		LogResponses:   true,
		LogHeaders:     false,
		LogTiming:      true,
		RotateLogFiles: true,
		MaxLogSize:     10,
		MaxLogFiles:    5,
	},
	Version:   "1.0.0",
	LastSaved: time.Now(),
}

func HandleSettingsManagement() {
	options := []utils.SelectionOption{
		{"Display Settings", "display"},
		{"Behavior Settings", "behavior"},
		{"Network Settings", "network"},
		{"Logging Settings", "logging"},
		{"Export Settings", "export"},
		{"Import Settings", "import"},
		{"Reset to Defaults", "reset"},
		{"Current Settings Overview", "overview"},
		{"Back to Main Menu", "back"},
	}

	selectedOption, err := utils.AskSelection("Settings Management:", options)
	if err != nil {
		utils.ShowError("Error running settings menu", err)
		return
	}

	handleSettingsSelection(selectedOption)
}

func handleSettingsSelection(selection string) {
	switch selection {
	case "display":
		handleDisplaySettings()
	case "behavior":
		handleBehaviorSettings()
	case "network":
		handleNetworkSettings()
	case "logging":
		handleLoggingSettings()
	case "export":
		handleExportSettings()
	case "import":
		handleImportSettings()
	case "reset":
		handleResetSettings()
	case "overview":
		showSettingsOverview()
	case "back":
		RunInteractiveMode()
	default:
		utils.ShowMessage("Unknown settings option")
	}
}

func handleDisplaySettings() {
	options := []utils.SelectionOption{
		{"Response Format", "response-format"},
		{"Color Output", "color-output"},
		{"Show Request Timing", "timing"},
		{"Show Headers", "headers"},
		{"Show Status Code", "status"},
		{"Max Response Size", "max-size"},
		{"Formatting Options", "formatting"},
		{"Back to Settings", "back"},
	}

	selectedOption, err := utils.AskSelection("Display Settings:", options)
	if err != nil {
		utils.ShowError("Error in display settings", err)
		return
	}

	switch selectedOption {
	case "response-format":
		handleResponseFormat()
	case "color-output":
		AppSettings.Display.ColorOutput = !AppSettings.Display.ColorOutput
		status := "disabled"
		if AppSettings.Display.ColorOutput {
			status = "enabled"
		}
		utils.ShowSuccess(fmt.Sprintf("Color output %s", status))
		askContinueOrReturnSettings()
	case "timing":
		AppSettings.Display.ShowTiming = !AppSettings.Display.ShowTiming
		status := "disabled"
		if AppSettings.Display.ShowTiming {
			status = "enabled"
		}
		utils.ShowSuccess(fmt.Sprintf("Request timing display %s", status))
		askContinueOrReturnSettings()
	case "headers":
		AppSettings.Display.ShowHeaders = !AppSettings.Display.ShowHeaders
		status := "disabled"
		if AppSettings.Display.ShowHeaders {
			status = "enabled"
		}
		utils.ShowSuccess(fmt.Sprintf("Headers display %s", status))
		askContinueOrReturnSettings()
	case "status":
		AppSettings.Display.ShowStatusCode = !AppSettings.Display.ShowStatusCode
		status := "disabled"
		if AppSettings.Display.ShowStatusCode {
			status = "enabled"
		}
		utils.ShowSuccess(fmt.Sprintf("Status code display %s", status))
		askContinueOrReturnSettings()
	case "max-size":
		handleMaxResponseSize()
	case "formatting":
		handleFormattingOptions()
	case "back":
		HandleSettingsManagement()
	}
}

func handleResponseFormat() {
	currentFormat := AppSettings.Display.ResponseFormat

	options := []utils.SelectionOption{
		{"Pretty JSON (formatted with indentation)", "pretty-json"},
		{"Raw (unmodified response)", "raw"},
		{"Headers Only (no body)", "headers-only"},
		{"Compact JSON (minified)", "compact-json"},
	}

	selectedFormat, err := utils.AskSelection(
		fmt.Sprintf("Response Format (current: %s):", strings.ToUpper(strings.ReplaceAll(currentFormat, "-", " "))),
		options,
	)
	if err != nil {
		utils.ShowError("Error selecting response format", err)
		return
	}

	AppSettings.Display.ResponseFormat = selectedFormat
	utils.ShowSuccess(fmt.Sprintf("Response format set to: %s", strings.ToUpper(strings.ReplaceAll(selectedFormat, "-", " "))))
	askContinueOrReturnSettings()
}

func handleMaxResponseSize() {
	currentSize := AppSettings.Display.MaxResponseSize
	sizeStr := "unlimited"
	if currentSize > 0 {
		sizeStr = fmt.Sprintf("%d KB", currentSize)
	}

	sizeConfig := utils.InputConfig{
		Title:       fmt.Sprintf("Max Response Size (current: %s):", sizeStr),
		Description: "Enter size in KB (0 = unlimited, max 10240 = 10MB)",
		Placeholder: "1024",
		Required:    false,
	}

	sizeInput, err := utils.AskInput(sizeConfig)
	if err != nil {
		utils.ShowError("Error setting response size", err)
		return
	}

	if sizeInput != "" {
		if size, parseErr := strconv.Atoi(sizeInput); parseErr == nil {
			if size < 0 || size > 10240 {
				utils.ShowWarning("Size must be between 0 and 10240 KB")
			} else {
				AppSettings.Display.MaxResponseSize = size
				if size == 0 {
					utils.ShowSuccess("Response size limit removed (unlimited)")
				} else {
					utils.ShowSuccess(fmt.Sprintf("Max response size set to %d KB", size))
				}
			}
		} else {
			utils.ShowError("Invalid size", parseErr)
		}
	}
	askContinueOrReturnSettings()
}

func handleFormattingOptions() {
	options := []utils.SelectionOption{
		{fmt.Sprintf("Indent Size (current: %d)", AppSettings.Display.IndentSize), "indent"},
		{fmt.Sprintf("Line Numbers (%s)", formatBoolStatus(AppSettings.Display.LineNumbers)), "line-numbers"},
		{fmt.Sprintf("Syntax Highlighting (%s)", formatBoolStatus(AppSettings.Display.SyntaxHighlight)), "syntax"},
		{"Back to Display Settings", "back"},
	}

	selectedOption, err := utils.AskSelection("Formatting Options:", options)
	if err != nil {
		utils.ShowError("Error in formatting options", err)
		return
	}

	switch selectedOption {
	case "indent":
		handleIndentSize()
	case "line-numbers":
		AppSettings.Display.LineNumbers = !AppSettings.Display.LineNumbers
		utils.ShowSuccess(fmt.Sprintf("Line numbers %s", formatBoolStatus(AppSettings.Display.LineNumbers)))
		askContinueOrReturnSettings()
	case "syntax":
		AppSettings.Display.SyntaxHighlight = !AppSettings.Display.SyntaxHighlight
		utils.ShowSuccess(fmt.Sprintf("Syntax highlighting %s", formatBoolStatus(AppSettings.Display.SyntaxHighlight)))
		askContinueOrReturnSettings()
	case "back":
		handleDisplaySettings()
	}
}

func handleIndentSize() {
	indentConfig := utils.InputConfig{
		Title:       fmt.Sprintf("Indent Size (current: %d):", AppSettings.Display.IndentSize),
		Description: "Number of spaces for indentation (2-8)",
		Placeholder: "2",
		Required:    false,
	}

	indentInput, err := utils.AskInput(indentConfig)
	if err != nil {
		utils.ShowError("Error setting indent size", err)
		return
	}

	if indentInput != "" {
		if size, parseErr := strconv.Atoi(indentInput); parseErr == nil {
			if size < 2 || size > 8 {
				utils.ShowWarning("Indent size must be between 2 and 8")
			} else {
				AppSettings.Display.IndentSize = size
				utils.ShowSuccess(fmt.Sprintf("Indent size set to %d spaces", size))
			}
		} else {
			utils.ShowError("Invalid indent size", parseErr)
		}
	}
	askContinueOrReturnSettings()
}

func handleBehaviorSettings() {
	options := []utils.SelectionOption{
		{fmt.Sprintf("Auto-save Requests (%s)", formatBoolStatus(AppSettings.Behavior.AutoSaveRequests)), "auto-save"},
		{fmt.Sprintf("Confirm DELETE Requests (%s)", formatBoolStatus(AppSettings.Behavior.ConfirmDeleteRequests)), "confirm-delete"},
		{fmt.Sprintf("Confirm Destructive Actions (%s)", formatBoolStatus(AppSettings.Behavior.ConfirmDestructive)), "confirm-destructive"},
		{fmt.Sprintf("Request Timeout (%ds)", AppSettings.Behavior.RequestTimeout), "timeout"},
		{fmt.Sprintf("Retry Settings (max: %d, delay: %ds)", AppSettings.Behavior.MaxRetries, AppSettings.Behavior.RetryDelay), "retry"},
		{fmt.Sprintf("Redirect Settings (%s, max: %d)", formatBoolStatus(AppSettings.Behavior.FollowRedirects), AppSettings.Behavior.MaxRedirects), "redirects"},
		{fmt.Sprintf("SSL Validation (%s)", formatBoolStatus(AppSettings.Behavior.ValidateSSL)), "ssl"},
		{fmt.Sprintf("Response Caching (%s)", formatBoolStatus(AppSettings.Behavior.CacheResponses)), "cache"},
		{"Advanced Behavior Options", "advanced"},
		{"Back to Settings", "back"},
	}

	selectedOption, err := utils.AskSelection("Behavior Settings:", options)
	if err != nil {
		utils.ShowError("Error in behavior settings", err)
		return
	}

	switch selectedOption {
	case "auto-save":
		AppSettings.Behavior.AutoSaveRequests = !AppSettings.Behavior.AutoSaveRequests
		utils.ShowSuccess(fmt.Sprintf("Auto-save requests %s", formatBoolStatus(AppSettings.Behavior.AutoSaveRequests)))
		askContinueOrReturnSettings()
	case "confirm-delete":
		AppSettings.Behavior.ConfirmDeleteRequests = !AppSettings.Behavior.ConfirmDeleteRequests
		utils.ShowSuccess(fmt.Sprintf("DELETE confirmation %s", formatBoolStatus(AppSettings.Behavior.ConfirmDeleteRequests)))
		askContinueOrReturnSettings()
	case "confirm-destructive":
		AppSettings.Behavior.ConfirmDestructive = !AppSettings.Behavior.ConfirmDestructive
		utils.ShowSuccess(fmt.Sprintf("Destructive action confirmation %s", formatBoolStatus(AppSettings.Behavior.ConfirmDestructive)))
		askContinueOrReturnSettings()
	case "timeout":
		handleRequestTimeout()
	case "retry":
		handleRetrySettings()
	case "redirects":
		handleRedirectSettings()
	case "ssl":
		AppSettings.Behavior.ValidateSSL = !AppSettings.Behavior.ValidateSSL
		utils.ShowSuccess(fmt.Sprintf("SSL validation %s", formatBoolStatus(AppSettings.Behavior.ValidateSSL)))
		askContinueOrReturnSettings()
	case "cache":
		handleCacheSettings()
	case "advanced":
		handleAdvancedBehavior()
	case "back":
		HandleSettingsManagement()
	}
}

func handleRequestTimeout() {
	timeoutConfig := utils.InputConfig{
		Title:       fmt.Sprintf("Request Timeout (current: %ds):", AppSettings.Behavior.RequestTimeout),
		Description: "Timeout in seconds (5-300)",
		Placeholder: "30",
		Required:    false,
	}

	timeoutInput, err := utils.AskInput(timeoutConfig)
	if err != nil {
		utils.ShowError("Error setting timeout", err)
		return
	}

	if timeoutInput != "" {
		if timeout, parseErr := strconv.Atoi(timeoutInput); parseErr == nil {
			if timeout < 5 || timeout > 300 {
				utils.ShowWarning("Timeout must be between 5 and 300 seconds")
			} else {
				AppSettings.Behavior.RequestTimeout = timeout
				utils.ShowSuccess(fmt.Sprintf("Request timeout set to %d seconds", timeout))
			}
		} else {
			utils.ShowError("Invalid timeout value", parseErr)
		}
	}
	askContinueOrReturnSettings()
}

func handleRetrySettings() {
	configs := []utils.InputConfig{
		{
			Title:       fmt.Sprintf("Max Retries (current: %d):", AppSettings.Behavior.MaxRetries),
			Description: "Number of retry attempts (0-10)",
			Placeholder: "3",
			Required:    false,
		},
		{
			Title:       fmt.Sprintf("Retry Delay (current: %ds):", AppSettings.Behavior.RetryDelay),
			Description: "Delay between retries in seconds (1-30)",
			Placeholder: "1",
			Required:    false,
		},
	}

	values, err := utils.AskMultipleInputs(configs)
	if err != nil {
		utils.ShowError("Error setting retry options", err)
		return
	}

	if values[0] != "" {
		if retries, parseErr := strconv.Atoi(values[0]); parseErr == nil && retries >= 0 && retries <= 10 {
			AppSettings.Behavior.MaxRetries = retries
		}
	}

	if values[1] != "" {
		if delay, parseErr := strconv.Atoi(values[1]); parseErr == nil && delay >= 1 && delay <= 30 {
			AppSettings.Behavior.RetryDelay = delay
		}
	}

	utils.ShowSuccess(fmt.Sprintf("Retry settings updated: %d retries with %ds delay",
		AppSettings.Behavior.MaxRetries, AppSettings.Behavior.RetryDelay))
	askContinueOrReturnSettings()
}

func handleRedirectSettings() {
	// Toggle follow redirects
	AppSettings.Behavior.FollowRedirects = !AppSettings.Behavior.FollowRedirects

	if AppSettings.Behavior.FollowRedirects {
		// Ask for max redirects
		maxConfig := utils.InputConfig{
			Title:       fmt.Sprintf("Max Redirects (current: %d):", AppSettings.Behavior.MaxRedirects),
			Description: "Maximum number of redirects to follow (1-20)",
			Placeholder: "5",
			Required:    false,
		}

		maxInput, err := utils.AskInput(maxConfig)
		if err == nil && maxInput != "" {
			if max, parseErr := strconv.Atoi(maxInput); parseErr == nil && max >= 1 && max <= 20 {
				AppSettings.Behavior.MaxRedirects = max
			}
		}
	}

	status := "disabled"
	if AppSettings.Behavior.FollowRedirects {
		status = fmt.Sprintf("enabled (max: %d)", AppSettings.Behavior.MaxRedirects)
	}
	utils.ShowSuccess(fmt.Sprintf("Redirect following %s", status))
	askContinueOrReturnSettings()
}

func handleCacheSettings() {
	AppSettings.Behavior.CacheResponses = !AppSettings.Behavior.CacheResponses

	if AppSettings.Behavior.CacheResponses {
		durationConfig := utils.InputConfig{
			Title:       fmt.Sprintf("Cache Duration (current: %dm):", AppSettings.Behavior.CacheDuration),
			Description: "Cache duration in minutes (1-60)",
			Placeholder: "10",
			Required:    false,
		}

		durationInput, err := utils.AskInput(durationConfig)
		if err == nil && durationInput != "" {
			if duration, parseErr := strconv.Atoi(durationInput); parseErr == nil && duration >= 1 && duration <= 60 {
				AppSettings.Behavior.CacheDuration = duration
			}
		}
	}

	status := "disabled"
	if AppSettings.Behavior.CacheResponses {
		status = fmt.Sprintf("enabled (%dm)", AppSettings.Behavior.CacheDuration)
	}
	utils.ShowSuccess(fmt.Sprintf("Response caching %s", status))
	askContinueOrReturnSettings()
}

func handleAdvancedBehavior() {
	options := []utils.SelectionOption{
		{fmt.Sprintf("Progress Bar (%s)", formatBoolStatus(AppSettings.Behavior.ShowProgressBar)), "progress"},
		{fmt.Sprintf("Verbose Mode (%s)", formatBoolStatus(AppSettings.Behavior.VerboseMode)), "verbose"},
		{fmt.Sprintf("Save Failed Requests (%s)", formatBoolStatus(AppSettings.Behavior.SaveFailedRequests)), "save-failed"},
		{fmt.Sprintf("Auto-add Headers (%s)", formatBoolStatus(AppSettings.Behavior.AutoAddHeaders)), "auto-headers"},
		{fmt.Sprintf("Default Content-Type (%s)", AppSettings.Behavior.DefaultContentType), "content-type"},
		{fmt.Sprintf("Preserve Cookies (%s)", formatBoolStatus(AppSettings.Behavior.PreserveSessionCookies)), "cookies"},
		{"Back to Behavior Settings", "back"},
	}

	selectedOption, err := utils.AskSelection("Advanced Behavior Options:", options)
	if err != nil {
		utils.ShowError("Error in advanced behavior settings", err)
		return
	}

	switch selectedOption {
	case "progress":
		AppSettings.Behavior.ShowProgressBar = !AppSettings.Behavior.ShowProgressBar
		utils.ShowSuccess(fmt.Sprintf("Progress bar %s", formatBoolStatus(AppSettings.Behavior.ShowProgressBar)))
		askContinueOrReturnSettings()
	case "verbose":
		AppSettings.Behavior.VerboseMode = !AppSettings.Behavior.VerboseMode
		utils.ShowSuccess(fmt.Sprintf("Verbose mode %s", formatBoolStatus(AppSettings.Behavior.VerboseMode)))
		askContinueOrReturnSettings()
	case "save-failed":
		AppSettings.Behavior.SaveFailedRequests = !AppSettings.Behavior.SaveFailedRequests
		utils.ShowSuccess(fmt.Sprintf("Save failed requests %s", formatBoolStatus(AppSettings.Behavior.SaveFailedRequests)))
		askContinueOrReturnSettings()
	case "auto-headers":
		AppSettings.Behavior.AutoAddHeaders = !AppSettings.Behavior.AutoAddHeaders
		utils.ShowSuccess(fmt.Sprintf("Auto-add headers %s", formatBoolStatus(AppSettings.Behavior.AutoAddHeaders)))
		askContinueOrReturnSettings()
	case "content-type":
		handleDefaultContentType()
	case "cookies":
		AppSettings.Behavior.PreserveSessionCookies = !AppSettings.Behavior.PreserveSessionCookies
		utils.ShowSuccess(fmt.Sprintf("Cookie preservation %s", formatBoolStatus(AppSettings.Behavior.PreserveSessionCookies)))
		askContinueOrReturnSettings()
	case "back":
		handleBehaviorSettings()
	}
}

func handleDefaultContentType() {
	options := []utils.SelectionOption{
		{"application/json", "application/json"},
		{"application/xml", "application/xml"},
		{"text/plain", "text/plain"},
		{"application/x-www-form-urlencoded", "application/x-www-form-urlencoded"},
		{"multipart/form-data", "multipart/form-data"},
		{"Custom", "custom"},
	}

	selectedType, err := utils.AskSelection(
		fmt.Sprintf("Default Content-Type (current: %s):", AppSettings.Behavior.DefaultContentType),
		options,
	)
	if err != nil {
		utils.ShowError("Error selecting content type", err)
		return
	}

	if selectedType == "custom" {
		customConfig := utils.InputConfig{
			Title:       "Custom Content-Type:",
			Description: "Enter custom MIME type",
			Placeholder: "application/vnd.api+json",
			Required:    true,
		}

		customType, err := utils.AskInput(customConfig)
		if err != nil {
			utils.ShowError("Error setting custom content type", err)
			return
		}
		selectedType = customType
	}

	AppSettings.Behavior.DefaultContentType = selectedType
	utils.ShowSuccess(fmt.Sprintf("Default content-type set to: %s", selectedType))
	askContinueOrReturnSettings()
}

func handleNetworkSettings() {
	utils.ShowMessage("Network settings functionality coming soon...")
	askContinueOrReturnSettings()
}

func handleLoggingSettings() {
	utils.ShowMessage("Logging settings functionality coming soon...")
	askContinueOrReturnSettings()
}

func handleExportSettings() {
	utils.ShowMessage("Export settings functionality coming soon...")
	askContinueOrReturnSettings()
}

func handleImportSettings() {
	utils.ShowMessage("Import settings functionality coming soon...")
	askContinueOrReturnSettings()
}

func handleResetSettings() {
	confirmed, err := utils.AskDangerousConfirmation(
		"Reset All Settings",
		"Reset all settings to default values? This will lose all your customizations",
		"all settings",
	)

	if err != nil {
		utils.ShowError("Error confirming reset", err)
		return
	}

	if confirmed {
		// Reset to defaults (create new instance)
		*AppSettings = model.GlobalSettings{
			Display: model.DisplaySettings{
				ResponseFormat:  "pretty-json",
				ColorOutput:     true,
				ShowTiming:      true,
				ShowHeaders:     false,
				ShowStatusCode:  true,
				MaxResponseSize: 1024,
				IndentSize:      2,
				LineNumbers:     false,
				SyntaxHighlight: true,
			},
			Behavior: model.BehaviorSettings{
				AutoSaveRequests:       false,
				ConfirmDeleteRequests:  true,
				ConfirmDestructive:     true,
				RequestTimeout:         30,
				MaxRetries:             3,
				RetryDelay:             1,
				FollowRedirects:        true,
				MaxRedirects:           5,
				ValidateSSL:            true,
				CacheResponses:         false,
				CacheDuration:          10,
				ShowProgressBar:        true,
				VerboseMode:            false,
				SaveFailedRequests:     true,
				AutoAddHeaders:         true,
				DefaultContentType:     "application/json",
				PreserveSessionCookies: true,
			},
			Version:   "1.0.0",
			LastSaved: time.Now(),
		}
		utils.ShowSuccess("All settings have been reset to defaults")
	} else {
		utils.ShowMessage("Settings reset cancelled")
	}
	askContinueOrReturnSettings()
}

func showSettingsOverview() {
	var overview strings.Builder

	overview.WriteString("Current Settings Overview\n")
	overview.WriteString("═══════════════════════════════════════\n\n")

	// Display Settings
	overview.WriteString("Display:\n")
	overview.WriteString("─────────────────────────────────\n")
	overview.WriteString(fmt.Sprintf("Response Format: %s\n", strings.ToUpper(strings.ReplaceAll(AppSettings.Display.ResponseFormat, "-", " "))))
	overview.WriteString(fmt.Sprintf("Color Output: %s\n", formatBoolStatus(AppSettings.Display.ColorOutput)))
	overview.WriteString(fmt.Sprintf("Show Timing: %s\n", formatBoolStatus(AppSettings.Display.ShowTiming)))
	overview.WriteString(fmt.Sprintf("Show Headers: %s\n", formatBoolStatus(AppSettings.Display.ShowHeaders)))
	overview.WriteString(fmt.Sprintf("Max Response Size: %s\n", formatSize(AppSettings.Display.MaxResponseSize)))
	overview.WriteString("\n")

	// Behavior Settings
	overview.WriteString("Behavior:\n")
	overview.WriteString("─────────────────────────────────\n")
	overview.WriteString(fmt.Sprintf("Auto-save Requests: %s\n", formatBoolStatus(AppSettings.Behavior.AutoSaveRequests)))
	overview.WriteString(fmt.Sprintf("Confirm DELETE: %s\n", formatBoolStatus(AppSettings.Behavior.ConfirmDeleteRequests)))
	overview.WriteString(fmt.Sprintf("Request Timeout: %ds\n", AppSettings.Behavior.RequestTimeout))
	overview.WriteString(fmt.Sprintf("Max Retries: %d\n", AppSettings.Behavior.MaxRetries))
	overview.WriteString(fmt.Sprintf("Follow Redirects: %s\n", formatBoolStatus(AppSettings.Behavior.FollowRedirects)))
	overview.WriteString(fmt.Sprintf("SSL Validation: %s\n", formatBoolStatus(AppSettings.Behavior.ValidateSSL)))
	overview.WriteString("\n")

	overview.WriteString("═══════════════════════════════════════")
	overview.WriteString(fmt.Sprintf("\nLast Updated: %s", AppSettings.LastSaved.Format("2006-01-02 15:04:05")))

	utils.DisplayFormattedText("Settings Overview", overview.String())
	askContinueOrReturnSettings()
}

func askContinueOrReturnSettings() {
	utils.AskContinueOrReturn(
		HandleSettingsManagement,
		RunInteractiveMode,
		"Continue with Settings",
		"Return to Main Menu",
	)
}

// Helper functions
func formatBoolStatus(value bool) string {
	if value {
		return "ON"
	}
	return "OFF"
}

func formatSize(sizeKB int) string {
	if sizeKB == 0 {
		return "Unlimited"
	}
	return fmt.Sprintf("%d KB", sizeKB)
}

// GetCurrentSettings returns the current settings
func GetCurrentSettings() *model.GlobalSettings {
	return AppSettings
}

// SaveSettings updates the last saved timestamp
func SaveSettings() {
	AppSettings.LastSaved = time.Now()
}
