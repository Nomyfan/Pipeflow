package core

import (
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
	Path    string
	Methods map[HttpMethod]bool
	Handle  Handler
}

// Handler's path equals to other's and HTTP methods have intersection
func (r *HttpHandler) Conflict(other *HttpHandler) bool {
	if r.Path == other.Path {
		for k := range r.Methods {
			if _, ok := other.Methods[k]; ok {
				return true
			}
		}
	}

	return false
}

func (r *HttpHandler) Match(path string, method string) bool {
	if r.Path != path {
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
		// Have conflict means matched
		return r.Conflict(&HttpHandler{Path: path, Methods: map[HttpMethod]bool{httpMethods[httpMethod]: true}})
	}

	return false
}
