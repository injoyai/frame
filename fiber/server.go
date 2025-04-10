package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame"
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
		func(r Requester)
		func(r Respondent)
	*/
	Handler = any

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
		ErrorHandler:              dealErr,
		DisableDefaultContentType: true,
	})
	_group.Router = app.Group("")
	ser := &Server{
		port:    frame.DefaultPort,
		App:     app,
		Grouper: _group,
	}
	ser.Use(WithRecover())
	ser.Use(use...)
	return ser
}

type Server struct {
	port int //端口
	App  *fiber.App
	Grouper
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
