package main

import (
	"embed"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/fiber"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/logs"
	"io"
	"time"
)

//go:embed dist/*
var dist embed.FS

func main() {

	s := fiber.Default(fiber.WithSwagger(&middle.Swagger{
		IndexPath:    "/swagger",
		JsonPath:     "/swagger/swagger.json",
		JsonFilename: "./example/testfiber/docs/swagger.json",
		UI:           middle.DefaultSwaggerUI,
	}))
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
		g.ALL("/requester", func(r fiber.Requester) {
			logs.Debug(r.GetString("key"))
		})

		g.ALL("/ws", func(c fiber.Ctx) {
			c.Websocket(func(ws *fiber.Websocket) {
				for {
					<-time.After(time.Second * 3)
					_, err := ws.WriteText("hello")
					if err != nil {
						return
					}
				}
			})
		})
		g.ALL("/sse", func(c fiber.Ctx) {
			c.SSE(func(w fiber.SSE) {
				for i := 0; i < 10; i++ {
					w.WriteString("data: " + conv.String(i))
					<-time.After(time.Second)
				}
			})
		})
		g.ALL("/:key", func(c fiber.Ctx) {
			c.Succ(c.GetString("key"))
		})
	})
	s.Use(fiber.WithEmbed("/dist", "dist", dist))
	s.Use(fiber.WithStatic("./example/testfiber/dist/"))

	s.Run()

}
