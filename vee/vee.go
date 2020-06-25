package vee

import (
	"fmt"
	"net/http"
)

//核心是实现http.ListenAndServe里的，Handler接口所定义的ServeHTTP方法
//即type Handler interface{
//		ServeHTTP(w ResponseWriter, r *Request)
//}

//func ListenAndServe(address string, h Handler) error

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

//Engine 是所有请求的uri处理函数
type Engine struct {
	router map[string]HandlerFunc  //router用来存储请求对应的处理函数
}

//New 函数是暴露给外部，用来创建Engine实例的
func New() *Engine {
	engine := &Engine{router: make(map[string]HandlerFunc)}
	return engine
}

//有了实例以后，就要向engine里的router添加路由对应的方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//这是对外暴露的GET请求函数，使用的时候会在engine.router里注册对应的路由处理函数
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//这是对外暴露的POST方法，使用时在router里注册对应的处理函数和请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(port string) {
	http.ListenAndServe(port, engine)
}
//首先要实现Engine的ServeHTTP方法，还要处理对应的请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//从请求里获得请求的方法和pattern
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		//如果对应请求的处理函数存在于engine.router中，那么就用该函数处理请求
		handler(w, r)
	} else {
		fmt.Fprint(w, "404 not found: %s \n", r.URL.Path)
	}
}
