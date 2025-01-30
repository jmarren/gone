package gone

import (
	// "errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"
)

type middleware func(http.Handler) http.Handler

type Route[T any] struct {

	// middlewares to apply
	middlewares []middleware

	// pattern to match
	pattern string

	// mux
	mux *http.ServeMux

	// get handler
	get http.Handler

	// post handler
	post http.Handler

	// put handler
	put http.Handler

	// put handler
	delete http.Handler

	// parent
	parent *Route[T]

	// whether the route is registered
	registered bool

	// subroutes
	subRoutes map[string]*Route[T]

	datastore T
}

// create an app
func New[T any]() *Route[T] {

	// the default pattern is "/"
	return newRoute[T]("/")

}

// Creates a new Route
func newRoute[T any](pattern string) *Route[T] {

	// create a new instance of Route
	r := new(Route[T])

	// assign it a serveMux
	r.mux = http.NewServeMux()

	// assign its pattern
	r.pattern = pattern

	// assign it an empty subRoutes map
	r.subRoutes = make(map[string]*Route[T])

	// return
	return r
}

func (r *Route[T]) SetData(data T) {
	r.datastore = data
}

func (r *Route[T]) GetData() T {
	return r.datastore
}

// Adds a subroute with the provided pattern and handler
func (r *Route[T]) Then(pattern string) *Route[T] {

	// create new route
	rt := newRoute[T](pattern)

	// add parent
	rt.parent = r

	// add parent middlewares
	rt.middlewares = r.middlewares

	// append parent pattern
	rt.pattern = fmt.Sprintf("%s%s", r.pattern, pattern)

	// add route as subroute of parent
	r.subRoutes[pattern] = rt

	// return
	return rt
}

// Appends middlewares to a route
func (r *Route[T]) Use(middlewares ...middleware) {

	// for each middleware supplied, append the route's middleware
	for _, m := range middlewares {
		r.middlewares = append(r.middlewares, m)
	}
}

// HANDLERS ///////////////////////////////////

// adds a get handler to the route
func (r *Route[T]) Get(handler http.Handler) {
	r.get = handler
}

// adds a post handler to the route
func (r *Route[T]) Post(handler http.Handler) {
	r.post = handler
}

// adds a put handler to the route
func (r *Route[T]) Put(handler http.Handler) {
	r.put = handler
}

// adds a delete handler to the route
func (r *Route[T]) Delete(handler http.Handler) {
	r.delete = handler
}

// Registers a route and all of its subRoutes recursively
func (r *Route[T]) Register(mux *http.ServeMux) {

	// apply the routes handlers to its servemux
	r.applyRoutes()

	// Handle the route with its pattern using
	// the provided servemux
	mux.Handle(r.pattern, r.mux)

	// for each subRoute, register
	for _, sbrt := range r.subRoutes {
		sbrt.Register(mux)
	}
}

// Serve //////////////////////////////////////
func (r *Route[T]) Serve(port string) {

	// create a new serveMux to handle all routes
	// (the orignal serveMux created is responsible
	// for handling its [GET /], [POST /], etc. routes,
	// so we need a separate mux).
	rootMux := http.NewServeMux()

	// register the route using the newly created
	// serveMux
	r.Register(rootMux)

	// create a server
	s := &http.Server{
		Addr:           port,
		Handler:        rootMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("starting GONE server")

	// Listen and Serve
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}

func (r *Route[T]) applyRoutes() {
	// apply middlewares
	r.applyMiddlewares()

	// handle routes
	r.mux.Handle(fmt.Sprintf("GET %s", r.pattern), r.get)
	r.mux.Handle(fmt.Sprintf("POST %s", r.pattern), r.post)
	r.mux.Handle(fmt.Sprintf("PUT %s", r.pattern), r.put)
	r.mux.Handle(fmt.Sprintf("DELETE %s", r.pattern), r.delete)

}

// apply a routes middleware by reassigning each handler
// to the result of passing it (the handler), along with
// the route objects middlewares into Chain()
func (r *Route[T]) applyMiddlewares() {
	// apply middleware
	r.get = Chain(r.get, r.middlewares...)
	r.post = Chain(r.post, r.middlewares...)
	r.put = Chain(r.put, r.middlewares...)
	r.delete = Chain(r.delete, r.middlewares...)
}

// Reverse the provided middleware slice, then loop through and
// resassign the handler to the ouput after passing it to the
// middleware
func Chain(h http.Handler, middlewares ...middleware) http.Handler {
	if len(middlewares) < 1 {
		return h
	}

	slices.Reverse(middlewares)
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
