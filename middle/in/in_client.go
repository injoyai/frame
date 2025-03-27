package in

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

type Option func(c *Client)

// WithCORS 设置响应CORS头部
func WithCORS() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
	}
}

// WithJson 设置响应json头部
func WithJson() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderJson()
		})
	}
}

// WithDefault 默认
func WithDefault() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
			e.SetHeaderJson()
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

func New(op ...Option) *Client {
	c := &Client{
		Safe:    maps.NewSafe(),
		BindMap: map[string]http.HandlerFunc{},
	}
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

type Client struct {
	*maps.Safe
	ExitOption []ExitOption
	BindMap    map[string]http.HandlerFunc //自定义接口绑定

	Succ         func(data interface{}, count ...int64)
	Fail         func(data interface{})
	Unauthorized func()
	Forbidden    func()
}

// SetHandlerWithCode 设置响应成功失败等
func (this *Client) SetHandlerWithCode(succ, fail, unauthorized, forbidden interface{}) *Client {
	this.Succ = this.NewSuccWithCode(succ)
	this.Fail = this.NewFailWithCode(fail)
	this.Unauthorized = this.NewUnauthorizedWithCode(unauthorized)
	this.Forbidden = this.NewForbiddenWithCode(forbidden)
	return this
}

func (this *Client) Bind(path string, handler http.HandlerFunc) *Client {
	this.BindMap[path] = handler
	return this
}

// SetExitOption 设置退出选项
func (this *Client) SetExitOption(f ...ExitOption) *Client {
	this.ExitOption = append(this.ExitOption, f...)
	return this
}

//=================================Response=================================//

func (this *Client) Redirect(httpCode int, url string) {
	this.NewExit(httpCode, &TEXT{}).SetHeader("Location", url).Exit()
}

func (this *Client) Custom(httpCode int, data interface{}, header http.Header) {
	this.NewExit(httpCode, &CUSTOM{Data: data, Headers: header}).Exit()
}

func (this *Client) Text(httpCode int, data interface{}) {
	this.NewExit(httpCode, &TEXT{Data: data}).Exit()
}

func (this *Client) Html(httpCode int, data interface{}) {
	this.NewExit(httpCode, &HTML{Data: data}).Exit()
}

// Json 返回json退出
func (this *Client) Json(httpCode int, data interface{}) {
	this.NewExit(httpCode, &JSON{Data: data}).Exit()
}

func (this *Client) File(name string, size int64, r io.ReadCloser) {
	this.NewExit(http.StatusOK, &FILE{
		Name:       name,
		Size:       size,
		ReadCloser: r}).Exit()
}

func (this *Client) Reader(httpCode int, r io.ReadCloser) {
	this.NewExit(httpCode, &READER{ReadCloser: r}).Exit()
}

// NewExit 自定义退出
func (this *Client) NewExit(httpCode int, i IMarshal) *Exit {
	return NewExit(httpCode, i, this.ExitOption...)
}

// Exit 直接退出,设置的跨域啥的应该是无效的
func (this *Client) Exit() {
	NewExit(-1, nil).Exit()
}

//=================================Proxy=================================//

func (this *Client) Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	defer r.Body.Close()
	req, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		this.Fail(err)
		return
	}
	req.Header = r.Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		this.Fail(err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

//=================================Middle=================================//

func (this *Client) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				this.MiddleRecover(e, w)
			}
		}()
		if handler, ok := this.BindMap[r.URL.Path]; ok {
			if handler != nil {
				handler(w, r)
			}
			return
		}
		h.ServeHTTP(w, r)
	})
}

// MiddleRecover 例gf等web框架只需要这一半即可,但是Bind会失效
func (this *Client) MiddleRecover(e interface{}, w http.ResponseWriter) {
	switch w2 := e.(type) {
	case *Exit:
		w2.WriteTo(w)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		s := conv.String(e)
		w.Write(*(*[]byte)(unsafe.Pointer(&s)))
	}
}
