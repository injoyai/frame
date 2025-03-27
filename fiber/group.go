package fiber

import "github.com/gofiber/fiber/v3"

type Grouper interface {
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
}

type group struct {
	fiber.Router
}

func (this *group) Use(use ...Middle) {
	for _, v := range use {
		switch val := v.(type) {
		case HandlerBase:
			this.Router.Use(val)
		case Handler:
			this.Router.Use(func(ctx Ctx) error {
				val(ctx)
				return nil
			})
		}
	}
}

func (this *group) Group(path string, handler func(g Grouper)) {
	g := this.Router.Group(path)
	handler(&group{g})
}

func (this *group) transfer(handler Handler) HandlerBase {
	return func(ctx Ctx) error {
		handler(ctx)
		return nil
	}
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
