package vin

import (
	"net/http"
	"strings"
)

//把router提取出来，便于修改
type router struct {
	roots map[string]*node  //roots用来存放每种晴子u方式的Trie树根节点
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	router := &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
	return router
}

//解析请求URL，只允许有一个*，找到一个*就返回
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	//因为split以后，可能产生空的字符串，所以要进行判断
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//如果出现了* 那么要直接返回
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//首先要将pattern进行拆分，分成单独的api
	parts := parsePattern(pattern)
	//然后还是构造key，用来存储路由
	key := method + "-" + pattern
	//首先判断该method对应的方法是否有其对应的trie树，没有的话要进行创建
	_, ok := r.roots[method]
	if !ok {
		//创建一个空节点，代表0层
		r.roots[method] = &node{}
	}
	//然后将整体的pattern和拆分后的parts，放入到trie树中去
	//插入路由的同时，已经为路由建立了处理函数的映射
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

//因为路由是放在trie树中的，既然有插入路由的选项，也有寻找路由的选项
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//首先将请求的api进行拆分
	searchParts := parsePattern(path)
	//params用来保存动态路由对应的实际的api
	params := make(map[string]string)

	//在get之前，判断是否有对应的方法树存在，如果不存在，则返回nil
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	//将解析过的searchParts进入到trie树中进行查找
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		//这里是将最子节点存储的动态路由里， :name   *file这种，替换为请求的路由
		for index, part := range parts {
			//如果匹配到了：，说明该part的冒号需要去掉，并且被替换成请求的api
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			//即保证该字符段，不止有*一个元素
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	//除了前面的trie树不存在，如果没有找到对应的路由，也要返回两个nil
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		//因为trie树中存储的是，原来定义的固定的路由格式
		key := c.Method + "-" + n.pattern
		//这一步，是将中间件放在所需要的请求函数之前进行执行
		//只有中间件的处理函数执行完了，才会执行handler的请求函数
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 not found, %s \n", c.Path)
		})
	}

	//当把请求的处理函数放在中间件以后，就可以调用Next函数，开始执行c.handlers里的函数了
	c.Next()
}