package core

// HttpHandlerDispatcher is a middleware in the end point of workflow
type HttpHandlerDispatcher struct {
	Handlers *[]HttpHandler
}

func (hd *HttpHandlerDispatcher) Handle(ctx HttpContext) bool {
	path := ctx.Request.URL.Path
	for _, v := range *hd.Handlers {
		if v.Match(path, ctx.Request.Method) {
			v.Handle(ctx)
		}
	}

	// Here is the end point of flow
	return false
}
