package fiber

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
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

func WithRecover() HandlerBase {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c Ctx, e any) {
			switch w := e.(type) {
			case *in.Exit:
				for i, v := range w.Headers {
					c.Response().Header.Set(i, strings.Join(v, ","))
				}
				if w.Code >= 0 {
					c.Status(w.Code)
				}
				c.Status(w.Code)
				io.Copy(c, w.Body)
			default:
				c.Status(http.StatusInternalServerError)
				s := conv.String(e)
				c.Write(*(*[]byte)(unsafe.Pointer(&s)))
			}
		},
	})
}

func WithPprof() HandlerBase {
	return pprof.New(pprof.Config{})
}

func WithCORS() HandlerBase {
	return func(c Ctx) error {
		c.Response().Header.Set("Access-Control-Allow-Origin", "*")
		c.Response().Header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		c.Response().Header.Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		c.Response().Header.Set("Access-Control-Allow-Max-Age", "3600")
		return c.Next()
	}
}

func WithSwagger(s *middle.Swagger) HandlerBase {
	return func(c Ctx) error {
		switch string(c.Request().URI().Path()) {
		case s.IndexPath:
			in.Text200(fmt.Sprintf(middle.DefaultSwaggerUI, s.JsonPath))
		case s.JsonPath:
			in.FileLocal("json", s.Filename)
		}
		return c.Next()
	}
}

func WithPing() HandlerBase {
	return func(c Ctx) error {
		switch c.Path() {
		case "/ping":
			in.Succ("pong")
		}
		return c.Next()
	}
}

func WithLog() HandlerBase {
	return func(c Ctx) error {
		start := time.Now()
		defer func() {
			frame.Log.Printf("%-7s 耗时:%-9s %s  \n", c.Method(), time.Now().Sub(start), c.Path())
		}()
		return c.Next()
	}
}

func WithEmbed(apiPrefix, filePrefix string, e embed.FS) HandlerBase {
	return func(c Ctx) error {
		filename, _ := strings.CutPrefix(c.Path(), apiPrefix)
		f, err := e.Open(path.Join(filePrefix, filename))
		if os.IsNotExist(err) {
			return c.Next()
		}
		in.CheckErr(err)
		defer f.Close()
		in.Html200(f)
		return nil
	}
}

func WithFS(apiPrefix, filePrefix string, fs fs.FS) HandlerBase {
	return func(c Ctx) error {
		filename, _ := strings.CutPrefix(c.Path(), apiPrefix)
		f, err := fs.Open(path.Join(filePrefix, filename))
		if os.IsNotExist(err) {
			return c.Next()
		}
		in.CheckErr(err)
		defer f.Close()
		in.Html200(f)
		return nil
	}
}

func WithCache() HandlerBase {
	type Message struct {
		Header []byte `json:"header"`
		Body   []byte `json:"body"`
	}
	cache := maps.NewSafe()
	return func(c Ctx) error {
		if c.Method() == fiber.MethodGet && len(c.Queries()) == 0 {
			data, err := cache.GetOrSetByHandler(c.Path(), func() (any, error) {
				if err := c.Next(); err != nil {
					return nil, err
				}
				return &Message{
					Header: c.Response().Header.Header(),
					Body:   c.Response().Body(),
				}, nil
			})
			in.CheckErr(err)
			header := http.Header{}
			json.Unmarshal(data.(*Message).Header, &header)
			in.Custom200(data.(*Message).Body, header)
		}
		return c.Next()
	}
}
