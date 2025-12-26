package fbr

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle/in"
	"github.com/injoyai/frame/middle/swagger"

	"io"
	"io/fs"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func dealErr(c fiber.Ctx, e any) {
	switch err := e.(type) {
	case nil:

	case in.Writer:
		c.Response().ResetBody()
		for i, v := range err.Header() {
			c.Set(i, strings.Join(v, ","))
		}
		if err.StatusCode() >= 0 {
			c.Status(err.StatusCode())
		}
		io.Copy(c, err)

	case *fiber.Error:
		c.Status(err.Code).SendString(err.Message)

	default:
		c.Response().ResetBody()
		c.Status(http.StatusInternalServerError)
		s := conv.String(e)
		c.Write(*(*[]byte)(unsafe.Pointer(&s)))

	}
}

// WithRecover 配合in包使用，可以提前捕获，方便计时能中间件德操作
func WithRecover() Middle {
	return func(c fiber.Ctx) error {
		defer func() {
			if e := recover(); e != nil {
				dealErr(c, e)
			}
		}()
		return c.Next()
	}
}

// WithPprof 开启pprof
func WithPprof() Middle {
	return pprof.New(pprof.Config{})
}

// WithCORS 设置响应CORS头部
func WithCORS() Middle {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		c.Set("Access-Control-Allow-Max-Age", "3600")
		return c.Next()
	}
}

// WithLog 打印请求日志,配合WithRecover使用
func WithLog() Middle {
	logger := frame.NewLogger()
	return func(c fiber.Ctx) error {
		start := time.Now()
		defer func() {
			logger.Printf("| %-7s | %-3d   | %-12s | %s  \n", c.Method(), c.Response().StatusCode(), time.Now().Sub(start), c.OriginalURL())
		}()
		return c.Next()
	}
}

func WithLog2(color ...bool) Middle {
	return func(c fiber.Ctx) error {
		start := time.Now()
		defer func() {
			nowStr := time.Now().Format("15:04:05")
			if len(color) > 0 && color[0] {
				code := c.Response().StatusCode()
				codeStr := conv.String(code)
				switch {
				case code < 300:
					codeStr = fmt.Sprintf("\x1b[32m%-3d\x1b[0m", code)
				case code < 400:
					codeStr = fmt.Sprintf("\x1b[33m%-3d\x1b[0m", code)
				default:
					codeStr = fmt.Sprintf("\x1b[31m%-3d\x1b[0m", code)
				}
				fmt.Printf("%s | \x1B[34m%-7s\x1B[0m | %s   | %-12s | %s  \n", nowStr, c.Method(), codeStr, time.Now().Sub(start), c.OriginalURL())
				return
			}
			fmt.Printf("%s | %-7s | %-3d   | %-12s | %s  \n", nowStr, c.Method(), c.Response().StatusCode(), time.Now().Sub(start), c.OriginalURL())
		}()
		return c.Next()
	}
}

func WithLimit(limit float64) Middle {
	return func(c fiber.Ctx) error {
		//todo
		return c.Next()
	}
}

// WithResponseCode 设置响应码
func WithResponseCode(succ, fail, unauthorized, forbidden string) Middle {
	return func(c in.Client) {
		c.SetHandlerWithCode(
			succ,
			fail,
			unauthorized,
			forbidden,
		)
	}
}

/*

Handler func(c fiber.Ctx)

*/

// WithSwagger 加载swagger
func WithSwagger(swag *swagger.Swagger) Handler {
	return func(c Ctx) {
		_, err := swag.Do(
			string(c.Request().URI().Path()),
			func(r io.Reader, contentType string) {
				c.Custom(http.StatusOK, r, http.Header{fiber.HeaderContentType: []string{contentType}})
			},
		)
		c.CheckErr(err)
		c.next()
	}
}

// WithPing 状态检查
func WithPing() Handler {
	return func(c Ctx) {
		switch c.Path() {
		case "/ping":
			c.Succ("pong")
		}
		c.next()
	}
}

// WithEmbed 加载嵌入文件
func WithEmbed(e embed.FS) Handler {
	return WithFS(e)
}

// WithFS 加载文件
func WithFS(e fs.FS) Handler {
	entries, err := fs.ReadDir(e, ".")
	if err != nil {
		panic(err)
	}
	// 只有一个顶层目录且是目录，自动去掉前缀
	if len(entries) == 1 && entries[0].IsDir() {
		e, err = fs.Sub(e, entries[0].Name())
		if err != nil {
			panic(err)
		}
	}
	return func(c Ctx) {
		filename, _ := strings.CutPrefix(c.Path(), c.Route().Path)
		filename = conv.Select(filename == "/" || filename == "", "index.html", filename)
		f, err := e.Open(filename)
		if os.IsNotExist(err) {
			c.next()
			return
		}
		c.CheckErr(err)
		defer f.Close()
		c.Custom200(f, nil)
	}
}

