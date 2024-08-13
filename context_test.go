package squirrel

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContext_RenderHTML(t *testing.T) {
	// 创建一个 ResponseRecorder 以记录响应
	rr := httptest.NewRecorder()

	// 创建一个模拟请求
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个 Template 实例
	tmpl := NewTemplate("./example/*.html")

	// 创建一个 Context 实例
	ctx := NewContext(rr, req, tmpl)

	// 渲染模板
	ctx.RenderHTML(http.StatusOK, "index.html", map[string]interface{}{
		"Title": "Test Title",
		"Body":  "This is a test body",
	})

	// 验证响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	// 验证响应内容
	expected := `<title>Test Title</title>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %v", expected)
	}
}
