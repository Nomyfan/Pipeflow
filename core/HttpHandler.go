package core

import (
	"net/http"
	"strings"
)

type HttpMethod int

type Handler func(ctx HttpContext)

const (
	HttpGet = iota
	HttpHead
	HttpPost
	HttpPut
	HttpDelete
	HttpConnect
	HttpOptions
	HttpTrace
)

type HttpHandler struct {
	Route   *Route
	Methods map[HttpMethod]bool
	Handle  Handler
}

// Handler's path equals to other's and HTTP methods have intersection
func (h *HttpHandler) Conflict(other *HttpHandler) bool {
	if h.Route.Equals(other.Route) {
		return h.HasInterMethod(other)
	}

	return false
}

func (h *HttpHandler) HasInterMethod(other *HttpHandler) bool {
	for k := range h.Methods {
		if _, ok := other.Methods[k]; ok {
			return true
		}
	}

	return false
}

func (h *HttpHandler) Match(request *http.Request) bool {
	path := request.URL.Path
	method := request.Method

	if !h.Route.PathReg.MatchString(path) {
		return false
	}

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
	httpMethods := []HttpMethod{HttpGet, HttpHead, HttpPost, HttpPut, HttpDelete, HttpConnect, HttpOptions, HttpTrace}

	method = strings.ToUpper(method)
	httpMethod := -1
	for i, v := range methods {
		if v == method {
			httpMethod = i
			break
		}
	}

	if -1 != httpMethod {
		hasInter := h.HasInterMethod(&HttpHandler{Methods: map[HttpMethod]bool{httpMethods[httpMethod]: true}})
		if !hasInter {
			return false
		}
	}

	if e := request.ParseForm(); e != nil {
		return false
	} else {
		for k := range h.Route.Params {
			if _, ok := request.Form[k]; !ok {
				return false
			}
		}
	}

	return true
}
