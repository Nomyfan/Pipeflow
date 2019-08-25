package pipeflow

import (
	"fmt"
)

// HttpHandlerDispatcher is a middleware in the end point of workflow
type HttpHandlerDispatcher struct {
	Handlers *[]HttpHandler
}

func (hd *HttpHandlerDispatcher) Handle(ctx HttpContext) bool {
	for _, v := range *hd.Handlers {
		if v.Match(ctx.Request) {
			getPathVars(&v, &ctx)
			v.Handle(ctx)
			return false
		}
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(ctx.ResponseWriter, "<h1>404</h1> <h3>The request path <span style='color: red;'>"+ctx.Request.RequestURI+"</span> cannot be found</h3>")
	return false
}

func getPathVars(handler *HttpHandler, ctx *HttpContext) {
	regex := handler.Route.PathReg

	vars := regex.FindAllStringSubmatch(ctx.Request.URL.Path, -1)
	groupNames := regex.SubexpNames()
	ctx.Vars = &map[string]string{}

	for i, name := range groupNames {
		if _, ok := handler.Route.Vars[name]; ok && name != "" {
			(*ctx.Vars)[name] = vars[0][i]
		}
	}
}
