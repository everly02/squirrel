package squirrel

import (
	"net/http"
	"regexp"
	"strings"
)

type HandlerFunc func(ctx *Context)

type node struct {
	pattern  string  // 完整匹配的路由路径，例如 /p/:lang/doc
	part     string  // 路径中的一部分，例如 :lang
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配，part 含有 : 或 * 时为 true
}

// matchChild 匹配单个子节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 匹配所有符合条件的子节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 插入路由规则
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search 查找路由规则
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

type route struct {
	pattern *regexp.Regexp
	params  []string
	handler HandlerFunc
}

type Router struct {
	roots       map[string]*node       // 保存每种请求方式的 Trie 树根节点
	handlers    map[string]HandlerFunc // 保存每个路由对应的处理函数
	middlewares []MiddlewareFunc       // 中间件
	template    *Template
}

func NewRouter(templatePattern string) *Router {
	tmpl := NewTemplate(templatePattern)
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
		template: tmpl,
	}
}

func (r *Router) SetTemplate(template *Template) {
	r.template = template
}

func (r *Router) Render(w http.ResponseWriter, name string, data interface{}) error {
	return r.template.Render(w, name, data)
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Use 添加中间件
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

func parsePattern(path string) []string {
	vs := strings.Split(path, "/") // 将路径按 / 分割

	parts := make([]string, 0) // 保存路径中的一部分
	for _, item := range vs {  // 遍历路径的每一部分
		if item != "" {
			parts = append(parts, item)
			if item[0] == ':' || item[0] == '*' {
				parts = append(parts, item)
			}
		}
	}
	return parts

}

// AddRoute 添加路由和处理函数
func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	parts := parsePattern(path)

	key := method + "-" + path
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

// GET 注册 GET 方法的路由
func (r *Router) GET(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodGet, path, handler)
}

// POST 注册 POST 方法的路由
func (r *Router) POST(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodDelete, path, handler)
}

// ServeHTTP 实现 http.Handler 接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

// matchRoute 匹配路径参数
func matchRoute(route, path string) (map[string]string, bool) {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i := range routeParts {
		if strings.HasPrefix(routeParts[i], ":") {
			params[routeParts[i][1:]] = pathParts[i]
		} else if routeParts[i] != pathParts[i] {
			return nil, false
		}
	}
	return params, true
}

// ServeStatic 提供静态文件服务
func (r *Router) ServeStatic(prefix, root string) {
	fs := http.FileServer(http.Dir(root))
	r.GET(prefix+"/*filepath", func(ctx *Context) {
		ctx.Request.URL.Path = ctx.Param("filepath")
		fs.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}
