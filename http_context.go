package pipeflow

import (
	"net/http"
	"reflect"
)

// HTTPContext is the request context wrapper
type HTTPContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Vars           map[string]string
	resource       map[string]interface{}
	resourceType   map[reflect.Type]interface{}
	Props          map[string]interface{}
}

// GetResource gets global singleton resource preset
func (ctx HTTPContext) GetResource(key string) interface{} {
	return ctx.resource[key]
}

// GetResourceByType gets global singleton resource preset by type
func (ctx HTTPContext) GetResourceByType(key reflect.Type) interface{} {
	return ctx.resourceType[key]
}
