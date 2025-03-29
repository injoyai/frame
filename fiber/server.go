package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"io"
	"net"
	"strings"
)

type (
	HandlerBase = fiber.Handler
	Error       = fiber.Error

	/*
		Handler 支持以下类型
		func(c Ctx)
		func(c Ctx) error
		func(c fiber.Ctx)
		func(c fiber.Ctx) error
		func(r in.Respondent)
		func(c in.Client)
		func(e *in.Exit)
	*/
	Handler = any

	Middle = Handler
)

func Default() *Server {
	return New(
		WithLog(),
		WithPprof(),
		WithCORS(),
		WithRecover(),
		WithPing(),
		WithSwagger(middle.DefaultSwagger),
	)
}

func New(use ...Middle) *Server {
	bindCode := make(map[int]func(c Ctx, body io.Reader))
	_group := &group{
		Router:     nil,
		Respondent: in.DefaultClient,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			defer func() { dealRecover(c, recover()) }()
			switch e := err.(type) {
			case in.Writer:
				if val, ok := bindCode[e.StatusCode()]; ok {
					cc := NewCtx(c, _group.Respondent)
					defer cc.free()
					val(cc, e)
					val = nil
					return nil
				}
			case *fiber.Error:
				//相应自定义的绑定
				if val, ok := bindCode[e.Code]; ok {
					cc := NewCtx(c, _group.Respondent)
					defer cc.free()
					val(cc, strings.NewReader(e.Message))
					val = nil
					return nil
				}
				return c.Status(e.Code).SendString(e.Message)
			}
			//相应成功
			return nil
		},
		DisableDefaultContentType: true,
	})
	_group.Router = app.Group("")
	ser := &Server{
		port:     frame.DefaultPort,
		App:      app,
		Grouper:  _group,
		bindCode: bindCode,
	}
	ser.Use(WithRecover())
	ser.Use(use...)
	return ser
}

type Server struct {
	port int //端口
	App  *fiber.App
	Grouper
	bindCode map[int]func(c Ctx, body io.Reader)
}

// BindCode 重置响应数据
func (this *Server) BindCode(code int, f func(c Ctx, body io.Reader)) {
	this.bindCode[code] = f
}

func (this *Server) Close() error {
	return this.App.Shutdown()
}

func (this *Server) SetPort(port int) {
	this.port = port
}

func (this *Server) Run() error {
	addr := fmt.Sprintf(":%d", this.port)
	frame.Log.Printf("[%s] 开启HTTP服务成功...\n", addr)
	return this.App.Listen(addr, fiber.ListenConfig{DisableStartupMessage: true})
}

func (this *Server) RunListener(l net.Listener) error {
	frame.Log.Printf("[%s] 开启服务成功...\n", l.Addr().String())
	return this.App.Listener(l)
}
