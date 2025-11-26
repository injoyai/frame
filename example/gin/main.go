package main

import (
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/gins"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/frame/middle/swagger"
)

func main() {
	s := gins.Default(
		gins.BindHtml(404, frame.Html404),
		gins.BindHtml(500, frame.Html500),
	)
	s.Use(gins.WithSwagger(swagger.Default))
	s.Group("/api", func(g *gins.Grouper) {
		g.ALL("/test", func(c *gins.Ctx) {
			in.Succ("ok")
		})
		g.ALL("/panic", func(c *gins.Ctx) {
			panic("panic")
		})
	})
	s.Static("/dist", "./example/dist/")
	s.Run()
}
