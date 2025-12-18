package gins

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/injoyai/frame"
)

type Option func(s *Server)

type Ctx = gin.Context

func Default(op ...Option) *Server {
	gin.SetMode(gin.ReleaseMode)
	return New(
		WithDefault(),
		WithOption(op...),
	)
}

func New(op ...Option) *Server {
	e := gin.New()
	e.Use(WithRecover())
	s := &Server{
		Engine: e,
		port:   frame.DefaultPort,
		Log:    frame.NewLogger(),
	}
	for _, o := range op {
		o(s)
	}
	return s
}

type Server struct {
	*gin.Engine              //实例
	port        int          //端口
	Log         frame.Logger //日志
}

func (this *Server) SetPort(port int) {
	this.port = port
}

func (this *Server) Group(path string, f func(g *Grouper)) {
	group := this.RouterGroup.Group(path)
	f(&Grouper{RouterGroup: group})
}

func (this *Server) Run() error {
	addr := fmt.Sprintf(":%d", this.port)
	if this.Log != nil {
		this.Log.Printf("[%s] 开启HTTP服务...\n", addr)
	}
	return this.Engine.Run(addr)
}
