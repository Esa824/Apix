package model

// HTTPResponse holds the response data and provides querying capabilities
type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    map[string]string
	Body       []byte
	IsJSON     bool
	ParsedJSON any
}
