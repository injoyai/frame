### 说明
基于fiber等框架的二次封装

### 如何使用

```go
package main

import(
	"github.com/injoyai/frame"
    "github.com/injoyai/frame/fiber"
)

func main(){

    s := fiber.Default(
            fiber.WithPort(frame.DefaultPort),
            fiber.BindHtml(401, frame.Html401),
            fiber.BindHtml(403, frame.Html403),
            fiber.BindHtml(404, frame.Html404),
            fiber.WithStatic("./dist/"),
        )
	
	s.Group("/api", func(g fiber.Grouper) {
        g.ALL("/test", func(c fiber.Ctx) {
            c.Succ(667)
        })
        g.ALL("/test2", func(c fiber.Ctx) {
            c.JSON(map[string]any{"key": "value"})
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
    })
    s.Run()
}
```