package squirrel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	PathParams     map[string]string // 存储路径参数
	StatusCode     int
	Template       *Template
}

// NewContext 创建一个新的 Context 对象
func NewContext(w http.ResponseWriter, r *http.Request, tmpl *Template) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		PathParams:     make(map[string]string),
		Template:       tmpl,
	}
}

// Query 读取查询参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// PostForm 读取表单参数
func (c *Context) PostForm(key string) string {
	return c.Request.PostFormValue(key)
}

// Param 读取路径参数
func (c *Context) Param(key string) string {
	return c.PathParams[key]
}

// BindJSON 绑定 JSON 数据到结构体
func (c *Context) BindJSON(obj interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, obj)
}

// FormFile 读取上传的文件
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.ResponseWriter.WriteHeader(code)
}

// JSON 发送 JSON 响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	json.NewEncoder(c.ResponseWriter).Encode(obj)
}

// HTML 发送 HTML 响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.ResponseWriter.Write([]byte(html))
}

// String 发送纯文本响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.ResponseWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

// Redirect 重定向
func (c *Context) Redirect(code int, location string) {
	http.Redirect(c.ResponseWriter, c.Request, location, code)
}

// AbortWithStatus 中止处理并设置状态码
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.ResponseWriter.Write([]byte(http.StatusText(code)))
}

func (c *Context) RenderHTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.Template.Render(c.ResponseWriter, name, data); err != nil {
		http.Error(c.ResponseWriter, "Template rendering error", http.StatusInternalServerError)
	}
}
