package mux

import (
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"time"
)

func WithLog() Middle {
	return func(c Ctx) {
		start := time.Now()
		defer func() {
			frame.Log.Printf("%-7s  耗时:%-9s %s  \n", c.Request().Method, time.Now().Sub(start), c.Request().RequestURI)
		}()
		c.Next()
	}
}

// WithCORS 设置跨域
func WithCORS() Option {
	return func(s *Server) {
		s.Client.SetWriterOption(func(w in.Writer) {
			w.SetCORS()
		})
	}
}

// WithSwagger 设置swagger
func WithSwagger(swag *middle.Swagger) Handler {
	return func(c Ctx) {
		if swag.Use(c.Writer(), c.Request()) {
			c.Exit()
		}
		c.Next()
	}
}

func WithPing(content ...any) Option {
	return func(s *Server) {
		s.GET("/ping", func(c Ctx) {
			if len(content) > 0 {
				c.Text200(content[0])
			}
			c.Succ(nil)
		})
	}
}

// WithPort 设置端口
func WithPort(port int) Option {
	return func(s *Server) { s.SetPort(port) }
}

// WithPrefix 设置全局前缀,注意使用
func WithPrefix(prefix string) Option {
	return func(s *Server) { s.Grouper.(*Group).prefix = prefix }
}
