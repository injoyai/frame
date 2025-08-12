package in

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"net/http"
	"net/http/httputil"
	"net/url"
	"unsafe"
)

func New(op ...Option) Client {
	c := &client{}
	c.SetHandlerWithCode(
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
		http.StatusForbidden,
	)
	for _, f := range op {
		f(c)
	}
	return c
}

type client struct {
	cache         *maps.Safe
	writerOptions []WriterOption
	bind          map[string]http.HandlerFunc //自定义接口绑定
	succ          func(data any, count ...int64)
	fail          func(data any)
	unauthorized  func()
	forbidden     func()
}

func (this *client) Cache() *maps.Safe {
	if this.cache == nil {
		this.cache = maps.NewSafe()
	}
	return this.cache
}

// SetHandlerWithCode 设置响应成功失败等
func (this *client) SetHandlerWithCode(succ, fail, unauthorized, forbidden any) {
	this.succ = this.NewSuccWithCode(succ)
	this.fail = this.NewFailWithCode(fail)
	this.unauthorized = this.NewUnauthorizedWithCode(unauthorized)
	this.forbidden = this.NewForbiddenWithCode(forbidden)
}

func (this *client) Bind(path string, handler http.HandlerFunc) {
	if this.bind == nil {
		this.bind = map[string]http.HandlerFunc{}
	}
	this.bind[path] = handler
}

// SetWriterOption 设置退出选项
func (this *client) SetWriterOption(f ...WriterOption) {
	this.writerOptions = append(this.writerOptions, f...)
}

//=================================Proxy=================================//

func (this *client) Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	u, err := url.Parse(uri)
	if err != nil {
		this.fail(err)
		return
	}
	p := httputil.NewSingleHostReverseProxy(u)
	p.ServeHTTP(w, r)
}

//=================================Middle=================================//

func (this *client) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				this.RecoverWriter(e, w)
			}
		}()
		if handler, ok := this.bind[r.URL.Path]; ok {
			if handler != nil {
				handler(w, r)
			}
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (this *client) RecoverFunc(f http.HandlerFunc) http.Handler {
	return this.Recover(f)
}

// RecoverWriter 例gf等web框架只需要这一半即可,但是Bind会失效
func (this *client) RecoverWriter(e any, w http.ResponseWriter) {
	switch w2 := e.(type) {
	case Writer:
		w2.WriteTo(w)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		s := conv.String(e)
		w.Write(*(*[]byte)(unsafe.Pointer(&s)))
	}
}
