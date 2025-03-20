package utility

import (
	"aria/backend/database"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"firebase.google.com/go/v4/auth"
)

// Context The data passed and can be filled
type Context struct {
	*auth.Token // Can be nil
	*database.Queries

	// This field allows you to abort right away and return this response
	// Default 0 means there is no status
	abortStatus struct {
		status int
		msg    string
	}

	http.ResponseWriter // Never nil
	*http.Request       // Never nil
}

func (ctx *Context) AbortWithStatus(status int, msg string) {
	if ctx == nil {
		panic("AbortWithStatus ctx is nil")
	}

	ctx.abortStatus.status = status
	ctx.abortStatus.msg = msg
}

func (ctx *Context) ShouldAbort() bool {
	return ctx == nil || ctx.abortStatus.status != 0
}

func (ctx *Context) Json(status int, v interface{}) {
	ctx.WriteHeader(status)
	ctx.ResponseWriter.Header().Set(ContentType, ApplicationJson)
	err := json.NewEncoder(ctx.ResponseWriter).Encode(v)

	if err != nil {
		// Handle encoding error
		log.Printf("Error encoding json: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError, err.Error())
	}
}

func (ctx *Context) DecodeJson(v interface{}) error {
	if ctx == nil {
		panic("DecodeJson ctx is nil")
	}

	return json.NewDecoder(ctx.Body).Decode(v)
}

// StreamDecodeJson Decodes a stream of json objects
// This function will close the channel when it's done
func StreamDecodeJson[T any](ctx *Context, c chan<- T) error {
	if ctx == nil {
		panic("StreamDecodeJson ctx is nil")
	}

	defer func() {
		close(c)
	}()

	decoder := json.NewDecoder(ctx.Body)

	for {
		var v T
		if err := decoder.Decode(&v); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		c <- v
	}

	return nil
}

func (ctx *Context) GetEmail() *string {
	if ctx == nil || ctx.Token == nil {
		panic("GetEmail ctx is nil")
	}

	emails, ok := ctx.Token.Firebase.Identities["email"].([]string)

	if !ok || len(emails) == 0 {
		return nil
	}

	return &emails[0]
}

// MiddlewareHandler A main handler for each middleware
type MiddlewareHandler func(*Context)

// Middleware Base data to pass into a main handler
type Middleware struct {
	handler MiddlewareHandler
}

func NewMiddleware(h MiddlewareHandler) *Middleware {
	return &Middleware{
		handler: h,
	}
}

// Router The starting base for a chain of middleware
type Router struct {
	middlewares []*Middleware
	path        string
	mux         *http.ServeMux
	lock        sync.Mutex
}

func NewRouter(mux *http.ServeMux, path string) *Router {
	path = mergePath("", path)

	return &Router{
		mux:  mux,
		path: path,
	}
}

func (b *Router) Use(m ...*Middleware) {
	b.middlewares = append(b.middlewares, m...)
}

// Branch Creates a new branch off of the current base, reusing every middleware
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

// RouteHandler Another name for a handler. The main end request
type RouteHandler MiddlewareHandler

func (b *Router) Handle(path string, h RouteHandler, methods ...string) {
	if b == nil {
		panic("Handle router pointer can't be null")
	}

	if len(methods) == 0 {
		methods = append(methods, http.MethodGet)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{
			ResponseWriter: w,
			Request:        req,
		}

		// handleAbort is a helper function to check if the context should abort
		// and if true should return
		handleAbort := func(handler MiddlewareHandler) bool {
			handler(ctx)
			if ctx.ShouldAbort() {
				http.Error(w, ctx.abortStatus.msg, ctx.abortStatus.status)
				log.Printf("Aborting request with status %d", ctx.abortStatus.status)
				return true
			}

			return false
		}

		for _, m := range b.middlewares {
			if handleAbort(m.handler) {
				return
			}
		}

		handleAbort(MiddlewareHandler(h)) // In case any errors occur
	})

	for _, method := range methods {
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

	if len(path1) > 0 && path1[len(path1)-1] == '/' {
		// Special case
		path1 = path1[:len(path1)-1]
	}

	return path1 + path2
}

/* A few default implementations */

var LoggerMiddleware = NewMiddleware(
	func(_ *Context) {
		log.Println("executing request")
	},
)
