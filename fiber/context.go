package fiber

import (
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame/middle/in"
	"sync"
)

type (
	Bind         = fiber.Bind
	Writer       = in.Writer
	WriterOption = in.WriterOption
	Respondent   = in.Respondent
)

var ctxPoll = sync.Pool{New: func() any {
	return &ctx{}
}}

type Ctx interface {
	fiber.Ctx
	in.Respondent

	// Parse 解析body数据到ptr,需要指针类型
	Parse(ptr any)

	// OnWebsocket websocket
	OnWebsocket(handler func(conn *Websocket))

	// free 释放内存
	free()
}

func NewCtx(c fiber.Ctx, r Respondent) Ctx {
	cc := ctxPoll.Get().(*ctx)
	cc.Ctx = c
	cc.Respondent = r
	return cc
}

type ctx struct {
	fiber.Ctx
	Respondent
}

func (this *ctx) Parse(ptr any) {
	err := conv.Unmarshal(this.Ctx.Body(), ptr)
	this.CheckErr(err)
}

// free 手动释放内存
func (this *ctx) free() {
	this.Ctx = nil
	this.Respondent = nil
	ctxPoll.Put(this)
}

func (this *ctx) OnWebsocket(handler func(conn *Websocket)) {
	err := DefaultUpgrader.Upgrade(this.Ctx.RequestCtx(), func(c *websocket.Conn) {
		defer c.Close()
		handler(&Websocket{
			Conn: c,
			ctx:  this.Ctx.Context(),
		})
	})
	this.CheckErr(err)
}

func (this *ctx) WriteHeader(statusCode int) {
	this.Ctx.Status(statusCode)
}
