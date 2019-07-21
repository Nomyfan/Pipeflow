package middlewares

import (
	"pipeflow/core"
	"strings"
)

type Cors struct {
	AllowedOrigins map[string]bool
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
}

// Cors runs after any middleware and before dispatcher middleware
func (cors *Cors) Handle(ctx core.HttpContext) {

	origin := ctx.Request.Header.Get("Origin")
	enabled := false
	if _, ok := cors.AllowedOrigins["*"]; ok {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		enabled = true
	} else if _, ok := cors.AllowedOrigins[origin]; ok {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)
		enabled = true
	}
	if enabled {
		if nil != cors.AllowedMethods && 0 != len(cors.AllowedMethods) {
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.AllowedMethods, ","))
		}
		if nil != cors.AllowedHeaders && 0 != len(cors.AllowedHeaders) {
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ","))
		}
		if nil != cors.ExposedHeaders && 0 != len(cors.ExposedHeaders) {
			ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposedHeaders, ","))
		}
	}
}
