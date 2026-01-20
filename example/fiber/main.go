package main

import (
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/fbr"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/frame/middle/swagger"
)

//go:embed dist/*
var dist embed.FS

func main() {

	s := fbr.Default(
		fbr.WithPort(frame.DefaultPort),
		fbr.WithPrintRoutes(false),
		fbr.WithSwagger(&swagger.Swagger{
			IndexPath:    "/swagger",
			JsonPath:     "/swagger/swagger.json",
			JsonFilename: "./example/fiber/docs/swagger.json",
			UI:           swagger.DefaultUI,
		}),

		//fiber.WithContext(func() context.Context {
		//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		//	return ctx
		//}()),
		//fiber.WithShutdown(func(err error) {
		//	logs.Err(err)
		//}),

		fbr.BindHtml(401, frame.Html401),
		fbr.BindHtml(403, frame.Html403),
		fbr.BindHtml(404, frame.Html404),
		fbr.BindCode(500, func(c fbr.Ctx) {
			log.Println("错误:", string(c.Response().Body()))
			c.Html(http.StatusInternalServerError, frame.Html500)
			c.Succ("系统开小差啦,请稍后再试")
		}),

		fbr.WithResponseCode("success", "fail", "unauthorized", "forbidden"),

		fbr.WithEmbed(dist),
		fbr.WithStatic("./example/dist/"),
		fbr.WithALL("/api", func(c fbr.Ctx) {
			c.Succ(123)
		}),
		fbr.WithGroup("/api2", func(g fbr.Grouper) {
			g.ALL("xxx", func(c fbr.Ctx) {
				c.Succ(c.Path())
			})
		}),
	)

	s.Embed("/dist", dist)
	s.Group("/api", func(g fbr.Grouper) {
		g.Group("/user", fbr.WithStruct(&User{1.88}))
		g.GET("/xxx", func(c fbr.Ctx) { c.Succ("xxx") })
		g.ALL("/succ", func(c fbr.Ctx) {
			c.Write([]byte("ctx message"))
			c.Succ(667)
			in.Succ(666)
		})
		g.ALL("/json", func(c fbr.Ctx) {
			c.JSON(map[string]any{"key": "value"})
		})
		g.ALL("/500", func(c fbr.Ctx) {
			in.Text(500, "500")
		})
		g.ALL("/401", func(c fbr.Ctx) {
			in.Text(401, "401")
		})
		g.ALL("/403", func(c fbr.Ctx) {
			in.Text(403, "403")
		})
		g.ALL("/404", func(c fbr.Ctx) {
			in.Text(404, "404")
		})
		g.ALL("/panic", func(c fbr.Ctx) {
			panic("panic")
		})
		g.ALL("/proxy", func(c fbr.Ctx) {
			c.Proxy("http://127.0.0.1:8080/api/json")
		})
		g.ALL("/redirect", func(c fbr.Ctx) {
			c.RedirectTo("/api/succ")
		})
		g.Redirect("/redirect2", "/api/succ")
		g.ALL("/ws", func(c fbr.Ctx) {
			c.Websocket(func(ws *fbr.Websocket) {
				for {
					<-time.After(time.Second * 3)
					_, err := ws.WriteText("hello")
					if err != nil {
						return
					}
				}
			})
		})
		g.ALL("/sse", func(c fbr.Ctx) {
			c.SSE(func(w fbr.SSE) {
				for i := 0; i < 10; i++ {
					w.WriteString("data: " + conv.String(i))
					<-time.After(time.Second)
				}
			})
		})
		g.ALL("/body", func(c fbr.Ctx) {
			_ = c.Body()
			in.Text200(string(c.Body()))
		})
		g.Static("", "./example/fiber/dist")
		g.Embed("/dist", dist)
		g.ALL("/dist", fbr.WithStatic("./example/dist/"))
		g.ALL("/:key", func(c fbr.Ctx) {
			c.Succ(c.GetString("key"))
		})
	})

	s.Run()

}

type User struct {
	high float64
}

// Name 指针下方法可以注册,只能指针对象注册
func (u *User) Name(c fbr.Ctx) {
	in.Succ("injoy")
}

// Age 指针对象可以注册非指针下方法
func (u User) Age(c fbr.Ctx) {
	in.Succ(18)
}

// High 测试引用对象字段,对象不能为nil
func (u User) High(c fbr.Ctx) {
	in.Succ(u.high)
}

// test 小写不会注册
func (u User) test(c fbr.Ctx) {
	in.Succ("test")
}
