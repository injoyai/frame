package main

import (
	"embed"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/fiber"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/logs"
	"time"
)

//go:embed dist/*
var dist embed.FS

func main() {

	s := fiber.Default(fiber.WithSwagger(&middle.Swagger{
		IndexPath:    "/swagger",
		JsonPath:     "/swagger/swagger.json",
		JsonFilename: "./example/fiber/docs/swagger.json",
		UI:           middle.DefaultSwaggerUI,
	}))
	s.SetPort(frame.DefaultPort)

	s.Use(
		fiber.BindCode(500, func(c fiber.Ctx) {
			logs.Err(string(c.Response().Body()))
			c.Succ("系统开小差啦,请稍后再试")
		}),
		fiber.BindCode(404, func(c fiber.Ctx) {
			c.Text(404, "bind 404")
		}),
	)

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
		g.ALL("/proxy", func(c fiber.Ctx) {
			c.Proxy("http://127.0.0.1:8080/api/json")
		})
		g.ALL("/redirect", func(c fiber.Ctx) {
			c.RedirectTo("/api/succ")
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
		g.Static("", "./example/fiber/dist")
		g.Embed("/dist", "dist", dist)
		g.ALL("/dist", fiber.WithStatic("./example/fiber/dist/"))
		g.ALL("/:key", func(c fiber.Ctx) {
			c.Succ(c.GetString("key"))
		})
	})

	s.Embed("/dist", "dist", dist)
	s.Use(fiber.WithEmbed("dist", dist))
	s.Use(fiber.WithStatic("./example/fiber/dist/"))

	s.Run()

}
