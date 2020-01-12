package pipeflow

import (
	"net/http"
	"reflect"
)

// Flow is main service register center
type Flow struct {
	cors              func(ctx HTTPContext)
	middleware        []func(ctx HTTPContext, next func())
	postMiddleware    []func(ctx HTTPContext, next func())
	resource          map[string]interface{}
	resourceType      map[reflect.Type]interface{}
	requestDispatcher HTTPRequestDispatcher
	notFound          func(ctx HTTPContext, next func())
	once              bool
	init              bool
}

func (flow *Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	// Procedure
	// middleware → http request dispatcher → post middleware
	//                                                     ↓
	// middleware ← http request dispatcher ← post middleware
	ctx := HTTPContext{Request: res, ResponseWriter: writer, resource: flow.resource, resourceType: flow.resourceType, Props: map[string]interface{}{}}
	if !flow.init {
		// Add CORS to the pipeline
		if flow.cors != nil {
			flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
				flow.cors(ctx)
				next()
			})
		}

		if flow.requestDispatcher != nil {
			flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
				flow.requestDispatcher.Handle(ctx)
				next()
			})
		}

		if flow.notFound != nil {
			flow.postMiddleware = append(flow.postMiddleware, flow.notFound)
		}

		if flow.postMiddleware != nil {
			flow.middleware = append(flow.middleware, flow.postMiddleware...)
		}

		flow.init = true
	}

	invoke(flow, ctx, 0)
}

func invoke(f *Flow, ctx HTTPContext, i int) {
	if i == len(f.middleware) {
		return
	}
	f.middleware[i](ctx, func() {
		invoke(f, ctx, i+1)
	})
}

// NewFlow returns a new instance of pipeflow
func NewFlow() *Flow {
	flow := Flow{}
	flow.middleware = []func(ctx HTTPContext, next func()){}
	flow.resource = map[string]interface{}{}
	flow.resourceType = map[reflect.Type]interface{}{}
	flow.once = true
	flow.requestDispatcher = NewDefaultRequestDispatcher()
	flow.SetNotFound(NotFoundMiddleware)

	return &flow
}

// Use registers middleware
func (flow *Flow) Use(middleware func(ctx HTTPContext, next func())) {
	if middleware != nil {
		flow.middleware = append(flow.middleware, middleware)
	}
}

// Run runnable typed middleware will always invoke next
func (flow *Flow) Run(middleware func(ctx HTTPContext)) {
	if middleware != nil {
		flow.Use(func(ctx HTTPContext, next func()) {
			middleware(ctx)
			next()
		})
	}
}

// UsePost add middleware to invoke after HTTP request dispatcher
func (flow *Flow) UsePost(middleware func(ctx HTTPContext, next func())) {
	if middleware != nil {
		flow.postMiddleware = append(flow.postMiddleware, middleware)
	}
}

// RunPost add middleware must be invoked after HTTP request dispatcher
func (flow *Flow) RunPost(middleware func(ctx HTTPContext)) {
	if middleware != nil {
		flow.UsePost(func(ctx HTTPContext, next func()) {
			middleware(ctx)
			next()
		})
	}
}

// UseCors registers CORS middleware
func (flow *Flow) UseCors(origins []string, methods []string, headers []string, expose []string) {
	cors := Cors{AllowedOrigins: map[string]bool{}, AllowedMethods: methods, AllowedHeaders: headers, ExposedHeaders: expose}
	if nil != origins {
		for _, v := range origins {
			cors.AllowedOrigins[v] = true
		}
	}
	flow.cors = func(ctx HTTPContext) {
		cors.Handle(ctx)
	}
}

func (flow *Flow) SetHTTPDispatcher(hd HTTPRequestDispatcher) {
	if flow.once {
		flow.requestDispatcher = hd
	}
}

func (flow *Flow) SetNotFound(nf func(ctx HTTPContext)) {
	if nf == nil {
		flow.notFound = nil
	} else {
		flow.notFound = func(ctx HTTPContext, next func()) {
			nf(ctx)
			next()
		}
	}
}

// Map is used to add request handler
func (flow *Flow) Map(path string, handler func(ctx HTTPContext), methods []HTTPMethod) {
	if flow.requestDispatcher != nil {
		flow.requestDispatcher.Map(path, handler, methods)
		// Once Map has been called, the dispatcher cannot be replaced any more.
		flow.once = false
	}
}

func (flow *Flow) GET(path string, handler func(ctx HTTPContext)) {
	flow.Map(path, handler, []HTTPMethod{HTTPGet})
}

func (flow *Flow) POST(path string, handler func(ctx HTTPContext)) {
	flow.Map(path, handler, []HTTPMethod{HTTPPost})
}

// SetResource sets global singleton resource
func (flow *Flow) SetResource(key string, value interface{}) {
	flow.resource[key] = value
}

// SetResourceWithType sets global singleton resource using it's type as key
func (flow *Flow) SetResourceWithType(key reflect.Type, value interface{}) {
	flow.resourceType[key] = value
}

// SetResourceAlsoWithType calls SetResource and SetResourceWithType
func (flow *Flow) SetResourceAlsoWithType(key string, value interface{}) {
	flow.SetResource(key, value)
	flow.SetResourceWithType(reflect.TypeOf(value), value)
}

// GetResource gets global singleton resource preset
func (flow *Flow) GetResource(key string) interface{} {
	return flow.resource[key]
}

// GetResourceByType gets global singleton resource preset by type
func (flow *Flow) GetResourceByType(key reflect.Type) interface{} {
	return flow.resourceType[key]
}
