package squirrel

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type Template struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

// NewTemplate 创建一个新的 Template 实例，支持缓存
func NewTemplate(pattern string) *Template {
	t := &Template{
		templates: make(map[string]*template.Template),
	}
	t.loadTemplates(pattern)
	return t
}

// loadTemplates 解析并加载所有模板文件
func (t *Template) loadTemplates(pattern string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := filepath.Base(file)
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			panic(err)
		}
		t.templates[name] = tmpl
	}
}

func (t *Template) Render(w http.ResponseWriter, name string, data interface{}) error {
	t.mu.RLock()
	tmpl, ok := t.templates[name]
	t.mu.RUnlock()

	if !ok {
		http.Error(w, "Template not found", http.StatusNotFound)
		return nil
	}

	return tmpl.Execute(w, data)
}

// AddFunc 添加自定义函数到模板
func (t *Template) AddFunc(name string, function interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, tmpl := range t.templates {
		tmpl.Funcs(template.FuncMap{name: function})
	}
}
