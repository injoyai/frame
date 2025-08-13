package main

import (
	"embed"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/fiber"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"log"
	"net/http"
	"time"
)

//go:embed dist/*
var dist embed.FS

func main() {

	s := fiber.Default(
		fiber.WithPort(frame.DefaultPort),
		fiber.WithPrintRoutes(false),
		fiber.WithSwagger(&middle.Swagger{
			IndexPath:    "/swagger",
			JsonPath:     "/swagger/swagger.json",
			JsonFilename: "./example/fiber/docs/swagger.json",
			UI:           middle.DefaultSwaggerUI,
		}),

		//fiber.WithContext(func() context.Context {
		//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		//	return ctx
		//}()),
		//fiber.WithShutdown(func(err error) {
		//	logs.Err(err)
		//}),

		fiber.BindHtml(401, frame.Html401),
		fiber.BindHtml(403, frame.Html403),
		fiber.BindHtml(404, frame.Html404),
		fiber.BindCode(500, func(c fiber.Ctx) {
			log.Println("错误:", string(c.Response().Body()))
			c.Html(http.StatusInternalServerError, frame.Html500)
			c.Succ("系统开小差啦,请稍后再试")
		}),

		fiber.WithResponseCode("success", "fail", "unauthorized", "forbidden"),

		fiber.WithEmbed("dist", dist),
		fiber.WithStatic("./example/fiber/dist/"),
	)

	s.Embed("/dist", "dist", dist)
	s.Group("/api", func(g fiber.Grouper) {
		g.Group("/user", fiber.WithStruct(&User{1.88}))
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
		g.ALL("/401", func(c fiber.Ctx) {
			in.Text(401, "401")
		})
		g.ALL("/403", func(c fiber.Ctx) {
			in.Text(403, "403")
		})
		g.ALL("/404", func(c fiber.Ctx) {
			in.Text(404, "404")
		})
		g.ALL("/panic", func(c fiber.Ctx) {
			panic("panic")
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
		g.ALL("/body", func(c fiber.Ctx) {
			_ = c.Body()
			in.Text200(string(c.Body()))
		})
		g.Static("", "./example/fiber/dist")
		g.Embed("/dist", "dist", dist)
		g.ALL("/dist", fiber.WithStatic("./example/fiber/dist/"))
		g.ALL("/:key", func(c fiber.Ctx) {
			c.Succ(c.GetString("key"))
		})
	})

	s.Run()

}

type User struct {
	high float64
}

// Name 指针下方法可以注册,只能指针对象注册
func (u *User) Name(c fiber.Ctx) {
	in.Succ("injoy")
}

// Age 指针对象可以注册非指针下方法
func (u User) Age(c fiber.Ctx) {
	in.Succ(18)
}

// High 测试引用对象字段,对象不能为nil
func (u User) High(c fiber.Ctx) {
	in.Succ(u.high)
}

// test 小写不会注册
func (u User) test(c fiber.Ctx) {
	in.Succ("test")
}
