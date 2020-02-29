package pipeflow

import "strings"

func splitPathIntoSegments(path string) []string {
	var segments []string
	if path[len(path)-1:] == "/" {
		// ensure there's not tailing slash
		path = path[0 : len(path)-1]
	}
	if path != "" {
		path = path[1:] // remove heading slash
		segments = strings.Split(path, "/")
	}

	return segments
}
