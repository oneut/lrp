package livereloadproxy

import (
	"net/http"
)

func NewRouter() *Router {
	return &Router{}
}

type route struct {
	pattern string
	handler http.Handler
}

type Router struct {
	routes []*route
}

func (h *Router) Handle(pattern string, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern == "*" {
			route.handler.ServeHTTP(w, r)
			return
		}
		if route.pattern == r.URL.Path {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
