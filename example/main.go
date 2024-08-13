package main

import (
	"github.com/everly02/squirrel"
	"net/http"
)

func main() {
	r := squirrel.NewRouter("./*.html")

	// 添加中间件
	r.Use(func(next squirrel.HandlerFunc) squirrel.HandlerFunc {
		return func(ctx *squirrel.Context) {
			next(ctx)
		}
	})

	// 注册路由
	r.GET("/", func(ctx *squirrel.Context) {
		ctx.RenderHTML(http.StatusOK, "index.html", map[string]interface{}{
			"Title": "Hello, World!",
			"Body":  "Welcome to my web framework!",
		})
	})

	server := squirrel.NewServer()
	server.Router = r
	server.Run(":8080")
}
