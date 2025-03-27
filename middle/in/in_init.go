package in

import (
	"bytes"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
	"time"
)

var DefaultClient = New(WithDefault())

// SetCacheByHandler 尝试从缓存中获取数据,如果不存在则通过函数获取,执行函数时,其他相同的key会等待此次结果
func SetCacheByHandler(key interface{}, handler func() interface{}, expiration time.Duration) interface{} {
	value, err := DefaultClient.GetOrSetByHandler(key, func() (interface{}, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

// DelCache 删除缓存数据
func DelCache(key ...interface{}) {
	for _, v := range key {
		DefaultClient.Del(v)
	}
}

// SetCache 设置缓存,覆盖缓存
func SetCache(key interface{}, value interface{}, expiration time.Duration) {
	DefaultClient.Set(key, value, expiration)
}

// HTTPHandler 对http.HandlerFunc使用中间件,http.ListenAndServe(":80",in.HTTPHandler(h))
func HTTPHandler(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return DefaultClient.Recover(http.HandlerFunc(h))
}

// Recover 对http.Handler使用中间件
func Recover(h http.Handler) http.Handler {
	return DefaultClient.Recover(h)
}

// RecoverFunc 对http.Handler使用中间件
func RecoverFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return DefaultClient.Recover(http.HandlerFunc(f))
}

// MiddleRecover 捕捉panic,或自定义panic,并输出到http.ResponseWriter
func MiddleRecover(e interface{}, w http.ResponseWriter) {
	DefaultClient.MiddleRecover(e, w)
}

// SetHandlerWithCode 设置常用响应的状态码
func SetHandlerWithCode(succ, fail, unauthorized, forbidden interface{}) *Client {
	return DefaultClient.SetHandlerWithCode(succ, fail, unauthorized, forbidden)
}

//=================================Return=================================//

func Custom(code int, data interface{}, header http.Header) { DefaultClient.Custom(code, data, header) }

func Custom200(data interface{}, h http.Header) { DefaultClient.Custom(http.StatusOK, data, h) }

func Text(code int, data interface{}) { DefaultClient.Text(code, data) }

func Text200(data interface{}) { DefaultClient.Text(http.StatusOK, data) }

func Html(code int, data interface{}) { DefaultClient.Html(code, data) }

func Html200(data interface{}) { DefaultClient.Html(http.StatusOK, data) }

func Json(httpCode int, data interface{}) { DefaultClient.Json(httpCode, data) }

func Json200(data interface{}) { Json(http.StatusOK, data) }

func Json400(data interface{}) { Json(http.StatusBadRequest, data) }

func Json401(data interface{}) { Json(http.StatusUnauthorized, data) }

func Json403(data interface{}) { Json(http.StatusForbidden, data) }

func Json404(data interface{}) { Json(http.StatusNotFound, data) }

func Json415(data interface{}) { Json(http.StatusUnsupportedMediaType, data) }

func Json500(data interface{}) { Json(http.StatusInternalServerError, data) }

func Reader(code int, r io.ReadCloser) { DefaultClient.Reader(code, r) }

func Reader200(r io.ReadCloser) { DefaultClient.Reader(http.StatusOK, r) }

// Redirect301 永久重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect301(addr string) { DefaultClient.Redirect(http.StatusMovedPermanently, addr) }

// Redirect302 临时重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect302(addr string) { DefaultClient.Redirect(http.StatusFound, addr) }

// Redirect307 Temporary Redirect临时重定向,不改变请求方式
func Redirect307(addr string) { DefaultClient.Redirect(http.StatusTemporaryRedirect, addr) }

// Redirect308 Permanent Redirect永久重定向,不改变请求方式
func Redirect308(addr string) { DefaultClient.Redirect(http.StatusPermanentRedirect, addr) }

//=================================File=================================//

// FileLocal 返回本地文件
func FileLocal(name, filename string) {
	f, err := os.Open(filename)
	CheckErr(err)
	i, err := f.Stat()
	CheckErr(err)
	DefaultClient.File(name, i.Size(), f)
}

// FileReader 返回文件流
func FileReader(name string, r io.ReadCloser) {
	DefaultClient.File(name, -1, r)
}

// FileBytes 返回文件,字节
func FileBytes(name string, bs []byte) {
	DefaultClient.File(name, int64(len(bs)), io.NopCloser(bytes.NewReader(bs)))
}

func Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	DefaultClient.Proxy(w, r, uri)
}

//=================================Injoy=================================//

var (

	// Succ 响应成功
	Succ = DefaultClient.Succ

	// Fail 响应失败
	Fail = DefaultClient.Fail

	// Unauthorized 响应验证失败
	Unauthorized = DefaultClient.Unauthorized

	// Forbidden 响应没有权限
	Forbidden = DefaultClient.Forbidden
)

// Err 退出,并校验错误
func Err(data interface{}, succData ...interface{}) {
	if data == nil {
		Succ(conv.DefaultInterface(nil, succData...))
	} else {
		Fail(data)
	}
}

// CheckErr 检测错误(遇到错误结束)
func CheckErr(err error, failMsg ...string) {
	if err != nil {
		Err(conv.DefaultString(err.Error(), failMsg...))
	}
}
