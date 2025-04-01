package main

import (
	"github.com/injoyai/frame/mux"
)

func main() {
	s := mux.Default()
	s.Group("/api", func(g mux.Grouper) {
		g.ALL("/test", func(c mux.Ctx) {
			c.Succ(666)
		})
		g.ALL("/500", func(c mux.Ctx) {
			panic(500)
		})
	})
	s.Run()
}
