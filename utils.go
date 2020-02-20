package pipeflow

import "strings"

func splitPathIntoSegments(path string) []string {
	path = path[1:] // remove heading slash
	var segments []string
	if path[len(path)-1:] == "/" {
		// ensure there's not tailing slash
		path = path[0 : len(path)-1]
	}
	if path != "" {
		segments = strings.Split(path, "/")
	}

	return segments
}
