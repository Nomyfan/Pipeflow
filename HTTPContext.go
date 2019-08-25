package pipeflow

import "net/http"

// HTTPContext is the request context wrapper
type HTTPContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Vars           *map[string]string
}
