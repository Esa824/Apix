package model

type Auth struct {
	Type      string // "bearer", "apikey", "basic"
	Primary   string // token, api key, or username
	Secondary string // password if needed
}