// WithStatic 加载静态文件,本地目录
func WithStatic(root string) Handler {
	h := static.New(root)
	return func(c Ctx) {
		dealErr(c, h(c))
	}
}

// WithCache 缓存无参的GET请求
func WithCache(expiration ...time.Duration) Handler {
	type Message struct {
		Header []byte `json:"header"`
		Body   []byte `json:"body"`
	}
	cache := maps.NewSafe()
	return func(c Ctx) {
		if c.Method() == fiber.MethodGet && len(c.Queries()) == 0 {
			data, err := cache.GetOrSetByHandler(c.Path(), func() (any, error) {
				if err := c.Next(); err != nil {
					return nil, err
				}
				return &Message{
					Header: c.Response().Header.Header(),
					Body:   c.Response().Body(),
				}, nil
			}, expiration...)
			c.CheckErr(err)
			header := http.Header{}
			err = json.Unmarshal(data.(*Message).Header, &header)
			c.CheckErr(err)
			c.Custom200(data.(*Message).Body, header)
		}
		c.next()
	}
}

// BindHtml 绑定响应状态码,响应html
func BindHtml(code int, html string) Handler {
	return BindCode(code, func(c Ctx) {
		c.Html(code, html)
	})
}

// BindCode 绑定响应状态码
// 需要在WithRecover之前,才能改变状态码
// 这个Bing可以让log打印准确状态,系统自带的是最后执行的,log打印不准
func BindCode(code int, handler Handler) Handler {
	return func(c Ctx) {
		if handler != nil {
			defer func() {
				if e := recover(); e != nil {
					dealErr(c, e)
				}
				if code == c.Response().StatusCode() {
					handler(c)
				}
			}()
		}
		c.next()
	}
}

// BindCodes 绑定响应状态码
func BindCodes(m map[int]Handler) Handler {
	return func(c Ctx) {
		defer func() {
			if e := recover(); e != nil {
				dealErr(c, e)
			}
			handler, ok := m[c.Response().StatusCode()]
			if ok && handler != nil {
				handler(c)
			}
		}()
		c.next()
	}
}

/*

Option func(s *Server)

*/

// WithPort 设置监听端口
func WithPort(port int) Option {
	return func(s *Server) {
		s.SetPort(port)
	}
}

// WithListenConfig 设置监听配置
func WithListenConfig(cfg ListenConfig) Option {
	return func(s *Server) {
		s.ListenConfig = cfg
	}
}

// WithPrintRoutes 打印路由信息
func WithPrintRoutes(b ...bool) Option {
	return func(s *Server) {
		s.ListenConfig.EnablePrintRoutes = len(b) == 0 || b[0]
	}
}

// WithShutdown 设置服务关闭事件
func WithShutdown(f func(err error)) Option {
	return func(s *Server) {
		s.ListenConfig.OnShutdownSuccess = func() { f(nil) }
		s.ListenConfig.OnShutdownError = f
	}
}

// WithContext 设置服务上下文
func WithContext(ctx context.Context) Option {
	return func(s *Server) {
		s.ListenConfig.GracefulContext = ctx
	}
}

// WithStruct 注册对象的方法到路由
func WithStruct(a any) func(g Grouper) {
	t := reflect.TypeOf(a)
	if t.Kind() != reflect.Struct && t.Elem().Kind() != reflect.Struct {
		panic("type must be *struct/struct !!!")
	}
	ctxType := reflect.TypeOf((*Ctx)(nil)).Elem()
	return func(g Grouper) {
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if m.Type.NumIn() == 2 && m.Type.In(1) == ctxType {
				path := strings.ToLower(m.Name)
				g.ALL(path, func(c Ctx) {
					m.Func.Call([]reflect.Value{
						reflect.ValueOf(a),
						reflect.ValueOf(c),
					})
				})
			}
		}
	}
}

// NewWithStruct 注册对象的方法到路由
func NewWithStruct(f func(g Grouper, funcName string, f Handler)) func(a any) func(g Grouper) {
	return func(a any) func(g Grouper) {
		t := reflect.TypeOf(a)
		if t.Kind() != reflect.Struct && t.Elem().Kind() != reflect.Struct {
			panic("type must be *struct/struct !!!")
		}
		ctxType := reflect.TypeOf((*Ctx)(nil)).Elem()
		return func(g Grouper) {
			for i := 0; i < t.NumMethod(); i++ {
				m := t.Method(i)
				if m.Type.NumIn() == 2 && m.Type.In(1) == ctxType {
					f(g, m.Name, func(c Ctx) {
						m.Func.Call([]reflect.Value{
							reflect.ValueOf(a),
							reflect.ValueOf(c),
						})
					})
				}
			}
		}
	}
}
