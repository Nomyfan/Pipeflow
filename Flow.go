package pipeflow

import (
	"net/http"
	"reflect"
)

// Flow is main service register center
type Flow struct {
	middleware   []func(ctx HTTPContext, next func())
	resource     map[string]interface{}
	resourceType map[reflect.Type]interface{}
}

func (flow *Flow) ServeHTTP(writer http.ResponseWriter, res *http.Request) {
	// Procedure
	// middleware → http request dispatcher → post middleware
	//                                                     ↓
	// middleware ← http request dispatcher ← post middleware
	ctx := HTTPContext{Request: res, ResponseWriter: writer, resource: flow.resource, resourceType: flow.resourceType, Props: map[string]interface{}{}}
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

type FlowBuilder struct {
	flow              *Flow
	cors              func(ctx HTTPContext)
	notFound          func(ctx HTTPContext)
	requestDispatcher HTTPRequestDispatcher
	postMiddleware    []func(ctx HTTPContext, next func())
	once              bool
}

// Build returns a flow instance
func (fb *FlowBuilder) Build() *Flow {
	flow := fb.flow
	// Add CORS to the pipeline
	if fb.cors != nil {
		cors := fb.cors
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			cors(ctx)
			next()
		})
	}

	if fb.requestDispatcher != nil {
		rd := fb.requestDispatcher
		flow.middleware = append(flow.middleware, func(ctx HTTPContext, next func()) {
			rd.Handle(ctx)
			next()
		})
		// Only when there is a request handler, not found middleware has meaning.
		if fb.notFound != nil {
			nf := fb.notFound
			fb.postMiddleware = append(fb.postMiddleware, func(ctx HTTPContext, next func()) {
				nf(ctx)
				next()
			})
		}
	}

	if fb.postMiddleware != nil {
		flow.middleware = append(flow.middleware, fb.postMiddleware...)
	}

	fb.cors = nil
	fb.requestDispatcher = nil
	fb.notFound = nil
	fb.postMiddleware = nil
	fb.flow = nil
	return flow
}

func NewBuilder() *FlowBuilder {
	flow := &Flow{}
	flow.middleware = []func(ctx HTTPContext, next func()){}
	flow.resource = map[string]interface{}{}
	flow.resourceType = map[reflect.Type]interface{}{}

	return &FlowBuilder{flow: flow, requestDispatcher: newDefaultHTTPRequestDispatcher(), notFound: NotFoundMiddleware, once: true}
}

// Use registers middleware
func (fb *FlowBuilder) Use(middleware func(ctx HTTPContext, next func())) {
	if middleware != nil {
		fb.flow.middleware = append(fb.flow.middleware, middleware)
	}
}

// Run runnable typed middleware will always invoke next
func (fb *FlowBuilder) Run(middleware func(ctx HTTPContext)) {
	if middleware != nil {
		fb.Use(func(ctx HTTPContext, next func()) {
			middleware(ctx)
			next()
		})
	}
}

// UsePost add middleware to invoke after HTTP request dispatcher
func (fb *FlowBuilder) UsePost(middleware func(ctx HTTPContext, next func())) {
	if middleware != nil {
		fb.postMiddleware = append(fb.postMiddleware, middleware)
	}
}

// RunPost add middleware must be invoked after HTTP request dispatcher
func (fb *FlowBuilder) RunPost(middleware func(ctx HTTPContext)) {
	if middleware != nil {
		fb.UsePost(func(ctx HTTPContext, next func()) {
			middleware(ctx)
			next()
		})
	}
}

// UseCors registers CORS middleware
func (fb *FlowBuilder) UseCors(origins []string, methods []string, headers []string, expose []string) {
	cors := Cors{AllowedOrigins: map[string]bool{}, AllowedMethods: methods, AllowedHeaders: headers, ExposedHeaders: expose}
	if nil != origins {
		for _, v := range origins {
			cors.AllowedOrigins[v] = true
		}
	}
	fb.cors = func(ctx HTTPContext) {
		cors.Handle(ctx)
	}
}

// SetHTTPDispatcher replaces default HTTP request handler. It can be nil.
func (fb *FlowBuilder) SetHTTPDispatcher(hd HTTPRequestDispatcher) {
	if fb.once {
		fb.requestDispatcher = hd
	}
}

// SetNotFound replaces the default not found middleware
func (fb *FlowBuilder) SetNotFound(nf func(ctx HTTPContext)) {
	fb.notFound = nf
}

// Map is used to add request handler
func (fb *FlowBuilder) Map(path string, handler func(ctx HTTPContext), methods ...HTTPMethod) {
	if fb.requestDispatcher != nil {
		fb.requestDispatcher.Map(path, handler, methods...)
		// Once Map has been called, the dispatcher cannot be replaced any more.
		fb.once = false
	}
}

func (fb *FlowBuilder) GET(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPGet)
}

func (fb *FlowBuilder) POST(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPPost)
}

func (fb *FlowBuilder) HEAD(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPHead)
}

func (fb *FlowBuilder) PUT(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPPut)
}

func (fb *FlowBuilder) DELETE(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPDelete)
}

func (fb *FlowBuilder) OPTIONS(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPOptions)
}

func (fb *FlowBuilder) TRACE(path string, handler func(ctx HTTPContext)) {
	fb.Map(path, handler, HTTPTrace)
}

// SetResource sets global singleton resource
func (fb *FlowBuilder) SetResource(key string, value interface{}) {
	fb.flow.resource[key] = value
}

// SetResourceWithType sets global singleton resource using it's type as key
func (fb *FlowBuilder) SetResourceWithType(key reflect.Type, value interface{}) {
	fb.flow.resourceType[key] = value
}

// SetResourceAlsoWithType calls SetResource and SetResourceWithType
func (fb *FlowBuilder) SetResourceAlsoWithType(key string, value interface{}) {
	fb.SetResource(key, value)
	fb.SetResourceWithType(reflect.TypeOf(value), value)
}

// GetResource gets global singleton resource preset
func (flow *Flow) GetResource(key string) interface{} {
	return flow.resource[key]
}

// GetResourceByType gets global singleton resource preset by type
func (flow *Flow) GetResourceByType(key reflect.Type) interface{} {
	return flow.resourceType[key]
}
