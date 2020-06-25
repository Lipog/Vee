package vee

import (
	"log"
	"net/http"
)

//把router提取出来，便于修改
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	router := &router{handlers: make(map[string]HandlerFunc)}
	return router
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	handler, ok := r.handlers[key]
	if ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 not found: %s \n", c.Path)
	}
}