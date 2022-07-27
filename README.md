# mini-web

web框架，参考gee框架

### function

- [x] 上下文
- [x] 前缀树路由
- [x] 分组路由
- [x] 中间件
- [x] 模板恢复
- [x] 错误处理

### go version

> go 1.17+

### start

```shell
go get -u github.com/cold-bin/mini-web
```

### 使用示例

```go
// @author cold bin
// @date 2022/7/26

package main

import (
	app "github.com/cold-bin/mini-web"
)

func main() {
	engine := app.New()
	engine.GET("/hello1", func(c *app.Context) {
		c.JSON(200, "hello 1")
	})
	
	engine.Run("")
}

```
