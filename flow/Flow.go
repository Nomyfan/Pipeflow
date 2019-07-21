package flow

import (
	"pipeflow/core"
	"strings"
)
import "net/http"

type Flow struct {
	handlers   []core.HttpHandler
	middleware []core.Middleware
	dispatcher *core.HttpHandlerDispatcher
}

func (flow Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	ctx := core.HttpContext{Request: res, ResponseWriter: writer}
	breaker := false
	for _, v := range flow.middleware {
		if !v.Handle(ctx) {
			breaker = true
			break
		}
	}
	if !breaker {
		flow.dispatcher.Handle(ctx)
	}
}

func NewFlow() Flow {
	flow := Flow{}
	flow.handlers = []core.HttpHandler{}
	flow.middleware = []core.Middleware{}
	flow.dispatcher = &core.HttpHandlerDispatcher{Handlers: &flow.handlers}

	return flow
}

func (flow *Flow) Use(middleware core.Middleware) {
	// Register middleware
	if nil != middleware {
		flow.middleware = append(flow.middleware, middleware)
	}
}

type wrappedMiddleware struct {
	Handler func(ctx core.HttpContext)
}

func (wmw *wrappedMiddleware) Handle(ctx core.HttpContext) bool {
	wmw.Handler(ctx)
	return true
}

// Runnable middleware always returns true
func (flow *Flow) Run(middleware core.RunnableMiddleware) {
	// Register middleware
	if nil != middleware {
		flow.Use(&wrappedMiddleware{Handler: middleware.Handle})
	}
}

func (flow *Flow) AddHandler(path string, handler core.Handler, methods []core.HttpMethod) bool {
	path = strings.Trim(path, " ")
	if "" == path || path[0] != '/' || nil == methods || nil == handler {
		return false
	}

	httpHandler := core.HttpHandler{Path: path, Handle: handler, Methods: map[core.HttpMethod]bool{}}
	for _, v := range methods {
		httpHandler.Methods[v] = true
	}

	if !flow.checkConflict(&httpHandler) {
		flow.appendHandler(httpHandler)
		return true
	}

	return false
}

func (flow *Flow) checkConflict(handler *core.HttpHandler) bool {
	for _, v := range flow.handlers {
		if v.Conflict(handler) {
			return true
		}
	}

	return false
}

func (flow *Flow) appendHandler(handler core.HttpHandler) {
	// Append handler and update dispatcher's handlers ref
	flow.handlers = append(flow.handlers, handler)
	flow.dispatcher.Handlers = &flow.handlers
}
