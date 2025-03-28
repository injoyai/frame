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
	app := fiber.New(fiber.Config{
		//这个需要声明下,会在最后执行,否则会返回panic的数据
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(e.Code).SendString(e.Message)
			}
			return nil
		},
		DisableDefaultContentType: true,
	})
	ser := &Server{
		port: frame.DefaultPort,
		App:  app,
		Grouper: &group{
			Router:     app.Group(""),
			Respondent: in.DefaultClient,
		},
	}
	ser.Use(withRecover)
	ser.Use(use...)
	return ser
}

type Server struct {
	port int //端口
	App  *fiber.App
	Grouper
	bindErr map[int]func(ctx fiber.Ctx, err *fiber.Error)
}

// BindErr 重置响应数据
func (this *Server) BindErr(code int, f func(ctx fiber.Ctx, err *fiber.Error)) {
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
