package core

type Middleware interface {
	Handle(ctx HttpContext) bool
}

// This type of middleware will always return true.
type RunnableMiddleware interface {
	Handle(ctx HttpContext)
}