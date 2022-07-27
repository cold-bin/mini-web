// @author cold bin
// @date 2022/7/27

package mini_web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/juju/ratelimit"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

func Cors() HandlerFunc {
	return func(c *Context) {
		method := c.Req.Method
		origin := c.Req.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Req.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.SetHeader("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有的域,也可以指定某几个特定的域
			c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//header的类型
			c.SetHeader("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//允许跨域设置                                                                                                      可以返回其他子段
			c.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.SetHeader("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.SetHeader("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                                 // 设置返回格式是json
		}

		////放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//	c.JSON(http.StatusOK, "Options Request!")
		//}
		if method == "OPTIONS" {
			c.Fail(http.StatusNoContent, http.StatusText(http.StatusNoContent))
			return
		}

		// 处理请求
		c.Next()
	}
}

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.Fail(http.StatusOK, "rate limit...")
			return
		}
		c.Next()
	}
}
