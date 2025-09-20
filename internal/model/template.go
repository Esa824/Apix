package model

type Template struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
	Files       map[string]string `json:"files,omitempty"`
	Body        any               `json:"body,omitempty"`
}
