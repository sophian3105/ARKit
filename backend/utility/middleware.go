package utility

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"firebase.google.com/go/v4/auth"
)

// The data passed and can be filled
type MiddlewareData struct{
	*auth.Token

	// This field allows you to abort right away and return this response
	// Default 0 means there is no status
	abortStatus int
}

func (d *MiddlewareData) AbortWithStatus(status int) {
	if d == nil {
		panic("AbortWithStatus d is nil")
	}

	d.abortStatus = status
}

func (d *MiddlewareData) ShouldAbort() bool {
	return d == nil || d.abortStatus != 0
}

// A main handler for each middleware
type MiddlewareHandler func(http.ResponseWriter, *http.Request, *MiddlewareData)

// Base data to pass into a main handler
type Middleware struct {
	handler MiddlewareHandler
}

func NewMiddleware(h MiddlewareHandler) *Middleware {
	return &Middleware{
		handler: h,
	}
}

// The starting base for a chain of middleware
type Router struct {
	middlewares []*Middleware
	path string
	mux *http.ServeMux
	lock sync.Mutex
}

func NewRouter(mux *http.ServeMux, path string) *Router {
	path = mergePath("", path)

	return &Router{
		mux: mux,
		path: path,
	}
}

func (b *Router) Use(m ...*Middleware) error {
	b.middlewares = append(b.middlewares, m...)
	return nil
}

// Creates a new branch off of the current base, reusing every middleware
// up to this point
func (b *Router) Branch(path string) *Router {
	if b == nil {
		// Ensure the pointer is valid
		panic("router is null when branching")
	}
	
	tmp := NewRouter(b.mux, mergePath(b.path, path))

	// Share the same middlewares but not array
	tmp.middlewares = append(tmp.middlewares, b.middlewares...)

	return tmp
}

// Another name for a handler. The main end request 
type RouteHandler MiddlewareHandler

func (b *Router) Handle(path string, h RouteHandler, methods ...string) {
	if b == nil {
		panic("Handle router pointer can't be null")
	}

	if len(methods) == 0 {
		methods = append(methods, http.MethodGet)
	}

	data := &MiddlewareData{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, m := range b.middlewares {
			m.handler(w, req, data)

			if data.ShouldAbort() {
				w.WriteHeader(data.abortStatus)
				return
			}
		}
		
		h(w, req, data)
	})

	for _, method  := range methods {
		p := fmt.Sprintf("%s %s", method, mergePath(b.path, path))
		log.Println(p)
		b.mux.Handle(p, handler)
	}
}

// Merge 2 paths together
// We assume that path1 is in a "correct" state
func mergePath(path1 string, path2 string) string {
	if path2 == "" || path2[0] != '/' {
		path2 = "/" + path2
	}

	if (len(path1) > 0 && path1[len(path1) - 1] == '/') {
		// Special case
		path1 = path1[:len(path1) - 1]
	}

	return path1 + path2
}

// A few default implementations
func LoggerMiddleware() *Middleware {
	return NewMiddleware(
		func(w http.ResponseWriter, req *http.Request, d *MiddlewareData) {
			log.Println("executing request")
		},
	)
}

func JsonMiddleware() *Middleware {
	return NewMiddleware(
		func(w http.ResponseWriter, req *http.Request, d *MiddlewareData) {
			w.Header().Set(ContentType, ApplicationJson)
		},
	)
}
