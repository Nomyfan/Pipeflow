package pipeflow

import (
	"fmt"
	"net/http"
)

// NotFoundMiddleware handles not found case
func NotFoundMiddleware(ctx HTTPContext) {
	if prop := ctx.Props["not_found_reason"]; prop != nil {
		reason := prop.(string)
		ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(ctx.ResponseWriter, "<h1 style='color: #171717;'>404 Not Found</h1> <h3 style='color: #171717;'>Cannot handle the request path <span style='color: #d82b21;'>"+ctx.Request.RequestURI+"</span></h3><p>Message: "+reason+"</p>")
	}
}
