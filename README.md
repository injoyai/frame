### 说明
基于fiber等框架的二次封装

### 如何使用

```go
package main

import(
	"github.com/injoyai/frame"
    "github.com/injoyai/frame/fbr"
)

func main(){

    s := fbr.Default(
            fbr.WithPort(frame.DefaultPort),
            fbr.BindHtml(401, frame.Html401),
            fbr.BindHtml(403, frame.Html403),
            fbr.BindHtml(404, frame.Html404),
            fbr.WithStatic("./dist/"),
        )
	
	s.Group("/api", func(g fbr.Grouper) {
        g.ALL("/test", func(c fbr.Ctx) {
            c.Succ(667)
        })
        g.ALL("/test2", func(c fbr.Ctx) {
            c.JSON(map[string]any{"key": "value"})
        })
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
    })
    s.Run()
}
```