package pipeflow

import (
	"github.com/pkg/errors"
	"strings"
)

type HTTPRequestDispatcher interface {
	Map(path string, handler func(ctx HTTPContext), methods ...HTTPMethod)
	Handle(ctx HTTPContext)
}

type defaultHTTPRequestDispatcher struct {
	handlers []requestHandler
}

func newDefaultHTTPRequestDispatcher() *defaultHTTPRequestDispatcher {
	return &defaultHTTPRequestDispatcher{handlers: []requestHandler{}}
}

func (m *defaultHTTPRequestDispatcher) Map(path string, handler func(ctx HTTPContext), methods ...HTTPMethod) {
	path = strings.Trim(path, " ")
	if len(path) == 0 || path[0] != '/' {
		panic(errors.New("path should starts with '/'"))
	}

	if len(methods) == 0 {
		panic(errors.New("no HTTP method was specified"))
	}

	if handler == nil {
		panic(errors.New("handler cannot be nil"))
	}

	route, err := buildRoute(path)
	if err != nil {
		panic(err)
	}

	httpHandler := requestHandler{route: &route, handle: handler, methods: map[HTTPMethod]bool{}}
	for _, v := range methods {
		httpHandler.methods[v] = true
	}

	if checkConflict(m.handlers, &httpHandler) {
		panic(errors.New("this handler conflicts with the existing one"))
	}

	appendHandler(m, httpHandler)
}

func (m *defaultHTTPRequestDispatcher) Handle(ctx HTTPContext) {

	reason := "path dose not match"
	for _, h := range m.handlers {
		if matchPath(&h, ctx.Request) {
			if matchMethod(&h, ctx.Request.Method) {
				fillPathVars(&h, &ctx)
				h.handle(ctx)
				return
			}
			reason = "HTTP method dose not match"
		}
	}
	ctx.Props["not_found_reason"] = reason
}

func checkConflict(handlers []requestHandler, handler *requestHandler) bool {
	for _, h := range handlers {
		if conflict(&h, handler) {
			return true
		}
	}
	return false
}

func appendHandler(dp *defaultHTTPRequestDispatcher, handler requestHandler) {
	dp.handlers = append(dp.handlers, handler)
}

func fillPathVars(handler *requestHandler, ctx *HTTPContext) {

	segments := splitPathIntoSegments(ctx.Request.URL.Path)
	ctx.Vars = map[string]string{}
	for i, v := range segments {
		if handler.route.segments[i].isVar {
			ctx.Vars[handler.route.segments[i].seg] = v
		}
	}
}
