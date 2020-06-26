package vin

import (
	"Vee/vlog"
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

var Vlog *vlog.FileLogger
func FileLogger(c *Context) {
		logger := vlog.NewFileLogger("DEBUG", "./", "logtest.log", 10 * 10 *1024)
		Vlog = logger
		c.Next()
}
