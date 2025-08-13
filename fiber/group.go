package fiber

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/frame/middle/in"
	"io/fs"
)

type Grouper interface {
	SetRespondent(respondent in.Respondent)
	Use(use ...Middle)
	Group(path string, handler func(g Grouper))
	ALL(path string, handler Handler)
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)
	DELETE(path string, handler Handler)
	OPTIONS(path string, handler Handler)
	HEAD(path string, handler Handler)
	CONNECT(path string, handler Handler)
	TRACE(path string, handler Handler)
	PATCH(path string, handler Handler)
	Static(path string, dir string)
	Embed(path, prefix string, e embed.FS)
	FS(path, prefix string, fs fs.FS)
	Redirect(path, to string)
}

type group struct {
	fiber.Router
	in.Respondent
}

func (this *group) SetRespondent(r in.Respondent) {
	this.Respondent = r
}

func (this *group) Use(use ...Middle) {
	for _, v := range use {
		switch f := v.(type) {
		case func(c in.Client):
			f(this.Respondent.(in.Client))
		case func(e Writer):
			this.Respondent.(in.Client).SetWriterOption(f)
		default:
			this.Router.Use(this.transfer(v))
		}
	}
}

func (this *group) transfer(handler Middle) HandlerBase {
	return func(ctx fiber.Ctx) (err error) {
		switch f := handler.(type) {
		case fiber.Handler:
			dealErr(ctx, f(ctx))
		case func(ctx fiber.Ctx):
			f(ctx)
		case func(Ctx) error:
			cc := NewCtx(ctx, this.Respondent)
			defer cc.free()
			dealErr(ctx, f(cc))
		case func(Ctx):
			cc := NewCtx(ctx, this.Respondent)
			defer cc.free()
			f(cc)
		case func(r Respondent):
			f(this.Respondent)
		case func(r Requester):
			cc := NewCtx(ctx, this.Respondent)
			defer cc.free()
			f(cc)
		case []any:
			for _, v := range f {
				if err = this.transfer(v)(ctx); err != nil {
					return
				}
			}
		case func(s *Server):
			//pass
		default:
			panic(fmt.Sprintf("unknown handler/middle %T", handler))
		}
		return
	}
}

func (this *group) Group(path string, handler func(g Grouper)) {
	g := this.Router.Group(path)
	handler(&group{Router: g, Respondent: this.Respondent})
}

func (this *group) ALL(path string, handler Handler) {
	this.Router.All(path, this.transfer(handler))
}

func (this *group) GET(path string, handler Handler) {
	this.Router.Get(path, this.transfer(handler))
}

func (this *group) POST(path string, handler Handler) {
	this.Router.Post(path, this.transfer(handler))
}

func (this *group) PUT(path string, handler Handler) {
	this.Router.Put(path, this.transfer(handler))
}

func (this *group) DELETE(path string, handler Handler) {
	this.Router.Delete(path, this.transfer(handler))
}

func (this *group) OPTIONS(path string, handler Handler) {
	this.Router.Options(path, this.transfer(handler))
}

func (this *group) HEAD(path string, handler Handler) {
	this.Router.Head(path, this.transfer(handler))
}

func (this *group) CONNECT(path string, handler Handler) {
	this.Router.Connect(path, this.transfer(handler))
}

func (this *group) TRACE(path string, handler Handler) {
	this.Router.Trace(path, this.transfer(handler))
}

func (this *group) PATCH(path string, handler Handler) {
	this.Router.Patch(path, this.transfer(handler))
}

func (this *group) Static(path, dir string) {
	this.ALL(path, WithStatic(dir))
}

func (this *group) Embed(path, prefix string, e embed.FS) {
	this.ALL(path, WithEmbed(prefix, e))
}

func (this *group) FS(path, prefix string, fs fs.FS) {
	this.ALL(path, WithFS(prefix, fs))
}

func (this *group) Redirect(path, to string) {
	this.ALL(path, func(c Ctx) { c.RedirectTo(to) })
}
