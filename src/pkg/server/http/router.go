package http

import (
	"context"
	"net/http"
)

type Router struct {
	rtr    *http.ServeMux
	prefix string
}

func NewRouter() *Router {
	return &Router{
		rtr: http.NewServeMux(),
	}
}

func (r *Router) WithPrefix(prefix string) *Router {
	return &Router{rtr: r.rtr, prefix: r.prefix + prefix}
}

func (r *Router) handle(handlerName string, h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		h(w, req.WithContext(ctx))
	}
}

func (r *Router) HandleRoute(path string, h http.HandlerFunc) {
	r.rtr.Handle(r.prefix+path, r.handle(path, h))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}
