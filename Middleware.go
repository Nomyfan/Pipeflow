package pipeflow

// Middleware is request processing unit of pipeflow
type Middleware interface {
	Handle(ctx HTTPContext) bool
}

// RunnableMiddleware will always return true.
type RunnableMiddleware interface {
	Handle(ctx HTTPContext)
}
