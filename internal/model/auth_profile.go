package model

import "time"

// AuthProfile represents an authentication profile
type AuthProfile struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Token    string     `json:"token"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	APIKey   string     `json:"api_key"`
	Header   string     `json:"header"`
	Expiry   *time.Time `json:"expiry"`
	Active   bool       `json:"active"`
}
