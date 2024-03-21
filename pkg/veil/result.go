package veil

import "net/http"

type Result struct {
	url        string
	statusCode int
	headers    http.Header
	body       string
}
