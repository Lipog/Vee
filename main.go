package main

import (
	"Vee/vee"
)

//首先看下标准net/http是如何处理一个请求的

func main() {
	r := vee.New()
	r.GET("/", func(c *vee.Context) {
		c.String(200, "la'llallalala")
	})
	r.GET("/hello/:name", func(c *vee.Context) {
		c.String(200, "hello %s, u're at %s \n", c.Param("name"), c.Path)
	})
	r.GET("/assets/*filepath", func(c *vee.Context) {
		c.JSON(200, vee.H{
			"filepath": c.Param("filepath"),
		})
	})
	//只有在所有的路由都注册成功以后，才能进行启动
	r.Run(":8080")
}
