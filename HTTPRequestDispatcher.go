package pipeflow

import (
	"fmt"
	"net/http"
)

// HTTPRequestDispatcher is a middleware in the end point of workflow
type HTTPRequestDispatcher struct {
	Handlers *[]RequestHandler
}

// Handle implements middleware
func (hd *HTTPRequestDispatcher) Handle(ctx HTTPContext) {
	for _, v := range *hd.Handlers {
		if v.Match(ctx.Request) {
			getPathVars(&v, &ctx)
			v.Handle(ctx)
			return
		}
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(ctx.ResponseWriter, "<h1>404</h1> <h3>Cannot found the request path <span style='color: red;'>"+ctx.Request.RequestURI)
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
