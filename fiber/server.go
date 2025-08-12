package fiber

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/logs"
	"net"
)

type (
	HandlerBase  = fiber.Handler
	Error        = fiber.Error
	ListenConfig = fiber.ListenConfig

	Handler = func(c Ctx)

	/*
		Middle 支持以下类型
		func(c Ctx)
		func(c Ctx) error
		func(c fiber.Ctx)
		func(c fiber.Ctx) error
		func(r Requester)
		func(r Respondent)
		func(c in.Client)
		func(e in.Writer)
		func(s *Server)
	*/
	Middle = any
)

func Default(use ...Middle) *Server {
	op := []Middle{
		WithLog(),
		WithPprof(),
		WithCORS(),
		WithRecover(),
		WithPing(),
	}
	op = append(op, use...)
	return New(op...)
}

func New(use ...Middle) *Server {
	_group := &group{
		Router:     nil,
		Respondent: in.DefaultClient,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			dealErr(c, err)
			return nil
		},
		DisableDefaultContentType: true,
	})
	_group.Router = app.Group("")
	ser := &Server{
		Grouper: _group,
		port:    frame.DefaultPort,
		App:     app,
		Log:     logs.NewEntity("").SetSelfLevel(logs.LevelInfo).SetColor(color.FgCyan).SetFormatter(logs.TimeFormatter),
		ListenConfig: &ListenConfig{
			DisableStartupMessage: true,
		},
	}
	ser.Use(WithRecover())
	ser.Use(use...)
	return ser
}

type Server struct {
	Grouper                    //路由
	Log          frame.Logger  //日志
	port         int           //端口
	App          *fiber.App    //实例
	ListenConfig *ListenConfig //监听配置
}

func (this *Server) Use(use ...Middle) {
	for _, v := range use {
		switch f := v.(type) {
		case func(s *Server):
			f(this)
		default:
			this.Grouper.Use(v)
		}
	}
}

func (this *Server) Close() error {
	return this.App.Shutdown()
}

func (this *Server) SetPort(port int) {
	this.port = port
}

func (this *Server) Run() error {
	addr := fmt.Sprintf(":%d", this.port)
	if this.Log != nil {
		this.Log.Printf("[%s] 开启HTTP服务...\n", addr)
	}
	return this.App.Listen(addr, this.listenConfig()...)
}

func (this *Server) RunListener(l net.Listener) error {
	if this.Log != nil {
		this.Log.Printf("[%s] 开启服务成功...\n", l.Addr().String())
	}
	return this.App.Listener(l, this.listenConfig()...)
}

/*






 */

func (this *Server) listenConfig() []ListenConfig {
	if this.ListenConfig == nil {
		return nil
	}
	return []ListenConfig{*this.ListenConfig}
}
