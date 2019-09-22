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
	matchMsg := ""
	for _, v := range *hd.Handlers {
		if v.MatchPath(ctx.Request) {
			if v.MatchMethod(ctx.Request) {
				getPathVars(&v, &ctx)
				v.Handle(ctx)
				return
			} else {
				matchMsg = "Request method doesn't match"
			}
			break
		} else {
			matchMsg = "Path doesn't match"
		}
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(ctx.ResponseWriter, "<h1 style='color: #171717;'>404 Not Found</h1> <h3 style='color: #171717;'>Cannot handle the request path <span style='color: #d82b21;'>"+ctx.Request.RequestURI+"</span></h3><p>Message: "+matchMsg+"</p>")
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
