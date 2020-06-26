package vee

import (
	"log"
	"net/http"
)

//核心是实现http.ListenAndServe里的，Handler接口所定义的ServeHTTP方法
//即type Handler interface{
//		ServeHTTP(w ResponseWriter, r *Request)
//}

//func ListenAndServe(address string, h Handler) error

type HandlerFunc func(c *Context)

//Engine 是所有请求的uri处理函数
type Engine struct {
	*RouterGroup  //保证engine实例，拥有RouterGroup所有的能力
	router *router  //router用来存储请求对应的处理函数
	groups []*RouterGroup //engie存储所有的分组
}

//New 函数是暴露给外部，用来创建Engine实例的
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//有了分组以后，就可以把与路由相关的函数，都交给RouterGroup来实现了
//newGroupy用来创建一个新的路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//有了实例以后，就要向engine里的router添加路由对应的方法
//有了分组以后，就是以分组进行添加路由了
func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	log.Printf("Route %4s - %s \n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//这是对外暴露的GET请求函数，使用的时候会在engine.router里注册对应的路由处理函数
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//这是对外暴露的POST方法，使用时在router里注册对应的处理函数和请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(port string) {
	//各种请求，要经过engine进行处理,即通过engine的ServeHTP方法来处理
	http.ListenAndServe(port, engine)
}
//首先要实现Engine的ServeHTTP方法，还要处理对应的请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//从请求里获得w和r，并且创建一个Context实例，赋值给Context
	c := newContext(w, r)
	//然后再处理请求，相当于一开始的engine.New是创建好各种路由guize
	//然后到ServeHTTP这里，根据来的请求，去已经定义号的路由里处理请求
	engine.router.handle(c)
}
