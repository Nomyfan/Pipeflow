package pipeflow

import (
	"errors"
	"regexp"
	"strings"
)

// Route is used to identify a route
type Route struct {
	PathReg *regexp.Regexp
	Params  map[string]bool
	Vars    map[string]bool
}

// buildRoute builds a route from given pattern
func buildRoute(pattern string) (Route, error) {
	route := Route{}
	err := parse(pattern, &route)

	return route, err
}

func parse(pattern string, route *Route) error {
	if len(pattern) == 0 {
		return errors.New("path should not be empty")
	}

	if "/" == pattern {
		route.PathReg = regexp.MustCompile("^/$")
		return nil
	}

	route.Params = make(map[string]bool, 0)
	route.Vars = make(map[string]bool, 0)
	var routePattern = "^"

	pathReg := regexp.MustCompile(`^[\w\x{4e00}-\x{9fa5}]+$`)
	varReg := regexp.MustCompile(`^{([\w\x{4e00}-\x{9fa5}]+)}$`)
	paramReg := regexp.MustCompile(`^(?P<lp>[\w\x{4e00}-\x{9fa5}]+)\?([\w\x{4e00}-\x{9fa5}]+=\?+&?)*`)
	kvReg := regexp.MustCompile(`(?P<key>[\w\x{4e00}-\x{9fa5}]+)=\?`)

	var parts []string
	if pattern[0] == '/' {
		parts = strings.Split(pattern, "/")[1:]
	} else {
		parts = strings.Split(pattern, "/")
	}
	for i, v := range parts {
		if len(v) == 0 {
			return errors.New("partial path cannot be empty")
		} else if pathReg.MatchString(v) {
			routePattern += "/" + v
		} else if varReg.MatchString(v) {
			routePattern += `/(?P<` + v[1:len(v)-1] + `>[^/]+)`
			route.Vars[v[1:len(v)-1]] = true
		} else if paramReg.MatchString(v) {
			if i != len(parts)-1 {
				return errors.New("params should be in the last")
			}
			routePattern += `/` + paramReg.FindStringSubmatch(v)[1]
			// Add params into map
			for _, m := range kvReg.FindAllStringSubmatch(v, -1) {
				route.Params[m[1]] = true
			}
		} else {
			return errors.New("invalid URL was give")
		}
	}

	routePattern += "/?$"
	pathReg, err := regexp.Compile(routePattern)
	if err == nil {
		route.PathReg = pathReg
	}
	return err
}

// equals checks whether two routes are equaled
func (route *Route) equals(other *Route) bool {
	if route.PathReg.String() != other.PathReg.String() {
		return false
	}

	if len(route.Params) != len(route.Params) {
		return false
	}

	for k := range route.Params {
		if ok := other.Params[k]; !ok {
			return false
		}
	}

	if len(route.Vars) != len(other.Vars) {
		return false
	}

	for _, rv := range route.Vars {
		contains := false
		for _, ov := range other.Vars {
			if rv == ov {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	return true
}
