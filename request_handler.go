package pipeflow

import (
	"net/http"
	"strings"
)

// HTTPMethod is enum of http methods
type HTTPMethod int

const (
	// HTTPGet GET
	HTTPGet = iota
	// HTTPHead HEAD
	HTTPHead
	// HTTPPost POST
	HTTPPost
	// HTTPPut PUT
	HTTPPut
	// HTTPDelete DELETE
	HTTPDelete
	// HTTPOptions OPTIONS
	HTTPOptions
	// HTTPTrace TRACE
	HTTPTrace
)

// requestHandler is used to register a request handler
type requestHandler struct {
	route   *route
	methods map[HTTPMethod]bool
	handle  func(ctx HTTPContext)
}

// conflict checks handler's path equals to other's and HTTP methods have intersection
func conflict(h *requestHandler, other *requestHandler) bool {
	if h.route.equals(other.route) {
		return hasInterMethod(h, other)
	}

	return false
}

// hasInterMethod checks whether http methods has intersection
func hasInterMethod(h *requestHandler, other *requestHandler) bool {
	for k := range h.methods {
		if _, ok := other.methods[k]; ok {
			return true
		}
	}

	return false
}

// matchPath checks whether request path is matched
func matchPath(h *requestHandler, request *http.Request) bool {

	segments := splitPathIntoSegments(request.URL.Path)
	if len(h.route.segments) != len(segments) {
		return false
	}
	for i, v := range segments {
		if !h.route.segments[i].isVar && h.route.segments[i].seg != v {
			return false
		}
	}

	if err := request.ParseForm(); err != nil {
		return false
	}
	if len(h.route.params) != len(request.Form) {
		// TODO
		//  or find the route with least params matching the request form
		return false
	}
	for k := range h.route.params {
		if _, ok := request.Form[k]; !ok {
			return false
		}
	}

	return true
}

// matchMethod checks whether request method is matched
func matchMethod(h *requestHandler, method string) bool {

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "TRACE"}
	httpMethods := []HTTPMethod{HTTPGet, HTTPHead, HTTPPost, HTTPPut, HTTPDelete, HTTPOptions, HTTPTrace}

	method = strings.ToUpper(method)
	httpMethod := -1
	for i, v := range methods {
		if v == method {
			httpMethod = i
			break
		}
	}

	if -1 != httpMethod {
		hasInter := hasInterMethod(h, &requestHandler{methods: map[HTTPMethod]bool{httpMethods[httpMethod]: true}})
		if !hasInter {
			return false
		}
	}

	return true
}
