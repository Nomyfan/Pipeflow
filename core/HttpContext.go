package core

import "net/http"

type HttpContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Vars           *map[string]string
}
