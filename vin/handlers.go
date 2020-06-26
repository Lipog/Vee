package vin

import "fmt"

func IndexHandler(c *Context) {
	a := fmt.Sprintf("aaa测试测试2222")
	c.Writer.Write([]byte(a))
	Vlog.Info("测试日志的运行情况")
}
