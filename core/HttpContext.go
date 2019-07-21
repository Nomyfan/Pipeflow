package core

import "net/http"

// HttpContext is the wrapper for http request and response writer
type HttpContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}
