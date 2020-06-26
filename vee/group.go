package vee

type RouterGroup struct {
	prefix string  //该分组所使用的前缀
	middlewares []HandlerFunc //用来存放该分组所使用的中间件
	parent *RouterGroup //当前分组的父亲是谁
	engine *Engine //所有的分组共享一个Engine实例
}
