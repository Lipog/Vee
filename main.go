package main

import (
	"Vee/vee"
	"fmt"
	"net/http"
)

//首先看下标准net/http是如何处理一个请求的

func main() {
	r := vee.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "哈哈哈，成功啦！可以进行上下文的测试啦！")
	})

	//只有在所有的路由都注册成功以后，才能进行启动
	r.Run(":8080")
}
