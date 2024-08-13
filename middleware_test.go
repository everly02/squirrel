package squirrel

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_Logger(t *testing.T) {
	router := NewRouter("../templates/*.html")

	// 日志中间件
	router.Use(Logger())

	// 注册路由
	router.GET("/", func(ctx *Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	// 创建一个 ResponseRecorder 以记录响应
	rr := httptest.NewRecorder()

	// 创建一个模拟请求
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 发送请求到路由器
	router.ServeHTTP(rr, req)

	// 验证响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	// 验证响应内容
	expected := `Hello, World!`
	if rr.Body.String() != expected {
		t.Errorf("expected body %v, got %v", expected, rr.Body.String())
	}
}
