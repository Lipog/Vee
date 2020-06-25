package vee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//定义H 使代码显得简洁
type H map[string]interface{}

type Context struct {
	//原来的w和r对象
	Writer http.ResponseWriter
	Req *http.Request
	//请求信息
	Method string
	Path string
	Params map[string]string  //存储的是对应的动态路由解析到的实例对象
	//返回信息
	StatusCode int
}

//newContext 是构造一个包含了w和r的上下文实例
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		Writer:     w,
		Req:        r,
		Method:     r.Method,
		Path:       r.URL.Path,
	}
	return c
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//上下文对象要能够解析post请求的数据，并将状态码写到w
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key , value string) {
	c.Writer.Header().Set(key, value)
}

//func (c *Context) SetData(code int, data []byte) {
//	c.SetStatus(code)
//	c.Writer.Write(data)
//}

//利用Context，向w里写入各种信息，写入的格式是string类型
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)

	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.Writer.Write([]byte(html))
}
