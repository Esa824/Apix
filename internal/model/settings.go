package model

import "time"

// DisplaySettings manages output formatting preferences
type DisplaySettings struct {
	ResponseFormat  string // "pretty-json", "raw", "headers-only", "compact-json"
	ColorOutput     bool
	ShowTiming      bool
	ShowHeaders     bool
	ShowStatusCode  bool
	MaxResponseSize int // in KB, 0 = unlimited
	IndentSize      int // for pretty formatting
	LineNumbers     bool
	SyntaxHighlight bool
}

// BehaviorSettings manages application behavior
type BehaviorSettings struct {
	AutoSaveRequests       bool
	ConfirmDeleteRequests  bool
	ConfirmDestructive     bool
	RequestTimeout         int // in seconds
	MaxRetries             int
	RetryDelay             int // in seconds
	FollowRedirects        bool
	MaxRedirects           int
	ValidateSSL            bool
	CacheResponses         bool
	CacheDuration          int // in minutes
	ShowProgressBar        bool
	VerboseMode            bool
	SaveFailedRequests     bool
	AutoAddHeaders         bool
	DefaultContentType     string
	PreserveSessionCookies bool
}

// NetworkSettings manages connection preferences
type NetworkSettings struct {
	DefaultTimeout     int // in seconds
	ConnectTimeout     int // in seconds
	ReadTimeout        int // in seconds
	WriteTimeout       int // in seconds
	MaxConnections     int
	UserAgent          string
	ProxyURL           string
	ProxyEnabled       bool
	KeepAlive          bool
	CompressionEnabled bool
}

// LoggingSettings manages request/response logging
type LoggingSettings struct {
	EnableLogging  bool
	LogLevel       string // "debug", "info", "warn", "error"
	LogFile        string
	LogRequests    bool
	LogResponses   bool
	LogHeaders     bool
	LogTiming      bool
	RotateLogFiles bool
	MaxLogSize     int // in MB
	MaxLogFiles    int
}

// GlobalSettings holds all application settings
type GlobalSettings struct {
	Display   DisplaySettings
	Behavior  BehaviorSettings
	Network   NetworkSettings
	Logging   LoggingSettings
	Version   string
	LastSaved time.Time
}
