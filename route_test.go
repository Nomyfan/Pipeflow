package pipeflow

import (
	"errors"
	"testing"
)

func shouldPanic() {
	if r := recover(); r == nil {
		panic(errors.New("should panic"))
	}
}

func TestFlowBuilder_GET(t *testing.T) {
	defer shouldPanic()

	fb := NewBuilder()
	fb.GET("/{foo}/{bar}/hello", func(ctx HTTPContext) {
	})
	fb.GET("/{bar}/{foo}/hello", func(ctx HTTPContext) {
	})
}

func TestFlowBuilder_GET2(t *testing.T) {
	fb := NewBuilder()
	fb.GET("/{foo}/hello?id&name", func(ctx HTTPContext) {
	})
	fb.GET("/{foo}/hello?id", func(ctx HTTPContext) {
	})
}

func TestFlowBuilder_Map(t *testing.T) {
	fb := NewBuilder()
	fb.GET("/{foo}/hello?id&name", func(ctx HTTPContext) {
	})
	fb.POST("/{foo}/hello?id&name", func(ctx HTTPContext) {
	})
}
