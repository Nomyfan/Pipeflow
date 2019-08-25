package pipeflow

import (
	"strings"
)
import "net/http"

type Flow struct {
	cors       RunnableMiddleware
	handlers   []HttpHandler
	middleware []Middleware
	dispatcher *HttpHandlerDispatcher
}

func (flow Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	ctx := HttpContext{Request: res, ResponseWriter: writer}
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

func NewFlow() Flow {
	flow := Flow{}
	flow.handlers = []HttpHandler{}
	flow.middleware = []Middleware{}
	flow.dispatcher = &HttpHandlerDispatcher{Handlers: &flow.handlers}

	return flow
}

func (flow *Flow) Use(middleware Middleware) {
	// Register middleware
	if nil != middleware {
		flow.middleware = append(flow.middleware, middleware)
	}
}

type runnableMiddleware struct {
	Handler func(ctx HttpContext)
}

func (rm *runnableMiddleware) Handle(ctx HttpContext) bool {
	rm.Handler(ctx)
	return true
}

// Run | Runnable middleware always returns true
func (flow *Flow) Run(middleware RunnableMiddleware) {
	if nil != middleware {
		flow.Use(&runnableMiddleware{Handler: middleware.Handle})
	}
}

func (flow *Flow) UseCors(origins []string, methods []string, headers []string, expose []string) {
	cors := Cors{AllowedOrigins: map[string]bool{}, AllowedMethods: methods, AllowedHeaders: headers, ExposedHeaders: expose}
	if nil != origins {
		for _, v := range origins {
			cors.AllowedOrigins[v] = true
		}
	}
	flow.cors = &cors
}

func (flow *Flow) Register(path string, handler Handler, methods []HttpMethod) error {
	path = strings.Trim(path, " ")
	if "" == path || path[0] != '/' || nil == methods || len(methods) == 0 || nil == handler {
		return BasicError{Message: "Args given are not valid"}
	}

	route, err := BuildRoute(path)
	if err != nil {
		return err
	}

	httpHandler := HttpHandler{Route: &route, Handle: handler, Methods: map[HttpMethod]bool{}}
	for _, v := range methods {
		httpHandler.Methods[v] = true
	}

	if flow.checkConflict(&httpHandler) {
		return BasicError{Message: "This handler conflicts with existing one"}
	}

	flow.appendHandler(httpHandler)
	return nil
}

func (flow *Flow) checkConflict(handler *HttpHandler) bool {
	for _, h := range flow.handlers {
		if h.Conflict(handler) {
			return true
		}
	}

	return false
}

func (flow *Flow) appendHandler(handler HttpHandler) {
	flow.handlers = append(flow.handlers, handler)
	flow.dispatcher.Handlers = &flow.handlers
}
