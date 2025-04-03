package fiber

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame"
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
	bindCode := make(map[int]func(c Ctx, body io.Reader))
	_group := &group{
		Router:     nil,
		Respondent: in.DefaultClient,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			defer func() { dealRecover(c, recover()) }()
			code := c.Response().StatusCode()
			if val, ok := bindCode[code]; ok {
				var r io.Reader
				switch e := err.(type) {
				case in.Writer:
					r = e
				case *fiber.Error:
					r = strings.NewReader(e.Message)
				default:
					r = bytes.NewReader(c.Response().Body())
				}
				cc := NewCtx(c, _group.Respondent)
				defer cc.free()
				val(cc, r)
				val = nil
			}
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
