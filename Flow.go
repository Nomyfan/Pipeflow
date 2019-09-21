package pipeflow

import (
	"strings"
)
import "net/http"

// Flow is main service register center
type Flow struct {
	cors       func(ctx HTTPContext)
	handlers   []RequestHandler
	middleware []func(ctx HTTPContext, next func())
	dispatcher *HTTPRequestDispatcher
	resource   map[string]interface{}
}

func (flow Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	ctx := HTTPContext{Request: res, ResponseWriter: writer, resource: flow.resource}

	// Add CORS to the pipeline
	if flow.cors != nil {
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			flow.cors(ctx)
			next()
		})
	}

	// Add dispatcher to the end of pipeline
	if flow.dispatcher != nil {
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			flow.dispatcher.Handle(ctx)
		})
	}

	invoke(&flow, ctx, 0)
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
func NewFlow() Flow {
	flow := Flow{}
	flow.handlers = []RequestHandler{}
	flow.middleware = []func(ctx HTTPContext, next func()){}
	flow.dispatcher = &HTTPRequestDispatcher{Handlers: &flow.handlers}
	flow.resource = map[string]interface{}{}

	return flow
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

// Register is used to add request handler
func (flow *Flow) Register(path string, handler func(ctx HTTPContext), methods []HTTPMethod) error {
	path = strings.Trim(path, " ")
	if "" == path || path[0] != '/' || nil == methods || len(methods) == 0 || nil == handler {
		return BasicError{Message: "Args given are not valid"}
	}

	route, err := BuildRoute(path)
	if err != nil {
		return err
	}

	httpHandler := RequestHandler{Route: &route, Handle: handler, Methods: map[HTTPMethod]bool{}}
	for _, v := range methods {
		httpHandler.Methods[v] = true
	}

	if flow.checkConflict(&httpHandler) {
		return BasicError{Message: "This handler conflicts with existing one"}
	}

	flow.appendHandler(httpHandler)
	return nil
}

// SetResource set global singleton resource
func (flow Flow) SetResource(key string, value interface{}) {
	flow.resource[key] = value
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
