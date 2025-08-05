package fiber

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
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

// WithSwagger 加载swagger
func WithSwagger(swag *middle.Swagger) Middle {
	return func(c Ctx) error {
		_, err := swag.Do(
			string(c.Request().URI().Path()),
			func(r io.Reader, contentType string) {
				c.Custom(http.StatusOK, r, http.Header{fiber.HeaderContentType: []string{contentType}})
			},
		)

		c.CheckErr(err)
		return c.Next()
	}
}

// WithPing 状态检查
func WithPing() Middle {
	return func(c Ctx) error {
		switch c.Path() {
		case "/ping":
			c.Succ("pong")
		}
		return c.Next()
	}
}

// WithLog 打印请求日志,配合WithRecover使用
func WithLog() Middle {
	return func(c fiber.Ctx) error {
		start := time.Now()
		defer func() {
			frame.Log.Printf("| %-7s | %-3d   | %-12s | %s  \n", c.Method(), c.Response().StatusCode(), time.Now().Sub(start), c.OriginalURL())
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

// WithEmbed 加载嵌入文件
func WithEmbed(apiPrefix, filePrefix string, e embed.FS) Middle {
	return WithFS(apiPrefix, filePrefix, e)
}

// WithFS 加载文件
func WithFS(apiPrefix, filePrefix string, fs fs.FS) Middle {
	return func(c fiber.Ctx) error {
		filename, ok := strings.CutPrefix(c.Path(), apiPrefix)
		if !ok {
			return c.Next()
		}
		if filename == "/" || filename == "" {
			filename = "/index.html"
		}
		f, err := fs.Open(path.Join(filePrefix, filename))
		if os.IsNotExist(err) {
			return c.Next()
		}
		if err != nil {
			return err
		}
		defer f.Close()
		bs, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = c.Status(fiber.StatusOK).Write(bs)
		return err
	}
}

// WithStatic 加载静态文件,本地目录
func WithStatic(root string) Middle {
	return static.New(root)
}

// WithCache 缓存无参的GET请求
func WithCache(expiration ...time.Duration) Middle {
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
		err := c.Next()
		dealErr(c, err)
	}
}

// BindCode 绑定响应状态码
// 需要在WithRecover之前,才能改变状态码
// 这个Bing可以让log打印准确状态,系统自带的是最后执行的,log打印不准
func BindCode(code int, handler Handler) Middle {
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
		err := c.Next()
		dealErr(c, err)
	}
}

// BindCodes 绑定响应状态码
func BindCodes(m map[int]Handler) Middle {
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
		err := c.Next()
		dealErr(c, err)
	}
}
