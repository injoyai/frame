package fiber

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"net"
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
	Middle  = any
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
	bindErr := make(map[int]func(c Ctx, err *Error))
	_group := &group{
		Router:     nil,
		Respondent: in.DefaultClient,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			defer func() {
				dealRecover(c, recover())
			}()
			var e *fiber.Error
			if errors.As(err, &e) {
				//相应自定义的绑定
				if val, ok := bindErr[e.Code]; ok {
					cc := NewCtx(c, _group.Respondent)
					defer cc.free()
					val(cc, e)
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
		port:    frame.DefaultPort,
		App:     app,
		Grouper: _group,
		bindErr: bindErr,
	}
	ser.Use(WithRecover())
	ser.Use(use...)
	return ser
}

type Server struct {
	port int //端口
	App  *fiber.App
	Grouper
	bindErr map[int]func(c Ctx, err *Error)
}

// BindErr 重置响应数据
func (this *Server) BindErr(code int, f func(c Ctx, err *Error)) {
	this.bindErr[code] = f
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
