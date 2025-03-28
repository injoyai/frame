package in

import "net/http"

type Option func(c Client)

// WithCORS 设置响应CORS头部
func WithCORS() Option {
	return func(c Client) {
		c.SetWriterOption(func(w Writer) {
			w.SetCORS()
		})
	}
}

func WithBind(path string, handler http.HandlerFunc) Option {
	return func(c Client) {
		c.Bind(path, handler)
	}
}

// WithDefault 默认
func WithDefault() Option {
	return func(c Client) {
		c.SetWriterOption(func(w Writer) {
			w.SetCORS()
		})
		c.Bind("/ping", func(w http.ResponseWriter, r *http.Request) { c.Succ(nil) })
		c.SetHandlerWithCode(
			http.StatusOK,
			http.StatusInternalServerError,
			http.StatusUnauthorized,
			http.StatusForbidden,
		)
	}
}
