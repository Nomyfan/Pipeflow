package pipeflow

import (
	"strings"
)

// Cors is used to handle CORS
type Cors struct {
	AllowedOrigins map[string]bool
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
}

// handle implements middleware
func (cors *Cors) Handle(ctx HTTPContext) {
	// Cors runs after any middleware and before dispatcher middleware
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
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-methods", strings.Join(cors.AllowedMethods, ","))
		}
		if nil != cors.AllowedHeaders && 0 != len(cors.AllowedHeaders) {
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ","))
		}
		if nil != cors.ExposedHeaders && 0 != len(cors.ExposedHeaders) {
			ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposedHeaders, ","))
		}
	}
}
