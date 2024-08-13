package squirrel

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {
	router := NewRouter("./example/*.html")

	// 注册路由
	router.GET("/", func(ctx *Context) {
		ctx.RenderHTML(http.StatusOK, "index.html", map[string]interface{}{
			"Title": "Home Page",
			"Body":  "Welcome to the home page",
		})
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
	expected := `<title>Home Page</title>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %v", expected)
	}
}

func TestRouter_PathParams(t *testing.T) {
	router := NewRouter(".example/*.html")

	// 注册路由
	router.GET("/users/:id", func(ctx *Context) {
		id := ctx.Param("id")
		ctx.String(http.StatusOK, "User ID: %s", id)
	})

	// 创建一个 ResponseRecorder 以记录响应
	rr := httptest.NewRecorder()

	// 创建一个模拟请求
	req, err := http.NewRequest("GET", "/users/123", nil)
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
	expected := `User ID: 123`
	if rr.Body.String() != expected {
		t.Errorf("expected body %v, got %v", expected, rr.Body.String())
	}
}
