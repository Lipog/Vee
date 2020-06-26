package main

import (
	"Vee/vee"
	"net/http"
)

//首先看下标准net/http是如何处理一个请求的

func main() {
	r := vee.New()
	r.Use(vee.Logger())
	r.Use(vee.FileLogger)
	r.GET("/", vee.IndexHandler)
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *vee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Vee</h1>")
		})

		v1.GET("/hello", func(c *vee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *vee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *vee.Context) {
			c.JSON(http.StatusOK, vee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	//只有在所有的路由都注册成功以后，才能进行启动
	r.Run(":8080")
}
