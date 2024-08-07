package squirrel

import (
	"net/http"
	"regexp"
	"strings"
)

type HandlerFunc func(ctx *Context)

type route struct {
	pattern *regexp.Regexp
	params  []string
	handler HandlerFunc
}

type Router struct {
	routes      map[string][]route
	middlewares []MiddlewareFunc
	template    *Template
}

func NewRouter(templatePattern string) *Router {
	tmpl := NewTemplate(templatePattern)
	return &Router{
		routes:   make(map[string][]route),
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

// parsePattern parses a given path and returns a regular expression pattern and a list of parameter names.
//
// The function takes a path as input and splits it into segments using "/" as the delimiter. It then iterates over each segment and checks if it starts with ":". If it does, it treats it as a parameter and appends a regular expression pattern "([^/]+)" to the segments slice. It also appends the parameter name (without the ":") to the params slice. If the segment does not start with ":", it simply appends it to the segments slice.
//
// After iterating over all segments, the function joins the segments slice with "/" to form the pattern string. It then compiles the pattern string into a regular expression using regexp.MustCompile and returns it along with the params slice.
//
// Parameters:
// - path: the path to be parsed (string)
//
// Returns:
// - *regexp.Regexp: the compiled regular expression pattern (regexp.Regexp)
// - []string: the list of parameter names ([]string)
func parsePattern(path string) (*regexp.Regexp, []string) {
	var segments []string
	var params []string
	for _, segment := range strings.Split(path, "/") {
		if strings.HasPrefix(segment, ":") {
			param := segment[1:]
			segments = append(segments, `([^/]+)`)
			params = append(params, param)
		} else {
			segments = append(segments, segment)
		}
	}
	pattern := "^" + strings.Join(segments, "/") + "$"
	return regexp.MustCompile(pattern), params
}

// AddRoute 添加路由和处理函数
func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	pattern, paramNames := parsePattern(path)
	if _, exists := r.routes[method]; !exists {
		r.routes[method] = []route{}
	}
	r.routes[method] = append(r.routes[method], route{pattern, paramNames, handler})
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
	path := req.URL.Path
	method := req.Method

	if routes, ok := r.routes[method]; ok {
		for _, route := range routes {
			if matches := route.pattern.FindStringSubmatch(path); matches != nil {
				ctx := NewContext(w, req, r.template)
				ctx.PathParams = make(map[string]string)
				for i, match := range matches[1:] {
					ctx.PathParams[route.params[i]] = match
				}
				handler := route.handler
				for i := len(r.middlewares) - 1; i >= 0; i-- {
					handler = r.middlewares[i](handler)
				}
				handler(ctx)
				return
			}
		}
	}

	http.NotFound(w, req)
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
