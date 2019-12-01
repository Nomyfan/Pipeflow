package pipeflow

import "errors"

// HTTPRequestDispatcher is a middleware in the end point of workflow
type HTTPRequestDispatcher struct {
	Handlers *[]RequestHandler
}

// Handle implements middleware
func (hd *HTTPRequestDispatcher) Handle(ctx HTTPContext) error {
	methodNotMatch := false
	for _, v := range *hd.Handlers {
		if v.MatchPath(ctx.Request) {
			if v.MatchMethod(ctx.Request) {
				getPathVars(&v, &ctx)
				v.Handle(ctx)
				return nil
			} else {
				methodNotMatch = true
			}
		}
	}

	if methodNotMatch {
		return errors.New("HTTP method dose not match")
	}
	return errors.New("path dose not match")
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
