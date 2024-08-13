package squirrel

import (
	"errors"
	"log"
	"net/http"
)

// Server 结构体表示一个 HTTP 服务器
type Server struct {
	Router *Router
}

// NewServer 创建一个新的 Server 实例
func NewServer() *Server {
	return &Server{}
}

// Run 启动服务器
func (s *Server) Run(addr string) error {

	if s.Router == nil {
		return ErrNoRouter
	}

	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.Router)
}

// 定义一个错误常量
var ErrNoRouter = errors.New("no router found")
