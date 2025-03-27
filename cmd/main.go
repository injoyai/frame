package main

import (
	"embed"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/injoyai/frame/fiber"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/frame/mux"
)

//go:embed dist/*
var dist embed.FS

func main() {
	{
		s := fiber.Default()
		s.Use(fiber.WithEmbed("/", "dist", dist))
		s.Group("/api", func(g fiber.Grouper) {
			g.ALL("/test", func(c fiber.Ctx) {
				in.Succ(666)
			})
			g.ALL("/test2", func(c fiber.Ctx) {
				c.JSON(map[string]any{"key": "value"})
			})
			g.Use(static.New("./cmd/dist"))

		})
		s.Run()
	}

	{
		s := mux.New()
		s.Group("/api", func(g *mux.Grouper) {
			g.ALL("/test", func(ctx *mux.Request) {
				in.Succ(666)
			})
			g.ALL("/test2", func(ctx *mux.Request) {
				in.Json200(map[string]any{"key": "value"})
			})
		})
		s.Run()
	}
}
