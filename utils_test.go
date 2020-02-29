package pipeflow

import (
	"testing"
)

func Test_SplitPathIntoSegments(t *testing.T) {
	if len(splitPathIntoSegments("/")) != 0 {
		t.Fatal("segments should be empty")
	}

	if len(splitPathIntoSegments("/hello/{var}")) != 2 {
		t.Fatal("segments length should be 2")
	}

	if len(splitPathIntoSegments("/hello/{var}/")) != 2 {
		t.Fatal("segments length should be 2")
	}
}

func Test_SplitPathIntoSegments_Panic(t *testing.T) {
	defer shouldPanic(t)
	splitPathIntoSegments("")
}
