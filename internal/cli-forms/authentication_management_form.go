package cliforms

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"apix/internal/utils"
)

// AuthProfile represents an authentication profile
type AuthProfile struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"` // "bearer", "apikey", "basic", "oauth"
	Token    string     `json:"token"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	APIKey   string     `json:"api_key"`
	Header   string     `json:"header"`
	Expiry   *time.Time `json:"expiry"`
	Active   bool       `json:"active"`
}

// Global storage for auth profiles (in a real app, this would be persisted)
var AuthProfiles = make(map[string]*AuthProfile)
var ActiveProfile string

// initAuthProfilesDir ensures the auth-profiles directory exists
func initAuthProfilesDir() error {
	authProfilesDir := filepath.Join(ConfigPath, "auth-profiles")
	return os.MkdirAll(authProfilesDir, 0755)
}

// saveAuthProfile saves a single auth profile to a JSON file
func saveAuthProfile(profile *AuthProfile) error {
	if err := initAuthProfilesDir(); err != nil {
		return fmt.Errorf("failed to create auth-profiles directory: %w", err)
	}

	authProfilesDir := filepath.Join(ConfigPath, "auth-profiles")
	filename := fmt.Sprintf("%s.json", profile.Name)
	filepath := filepath.Join(authProfilesDir, filename)

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth profile: %w", err)
	}

	return os.WriteFile(filepath, data, 0600)
}

// deleteAuthProfileFile deletes the JSON file for an auth profile
func deleteAuthProfileFile(profileName string) error {
	authProfilesDir := filepath.Join(ConfigPath, "auth-profiles")
	filename := fmt.Sprintf("%s.json", profileName)
	filepath := filepath.Join(authProfilesDir, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete auth profile file: %w", err)
	}
	return nil
}

// loadAuthProfiles loads all auth profiles from the auth-profiles directory
func loadAuthProfiles() error {
	authProfilesDir := filepath.Join(ConfigPath, "auth-profiles")

	// Check if directory exists
	if _, err := os.Stat(authProfilesDir); os.IsNotExist(err) {
		return nil // No profiles directory yet, that's okay
	}

	entries, err := os.ReadDir(authProfilesDir)
	if err != nil {
		return fmt.Errorf("failed to read auth-profiles directory: %w", err)
	}

	AuthProfiles = make(map[string]*AuthProfile)
	ActiveProfile = ""

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filepath := filepath.Join(authProfilesDir, entry.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			utils.ShowWarning(fmt.Sprintf("Failed to read auth profile file %s: %v", entry.Name(), err))
			continue
		}

		var profile AuthProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			utils.ShowWarning(fmt.Sprintf("Failed to parse auth profile file %s: %v", entry.Name(), err))
			continue
		}

		AuthProfiles[profile.Name] = &profile
		if profile.Active {
			if ActiveProfile != "" {
				// Multiple active profiles found, deactivate others
				AuthProfiles[ActiveProfile].Active = false
			}
			ActiveProfile = profile.Name
		}
	}

	return nil
}

func HandleAuthenticationManagement() {
	// Load existing profiles when starting auth management
	if err := loadAuthProfiles(); err != nil {
		utils.ShowError("Error loading auth profiles", err)
	}

	options := []utils.SelectionOption{
		{"Create New Auth Profile", "create-profile"},
		{"Select Active Profile", "select-profile"},
		{"Edit Existing Profile", "edit-profile"},
		{"Delete Profile", "delete-profile"},
		{"View All Profiles", "view-profiles"},
		{"Back to Main Menu", "back"},
	}

	selectedOption, err := utils.AskSelection("🔐 Authentication Management:", options)
	if err != nil {
		utils.ShowError("Error running authentication form", err)
		return
	}

	handleAuthSelection(selectedOption)
}

func handleAuthSelection(selection string) {
	switch selection {
	case "create-profile":
		handleCreateProfile()
	case "select-profile":
		handleSelectProfile()
	case "edit-profile":
		handleEditProfile()
	case "delete-profile":
		handleDeleteProfile()
	case "view-profiles":
		handleViewProfiles()
	case "back":
		RunInteractiveMode()
	default:
		utils.ShowMessage("Unknown authentication option")
	}
}

func handleCreateProfile() {
	// Step 1: Get profile name and auth type
	configs := []utils.InputConfig{
		{
			Title:       "Profile Name:",
			Description: "Give this auth profile a memorable name",
			Placeholder: "production-api",
			Required:    true,
		},
	}

	values, err := utils.AskMultipleInputs(configs)
	if err != nil {
		utils.ShowError("Error creating profile", err)
		return
	}

	profileName := values[0]

	if _, exists := AuthProfiles[profileName]; exists {
		utils.ShowMessage(fmt.Sprintf("Profile '%s' already exists. Choose a different name.", profileName))
		askContinueOrReturnAuth()
		return
	}

	// Step 2: Select auth type
	authOptions := []utils.SelectionOption{
		{"Bearer Token", "bearer"},
		{"API Key", "apikey"},
		{"Basic Authentication", "basic"},
		{"OAuth 2.0 (Coming Soon)", "oauth"},
	}

	authType, err := utils.AskSelection("Authentication Type:", authOptions)
	if err != nil {
		utils.ShowError("Error selecting authentication type", err)
		return
	}

	if authType == "oauth" {
		utils.ShowMessage("OAuth 2.0 support is coming soon!")
		askContinueOrReturnAuth()
		return
	}

	// Step 3: Setup auth-specific details
	profile := &AuthProfile{
		Name: profileName,
		Type: authType,
	}

	var success bool
	switch authType {
	case "bearer":
		success = handleBearerTokenSetup(profile)
	case "apikey":
		success = handleAPIKeySetup(profile)
	case "basic":
		success = handleBasicAuthSetup(profile)
	}

	if !success {
		return
	}

	// Save profile to memory and disk
	AuthProfiles[profileName] = profile
	if err := saveAuthProfile(profile); err != nil {
		utils.ShowError("Failed to save auth profile", err)
		delete(AuthProfiles, profileName) // Remove from memory if save failed
		return
	}

	utils.ShowSuccess(fmt.Sprintf("Auth profile '%s' created successfully!", profileName))

	// Ask if this should be the active profile
	setAsActive, err := utils.AskConfirmation(
		"Set as Active Profile",
		fmt.Sprintf("Make '%s' the active authentication profile?", profileName),
		"Yes", "No",
	)

	if err == nil && setAsActive {
		if ActiveProfile != "" {
			AuthProfiles[ActiveProfile].Active = false
			if err := saveAuthProfile(AuthProfiles[ActiveProfile]); err != nil {
				utils.ShowWarning("Failed to update previous active profile")
			}
		}
		ActiveProfile = profileName
		AuthProfiles[ActiveProfile].Active = true
		if err := saveAuthProfile(AuthProfiles[ActiveProfile]); err != nil {
			utils.ShowWarning("Failed to save active profile status")
		}
		utils.ShowSuccess(fmt.Sprintf("'%s' is now the active profile", profileName))
	}

	askContinueOrReturnAuth()
}

func handleBearerTokenSetup(profile *AuthProfile) bool {
	// Get bearer token
	tokenConfig := utils.InputConfig{
		Title:       "Bearer Token:",
		Description: "Enter your bearer token",
		Placeholder: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		Required:    true,
		Password:    true,
	}

	token, err := utils.AskInput(tokenConfig)
	if err != nil {
		utils.ShowError("Error setting up bearer token", err)
		return false
	}

	profile.Token = token

	// Ask about expiry
	hasExpiry, err := utils.AskConfirmation(
		"Token Expiry",
		"Does this token have an expiry?",
		"Yes", "No",
	)

	if err != nil {
		utils.ShowError("Error asking about expiry", err)
		return false
	}

	if hasExpiry {
		expiryConfig := utils.InputConfig{
			Title:       "Token Expiry (optional):",
			Description: "Format: YYYY-MM-DD HH:MM or leave empty",
			Placeholder: "2024-12-31 23:59",
			Required:    false,
		}

		expiryStr, err := utils.AskInput(expiryConfig)
		if err == nil && expiryStr != "" {
			if expiry, parseErr := utils.ParseTimeInput(expiryStr); parseErr == nil && expiry != nil {
				profile.Expiry = expiry
			} else {
				utils.ShowWarning("Invalid date format, expiry not set")
			}
		}
	}

	return true
}

func handleAPIKeySetup(profile *AuthProfile) bool {
	configs := []utils.InputConfig{
		{
			Title:       "API Key:",
			Description: "Enter your API key",
			Placeholder: "sk_live_51H...",
			Required:    true,
			Password:    true,
		},
		{
			Title:       "Header Name:",
			Description: "The header name for the API key",
			Placeholder: "X-API-Key",
			Required:    false,
		},
	}

	values, err := utils.AskMultipleInputs(configs)
	if err != nil {
		utils.ShowError("Error setting up API key", err)
		return false
	}

	profile.APIKey = values[0]
	if values[1] == "" {
		profile.Header = "X-API-Key"
	} else {
		profile.Header = values[1]
	}

	return true
}

func handleBasicAuthSetup(profile *AuthProfile) bool {
	configs := []utils.InputConfig{
		{
			Title:       "Username:",
			Placeholder: "admin",
			Required:    true,
		},
		{
			Title:    "Password:",
			Required: true,
			Password: true,
		},
	}

	values, err := utils.AskMultipleInputs(configs)
	if err != nil {
		utils.ShowError("Error setting up basic authentication", err)
		return false
	}

	profile.Username = values[0]
	profile.Password = values[1]
	return true
}

func handleSelectProfile() {
	if err := loadAuthProfiles(); err != nil {
		utils.ShowError("Error loading auth profiles", err)
		return
	}

	if len(AuthProfiles) == 0 {
		utils.ShowMessage("No auth profiles found. Create one first.")
		askContinueOrReturnAuth()
		return
	}

	options := utils.CreateOptionsFromMap(AuthProfiles, func(name string, profile *AuthProfile) string {
		label := fmt.Sprintf("%s (%s)", name, strings.ToUpper(profile.Type))
		if name == ActiveProfile {
			label += " [ACTIVE]"
		}
		return label
	})

	selectedProfile, err := utils.AskSelection("Select Active Profile:", options)
	if err != nil {
		utils.ShowError("Error selecting profile", err)
		return
	}

	// Update active status
	if ActiveProfile != "" {
		AuthProfiles[ActiveProfile].Active = false
		if err := saveAuthProfile(AuthProfiles[ActiveProfile]); err != nil {
			utils.ShowWarning("Failed to update previous active profile")
		}
	}

	ActiveProfile = selectedProfile
	AuthProfiles[ActiveProfile].Active = true
	if err := saveAuthProfile(AuthProfiles[ActiveProfile]); err != nil {
		utils.ShowWarning("Failed to save active profile status")
	}

	utils.ShowSuccess(fmt.Sprintf("'%s' is now the active profile", selectedProfile))
	askContinueOrReturnAuth()
}

func handleEditProfile() {
	if err := loadAuthProfiles(); err != nil {
		utils.ShowError("Error loading auth profiles", err)
		return
	}

	if len(AuthProfiles) == 0 {
		utils.ShowMessage("No auth profiles found. Create one first.")
		askContinueOrReturnAuth()
		return
	}

	options := utils.CreateOptionsFromMap(AuthProfiles, func(name string, profile *AuthProfile) string {
		return fmt.Sprintf("%s (%s)", name, strings.ToUpper(profile.Type))
	})

	selectedProfile, err := utils.AskSelection("Select Profile to Edit:", options)
	if err != nil {
		utils.ShowError("Error selecting profile to edit", err)
		return
	}

	profile := AuthProfiles[selectedProfile]

	// Edit based on auth type
	var success bool
	switch profile.Type {
	case "bearer":
		success = handleBearerTokenSetup(profile)
	case "apikey":
		success = handleAPIKeySetup(profile)
	case "basic":
		success = handleBasicAuthSetup(profile)
	}

	if success {
		if err := saveAuthProfile(profile); err != nil {
			utils.ShowError("Failed to save updated profile", err)
			return
		}
		utils.ShowSuccess(fmt.Sprintf("Profile '%s' updated successfully!", selectedProfile))
	}
	askContinueOrReturnAuth()
}

func handleDeleteProfile() {
	if err := loadAuthProfiles(); err != nil {
		utils.ShowError("Error loading auth profiles", err)
		return
	}

	if len(AuthProfiles) == 0 {
		utils.ShowMessage("No auth profiles found.")
		askContinueOrReturnAuth()
		return
	}

	options := utils.CreateOptionsFromMap(AuthProfiles, func(name string, profile *AuthProfile) string {
		return fmt.Sprintf("%s (%s)", name, strings.ToUpper(profile.Type))
	})

	selectedProfile, err := utils.AskSelection("Select Profile to Delete:", options)
	if err != nil {
		utils.ShowError("Error selecting profile to delete", err)
		return
	}

	confirmDelete, err := utils.AskDangerousConfirmation(
		"Delete Profile",
		"Are you sure you want to delete profile",
		selectedProfile,
	)

	if err != nil {
		utils.ShowError("Error confirming deletion", err)
		return
	}

	if confirmDelete {
		// Delete from disk first
		if err := deleteAuthProfileFile(selectedProfile); err != nil {
			utils.ShowError("Failed to delete profile file", err)
			return
		}

		// Then delete from memory
		delete(AuthProfiles, selectedProfile)
		if ActiveProfile == selectedProfile {
			ActiveProfile = ""
		}
		utils.ShowSuccess(fmt.Sprintf("Profile '%s' deleted successfully!", selectedProfile))
	} else {
		utils.ShowMessage("Deletion cancelled.")
	}

	askContinueOrReturnAuth()
}

func handleViewProfiles() {
	if err := loadAuthProfiles(); err != nil {
		utils.ShowError("Error loading auth profiles", err)
		return
	}

	if len(AuthProfiles) == 0 {
		utils.ShowMessage("📋 No auth profiles configured.")
		askContinueOrReturnAuth()
		return
	}

	var profilesText strings.Builder
	profilesText.WriteString("📋 Authentication Profiles:\n")
	profilesText.WriteString("─────────────────────────────────\n\n")

	for name, profile := range AuthProfiles {
		status := ""
		if name == ActiveProfile {
			status = " [ACTIVE] "
		}

		profilesText.WriteString(fmt.Sprintf("🔑 %s%s\n", name, status))
		profilesText.WriteString(fmt.Sprintf("   Type: %s\n", strings.ToUpper(profile.Type)))

		switch profile.Type {
		case "bearer":
			profilesText.WriteString(utils.FormatKeyValue("Token", profile.Token, true))
			profilesText.WriteString("\n")
			if profile.Expiry != nil {
				profilesText.WriteString(utils.FormatKeyValue("Expires", utils.FormatTimeForDisplay(profile.Expiry), false))
				profilesText.WriteString("\n")
			}
		case "apikey":
			profilesText.WriteString(utils.FormatKeyValue("API Key", profile.APIKey, true))
			profilesText.WriteString("\n")
			profilesText.WriteString(utils.FormatKeyValue("Header", profile.Header, false))
			profilesText.WriteString("\n")
		case "basic":
			profilesText.WriteString(utils.FormatKeyValue("Username", profile.Username, false))
			profilesText.WriteString("\n")
			profilesText.WriteString(utils.FormatKeyValue("Password", profile.Password, true))
			profilesText.WriteString("\n")
		}
		profilesText.WriteString("\n")
	}

	profilesText.WriteString("─────────────────────────────────")

	utils.DisplayFormattedText("Authentication Profiles", profilesText.String())
	askContinueOrReturnAuth()
}

func askContinueOrReturnAuth() {
	utils.AskContinueOrReturn(
		HandleAuthenticationManagement,
		RunInteractiveMode,
		"Continue with Authentication",
		"Return to Main Menu",
	)
}

// GetActiveAuthProfile returns the currently active auth profile
func GetActiveAuthProfile() *AuthProfile {
	// Load profiles if not already loaded
	if len(AuthProfiles) == 0 {
		if err := loadAuthProfiles(); err != nil {
			return nil
		}
	}

	if ActiveProfile == "" {
		return nil
	}
	return AuthProfiles[ActiveProfile]
}

// GetAllAuthProfiles returns all auth profiles
func GetAllAuthProfiles() map[string]*AuthProfile {
	// Load profiles if not already loaded
	if len(AuthProfiles) == 0 {
		if err := loadAuthProfiles(); err != nil {
			return make(map[string]*AuthProfile)
		}
	}
	return AuthProfiles
}

