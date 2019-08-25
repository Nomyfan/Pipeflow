package pipeflow

type Middleware interface {
	Handle(ctx HttpContext) bool
}

// RunnableMiddleware will always return true.
type RunnableMiddleware interface {
	Handle(ctx HttpContext)
}