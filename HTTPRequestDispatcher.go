package pipeflow

import (
	"github.com/pkg/errors"
	"strings"
)

type HTTPRequestDispatcher interface {
	Map(path string, handler func(ctx HTTPContext), methods []HTTPMethod)
	Handle(ctx HTTPContext)
}

type DefaultHTTPRequestDispatcher struct {
	handlers []RequestHandler
}

func NewDefaultRequestDispatcher() *DefaultHTTPRequestDispatcher {
	return &DefaultHTTPRequestDispatcher{handlers: []RequestHandler{}}
}

func (m *DefaultHTTPRequestDispatcher) Map(path string, handler func(ctx HTTPContext), methods []HTTPMethod) {
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

	if checkConflict(m.handlers, &httpHandler) {
		panic(errors.New("this handler conflicts with the existing one"))
	}

	appendHandler(m, httpHandler)
}

func (m *DefaultHTTPRequestDispatcher) Handle(ctx HTTPContext) {
	for _, v := range m.handlers {
		if v.MatchPath(ctx.Request) {
			if v.MatchMethod(ctx.Request) {
				getPathVars(&v, &ctx)
				v.Handle(ctx)
				return
			} else {
				ctx.Props["not_found_reason"] = "HTTP method dose not match"
				return
			}
		}
	}
	ctx.Props["not_found_reason"] = "path dose not match"
}

func checkConflict(handlers []RequestHandler, handler *RequestHandler) bool {
	for _, h := range handlers {
		if h.Conflict(handler) {
			return true
		}
	}
	return false
}

func appendHandler(dp *DefaultHTTPRequestDispatcher, handler RequestHandler) {
	dp.handlers = append(dp.handlers, handler)
}

func getPathVars(handler *RequestHandler, ctx *HTTPContext) {
	regex := handler.Route.PathReg

	vars := regex.FindAllStringSubmatch(ctx.Request.URL.Path, -1)
	groupNames := regex.SubexpNames()
	ctx.Vars = map[string]string{}

	for i, name := range groupNames {
		if _, ok := handler.Route.Vars[name]; ok && name != "" {
			ctx.Vars[name] = vars[0][i]
		}
	}
}
