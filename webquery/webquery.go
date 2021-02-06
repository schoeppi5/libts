package webquery

// This file contains the basic structs and helper funcs

import (
	"net/http"
)

// WebQuery stores info for webquery connections
type WebQuery struct {
	Host       string
	Port       int
	Key        string
	TLS        bool
	HTTPClient *http.Client
}
