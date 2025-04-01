package mux

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle/in"
	"net/http"
)

func Default() *Server {
	return New(
		WithLog(),
		WithCORS(),
		WithPing("pong"),
	)
}

func New(op ...Option) *Server {
	s := &Server{
		port: frame.DefaultPort,
		Grouper: &Group{
			Router:     mux.NewRouter(),
			prefix:     "",
			middle:     nil,
			Respondent: in.DefaultClient,
		},
		Client: in.DefaultClient,
	}
	for _, v := range op {
		switch f := v.(type) {
		case func(s *Server):
			f(s)
		case func(c Ctx):
			s.Use(f)
		case Handler:
			s.Use(f)
		default:
			panic("unknown option " + fmt.Sprintf("%T", v))
		}
	}
	return s
}

type Server struct {
	Grouper
	in.Client
	port  int
	close func()
}

func (this *Server) Close() error {
	this.close()
	return nil
}

// SetPort 设置端口
func (this *Server) SetPort(port int) *Server {
	this.port = port
	return this
}

func (this *Server) Run() (err error) {
	defer func() { frame.Log.Println("HTTP服务结束监听:", err) }()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: this.Client.Recover(this.Grouper),
	}
	this.close = func() { s.Close() }
	frame.Log.Printf("[%s] 开启HTTP服务成功...\n", s.Addr)
	return s.ListenAndServe()
}
