package pipeflow

import "net/http"

// HTTPContext is the request context wrapper
type HTTPContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Vars           map[string]string
	resource       map[string]interface{}
	Props          map[string]interface{}
}

// GetResource get global singleton resource preset
func (ctx HTTPContext) GetResource(key string) interface{} {
	return ctx.resource[key]
}
