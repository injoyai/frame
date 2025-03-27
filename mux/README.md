### 轻量的HTTP服务
适用于嵌入式系统,使用upx打包之后,体积只有2-3MB
引用的第三方包
- [HTTP路由 `github.com/gorilla/mux`](https://github.com/gorilla/mux)
- [基础类型 `github.com/injoyai/base/maps`](https://github.com/injoyai/base/maps)
- [类型转换 `github.com/injoyai/conv`](https://github.com/injoyai/conv)
- [响应工具 `github.com/injoyai/goutil/frame/in/v3`](https://github.com/injoyai/goutil/frame/in/v3)

### 如何使用

```go
package main

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
)

func main() {
	s := mux.New()
	s.GET("/ping", func(r *mux.Request) {
		in.Text200("pong")
	})
	s.Run()
}

```