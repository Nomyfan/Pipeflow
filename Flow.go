package pipeflow

import (
	"strings"
)
import "net/http"

// Flow is main service register center
type Flow struct {
	cors       RunnableMiddleware
	handlers   []RequestHandler
	middleware []Middleware
	dispatcher *HTTPRequestDispatcher
}

func (flow Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	ctx := HTTPContext{Request: res, ResponseWriter: writer}
	toBreak := false

	for _, v := range flow.middleware {
		if !v.Handle(ctx) {
			toBreak = true
			break
		}
	}

	if nil != flow.cors {
		flow.cors.Handle(ctx)
	}

	if !toBreak {
		flow.dispatcher.Handle(ctx)
	}
}

// NewFlow returns a new instance of pipeflow
func NewFlow() Flow {
	flow := Flow{}
	flow.handlers = []RequestHandler{}
	flow.middleware = []Middleware{}
	flow.dispatcher = &HTTPRequestDispatcher{Handlers: &flow.handlers}

	return flow
}

// Use registers middleware
func (flow *Flow) Use(middleware Middleware) {
	// Register middleware
	if nil != middleware {
		flow.middleware = append(flow.middleware, middleware)
	}
}

type runnableMiddleware struct {
	Handler func(ctx HTTPContext)
}

func (rm *runnableMiddleware) Handle(ctx HTTPContext) bool {
	rm.Handler(ctx)
	return true
}

// Run | Runnable middleware always returns true
func (flow *Flow) Run(middleware RunnableMiddleware) {
	if nil != middleware {
		flow.Use(&runnableMiddleware{Handler: middleware.Handle})
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
	flow.cors = &cors
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
