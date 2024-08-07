package squirrel

import (
	"log"
	"net/http"
	"time"
)

// applyMiddlewares 应用中间件到处理函数
func applyMiddlewares(handler HandlerFunc, middlewares []MiddlewareFunc) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// Logger returns a MiddlewareFunc that logs the request details.
//
// It takes a HandlerFunc as a parameter and returns a new HandlerFunc.
// The returned HandlerFunc wraps the provided HandlerFunc and logs the request details before and after calling the next HandlerFunc.
// The logged details include the request method, URL path, protocol, status code, and the time taken to process the request.
func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			start := time.Now()
			next(ctx)
			log.Printf("%s %s %s %d %s",
				ctx.Request.Method,
				ctx.Request.URL.Path,
				ctx.Request.Proto,
				ctx.StatusCode,
				time.Since(start),
			)
		}
	}
}

// Recovery returns a MiddlewareFunc that recovers from panics.
func Recovery() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %v", err)
					ctx.AbortWithStatus(500)
				}
			}()
			next(ctx)
		}
	}
}

func ErrorHandling() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.String(http.StatusInternalServerError, "Internal Server Error")
				}
			}()
			next(ctx)
		}
	}
}

func CORS(origin string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			ctx.SetHeader("Access-Control-Allow-Origin", origin)
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if ctx.Request.Method == "OPTIONS" {
				ctx.Status(http.StatusOK)
				return
			}
			next(ctx)
		}
	}
}

func RateLimiter(limit int) MiddlewareFunc {
	var requests = make(map[string]int)
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			ip := ctx.Request.RemoteAddr
			if requests[ip] >= limit {
				ctx.Status(http.StatusTooManyRequests)
				return
			}
			requests[ip]++
			next(ctx)
			go func() {
				time.Sleep(time.Minute)
				requests[ip]--
			}()
		}
	}
}
