package pipeflow

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type segment struct {
	seg   string
	isVar bool
}

// route is used to identify a request URI
type route struct {
	params   map[string]bool
	vars     map[string]bool
	segments []segment
}

// buildRoute builds a route from given pattern
func buildRoute(pattern string) (route, error) {
	route := route{}
	err := parse(pattern, &route)

	return route, err
}

func parse(pattern string, route *route) error {
	// pattern example: /p1/{p2}/p3?p5&p6

	if len(pattern) == 0 {
		return errors.New("pathPattern should not be empty")
	}

	if "/" == pattern {
		return nil
	}

	route.params = make(map[string]bool, 0)
	route.vars = make(map[string]bool, 0)

	paramDelimiterIdx := len(pattern)
	fragmentDelimiterIdx := len(pattern) // ignore the fragments which means ignoring the part range in [fragmentDelimiterIdx, len(pattern)]
	for i, runeVal := range pattern {
		if string(runeVal) == "?" {
			paramDelimiterIdx = i
		} else if string(runeVal) == "#" {
			fragmentDelimiterIdx = i
		}
	}
	if fragmentDelimiterIdx < paramDelimiterIdx {
		panic(errors.New(fmt.Sprintf("invalid pattern: %s", pattern)))
	}

	pathPattern := pattern[1:paramDelimiterIdx]
	segments := strings.Split(pathPattern, "/")
	pathSegReg := regexp.MustCompile(`^[^{}/?:#\[\]@!$&'()*+,;=]+$`)
	varReg := regexp.MustCompile(`^{[^{}/?:#\[\]@!$&'()*+,;=]+}$`)
	for _, seg := range segments {
		if pathSegReg.MatchString(seg) {
			route.segments = append(route.segments, segment{seg: seg})
		} else if varReg.MatchString(seg) {
			v := seg[1 : len(seg)-1]
			route.vars[v] = true
			route.segments = append(route.segments, segment{seg: v, isVar: true})
		} else {
			panic(errors.New(fmt.Sprintf("pattern[%s] contains invalid path segment[%s]", pattern, seg)))
		}
	}

	if paramDelimiterIdx != fragmentDelimiterIdx {
		paramDelimiterIdx += 1
	}
	if paramPattern := pattern[paramDelimiterIdx:fragmentDelimiterIdx]; paramPattern != "" {
		params := strings.Split(paramPattern, "&")
		paramReg := regexp.MustCompile(`[^{}/?:#\[\]@!$&'()*+,;=]+`)
		for _, p := range params {
			if !paramReg.MatchString(p) {
				panic(fmt.Sprintf("pattern[%s] contains invalid formatted param[%s]", pattern, p))
			}
			route.params[p] = true
		}
	}

	return nil
}

// equals checks whether two route are equaled
func (r *route) equals(other *route) bool {

	if len(r.segments) != len(other.segments) {
		return false
	}

	if len(r.params) != len(other.params) {
		return false
	}

	for i := 0; i < len(r.segments); i++ {
		s1 := r.segments[i]
		s2 := other.segments[i]
		if s1.isVar != s2.isVar || (!s1.isVar && s1.seg != s2.seg) {
			return false
		}
	}

	for k := range r.params {
		if ok, _ := other.params[k]; !ok {
			return false
		}
	}

	return true
}
