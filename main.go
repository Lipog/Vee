package main

import (
	"Vee/vin"
	"net/http"
)

//首先看下标准net/http是如何处理一个请求的

func main() {
	r := vin.New()
	r.Use(vin.Logger())
	r.Use(vin.FileLogger)
	r.GET("/", vin.IndexHandler)
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *vin.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Vee</h1>")
		})

		v1.GET("/hello", func(c *vin.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *vin.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *vin.Context) {
			c.JSON(http.StatusOK, vin.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	//只有在所有的路由都注册成功以后，才能进行启动
	r.Run(":8080")
}
