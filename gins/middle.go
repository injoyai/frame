package gins

import (
	"github.com/gin-gonic/gin"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"io"
	"net/http"
	"time"
)

type Handler = gin.HandlerFunc

func WithDefault() func(s *Server) {
	return func(s *Server) {
		s.Use(
			WithLog(),
			WithCORS(),
			WithRecover(),
			WithPing(),
		)
	}
}

func WithOption(op ...Option) func(s *Server) {
	return func(s *Server) {
		for _, o := range op {
			o(s)
		}
	}
}

func WithRecover() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				//可能使用in包退出
				in.RecoverWriter(e, c.Writer)
			}
		}()
		c.Next()
	}
}

func WithPort(port int) func(s *Server) {
	return func(s *Server) {
		s.port = port
	}
}

// WithPing 状态检查
func WithPing() Handler {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case "/ping":
			c.Abort()
			in.Succ("pong")
		}
		c.Next()
	}
}

// WithLog 打印请求日志,配合WithRecover使用
func WithLog() Handler {
	logger := frame.NewLogger()
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			logger.Printf("| %-7s | %-3d   | %-12s | %s  \n", c.Request.Method, c.Writer.Status(), time.Now().Sub(start), c.Request.URL.String())
		}()
		c.Next()
	}
}

// WithCORS 设置响应CORS头部
func WithCORS() Handler {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Max-Age", "3600")
		c.Next()
	}
}

func BindHtml(code int, html string) func(s *Server) {
	return BindCode(code, func(c *gin.Context) {
		in.Html(code, html)
	})
}

func BindCode(code int, handler Handler) func(s *Server) {
	return func(s *Server) {
		switch code {
		case http.StatusNotFound:
			s.NoRoute(handler)
		case http.StatusMethodNotAllowed:
			s.Use(handler)
		case http.StatusInternalServerError:
			s.Use(func(c *gin.Context) {
				defer func() {
					if e := recover(); e != nil {
						switch x := e.(type) {
						case in.Writer:
							if x.StatusCode() == http.StatusInternalServerError {
								handler(c)
							}
							x.Exit()
						default:
							handler(c)
						}
					}
				}()
				c.Next()
			})
		default:
			s.Use(func(c *gin.Context) {
				c.Next()
				if code == c.Writer.Status() {
					handler(c)
				}
			})
		}
	}
}

func WithSwagger(swag *middle.Swagger) Handler {
	return func(c *gin.Context) {
		ok, err := swag.Do(
			c.Request.URL.Path,
			func(r io.Reader, contentType string) {
				c.Writer.WriteHeader(http.StatusOK)
				c.Writer.Header().Set("Content-Type", contentType)
				io.Copy(c.Writer, r)
			},
		)
		if err != nil {
			//中途退出需要使用这个函数,否则状态码会变成404,然后触发NoRoute事件
			c.Abort()
			in.Err(err)
		}
		if ok {
			c.Abort()
			return
		}
		c.Next()
	}
}
