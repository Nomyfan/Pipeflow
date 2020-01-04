package pipeflow

import (
	"errors"
	"reflect"
	"strings"
)
import "net/http"

// Flow is main service register center
type Flow struct {
	cors         func(ctx HTTPContext)
	notfound     func(ctx HTTPContext)
	handlers     []RequestHandler
	middleware   []func(ctx HTTPContext, next func())
	dispatcher   *HTTPRequestDispatcher
	resource     map[string]interface{}
	resourceType map[reflect.Type]interface{}
}

func (flow *Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	ctx := HTTPContext{Request: res, ResponseWriter: writer, resource: flow.resource, resourceType: flow.resourceType, Props: map[string]interface{}{}}

	// Add CORS to the pipeline
	if flow.cors != nil {
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			flow.cors(ctx)
			next()
		})
	}

	// Add HTTP dispatcher
	if flow.dispatcher != nil {
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			if err := flow.dispatcher.Handle(ctx); err != nil && flow.notfound != nil {
				ctx.Props["crash_reason"] = err.Error()
				flow.notfound(ctx)
			}
		})
	}

	invoke(flow, ctx, 0)
}

func invoke(f *Flow, ctx HTTPContext, i int) {
	if i == len(f.middleware) {
		return
	}
	f.middleware[i](ctx, func() {
		invoke(f, ctx, i+1)
	})
}

// NewFlow returns a new instance of pipeflow
func NewFlow() *Flow {
	flow := Flow{}
	flow.handlers = []RequestHandler{}
	flow.middleware = []func(ctx HTTPContext, next func()){}
	flow.dispatcher = &HTTPRequestDispatcher{Handlers: &flow.handlers}
	flow.resource = map[string]interface{}{}
	flow.resourceType = map[reflect.Type]interface{}{}
	flow.notfound = NotFoundMiddleware

	return &flow
}

// Use registers middleware
func (flow *Flow) Use(middleware func(ctx HTTPContext, next func())) {
	if middleware != nil {
		flow.middleware = append(flow.middleware, middleware)
	}
}

// Run runnable typed middleware will always invoke next
func (flow *Flow) Run(middleware func(ctx HTTPContext)) {
	if nil != middleware {
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			middleware(ctx)
			next()
		})
	}
}

// UseCors registers CORS middleware
func (flow *Flow) UseCors(origins []string, methods []string, headers []string, expose []string) {
	cors := Cors{AllowedOrigins: map[string]bool{}, AllowedMethods: methods, AllowedHeaders: headers, ExposedHeaders: expose}
	if nil != origins {
		for _, v := range origins {
			cors.AllowedOrigins[v] = true
		}
	}
	flow.cors = func(ctx HTTPContext) {
		cors.Handle(ctx)
	}
}

// Map is used to add request handler
func (flow *Flow) Map(path string, handler func(ctx HTTPContext), methods []HTTPMethod) {
	path = strings.Trim(path, " ")
	if "" == path || path[0] != '/' || nil == methods || len(methods) == 0 || nil == handler {
		panic(errors.New("args given are not valid"))
	}

	route, err := BuildRoute(path)
	if err != nil {
		panic(err)
	}

	httpHandler := RequestHandler{Route: &route, Handle: handler, Methods: map[HTTPMethod]bool{}}
	for _, v := range methods {
		httpHandler.Methods[v] = true
	}

	if flow.checkConflict(&httpHandler) {
		panic(errors.New("this handler conflicts with the existing one"))
	}

	flow.appendHandler(httpHandler)
}

func (flow *Flow) GET(path string, handler func(ctx HTTPContext)) {
	flow.Map(path, handler, []HTTPMethod{HTTPGet})
}

func (flow *Flow) POST(path string, handler func(ctx HTTPContext)) {
	flow.Map(path, handler, []HTTPMethod{HTTPPost})
}

// SetResource sets global singleton resource
func (flow *Flow) SetResource(key string, value interface{}) {
	flow.resource[key] = value
}

// SetResourceWithType sets global singleton resource using it's type as key
func (flow *Flow) SetResourceWithType(key reflect.Type, value interface{}) {
	flow.resourceType[key] = value
}

// SetResourceAlsoWithType calls SetResource and SetResourceWithType
func (flow *Flow) SetResourceAlsoWithType(key string, value interface{}) {
	flow.SetResource(key, value)
	flow.SetResourceWithType(reflect.TypeOf(value), value)
}

func (flow *Flow) checkConflict(handler *RequestHandler) bool {
	for _, h := range flow.handlers {
		if h.Conflict(handler) {
			return true
		}
	}

	return false
}

func (flow *Flow) appendHandler(handler RequestHandler) {
	flow.handlers = append(flow.handlers, handler)
	flow.dispatcher.Handlers = &flow.handlers
}
