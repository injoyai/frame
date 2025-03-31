package main

import (
	"embed"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/fiber"
	"github.com/injoyai/frame/middle/in"
	"io"
	"time"
)

//go:embed dist/*
var dist embed.FS

func main() {

	s := fiber.Default()
	s.SetPort(frame.DefaultPort)

	s.BindCode(404, func(c fiber.Ctx, body io.Reader) {
		c.Succ("bind 404")
	})

	s.BindCode(500, func(c fiber.Ctx, body io.Reader) {
		c.Response().ResetBody()
		c.Succ("bind 500")
	})

	//s.BindCode(200, func(c fiber.Ctx, body io.Reader) {
	//	c.Response().ResetBody()
	//	c.Succ("bind 200")
	//})

	s.Use(func(c in.Client) {
		c.SetHandlerWithCode(
			"success",
			"fail",
			"unauthorized",
			"forbidden",
		)
	})

	s.Group("/api", func(g fiber.Grouper) {
		g.ALL("/succ", func(c fiber.Ctx) {
			c.Write([]byte("ctx message"))
			c.Succ(667)
			in.Succ(666)
		})
		g.ALL("/json", func(c fiber.Ctx) {
			c.JSON(map[string]any{"key": "value"})
		})
		g.ALL("/500", func(c fiber.Ctx) {
			in.Text(500, "500")
		})
		g.ALL("/respondent", func(r fiber.Respondent) {
			r.Succ("respondent")
		})
		g.ALL("/ws", func(c fiber.Ctx) {
			c.OnWebsocket(func(conn *fiber.Websocket) {
				for {
					<-time.After(time.Second * 3)
					_, err := conn.WriteText("hello")
					if err != nil {
						return
					}
				}
			})
		})

	})
	s.Use(fiber.WithEmbed("/dist", "dist", dist))
	s.Use(fiber.WithStatic("./example/testfiber/dist/"))

	s.Run()

}
