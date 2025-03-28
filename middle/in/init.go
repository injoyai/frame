package in

import (
	"io"
	"net/http"
	"time"
)

var DefaultClient = New(WithDefault())

// SetCacheByHandler 尝试从缓存中获取数据,如果不存在则通过函数获取,执行函数时,其他相同的key会等待此次结果
func SetCacheByHandler(key any, handler func() any, expiration time.Duration) interface{} {
	value, err := DefaultClient.Cache().GetOrSetByHandler(key, func() (any, error) { return handler(), nil }, expiration)
	DefaultClient.CheckErr(err)
	return value
}

// DelCache 删除缓存数据
func DelCache(key ...any) {
	for _, v := range key {
		DefaultClient.Cache().Del(v)
	}
}

// SetCache 设置缓存,覆盖缓存
func SetCache(key any, value any, expiration ...time.Duration) {
	DefaultClient.Cache().Set(key, value, expiration...)
}

// Recover 对http.Handler使用中间件
func Recover(h http.Handler) http.Handler {
	return DefaultClient.Recover(h)
}

// RecoverFunc 对http.Handler使用中间件
func RecoverFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return DefaultClient.RecoverFunc(f)
}

// RecoverWriter 捕捉panic,或自定义panic,并输出到http.ResponseWriter
func RecoverWriter(e any, w http.ResponseWriter) { DefaultClient.RecoverWriter(e, w) }

// SetHandlerWithCode 设置常用响应的状态码
func SetHandlerWithCode(succ, fail, unauthorized, forbidden any) {
	DefaultClient.SetHandlerWithCode(succ, fail, unauthorized, forbidden)
}

//=================================Return=================================//

func Custom(code int, data any, header http.Header) { DefaultClient.Custom(code, data, header) }

func Custom200(data any, h http.Header) { DefaultClient.Custom(http.StatusOK, data, h) }

func Text(code int, data any) { DefaultClient.Text(code, data) }

func Text200(data any) { DefaultClient.Text(http.StatusOK, data) }

func Html(code int, data any) { DefaultClient.Html(code, data) }

func Html200(data any) { DefaultClient.Html(http.StatusOK, data) }

func HtmlLocal(filename string) { DefaultClient.HtmlLocal(filename) }

func HtmlReader(r io.ReadCloser) { DefaultClient.Html(http.StatusOK, r) }

func Json(httpCode int, data any) { DefaultClient.Json(httpCode, data) }

func Json200(data any) { Json(http.StatusOK, data) }

func Json400(data any) { Json(http.StatusBadRequest, data) }

func Json401(data any) { Json(http.StatusUnauthorized, data) }

func Json403(data any) { Json(http.StatusForbidden, data) }

func Json404(data any) { Json(http.StatusNotFound, data) }

func Json415(data any) { Json(http.StatusUnsupportedMediaType, data) }

func Json500(data any) { DefaultClient.Json500(data) }

func Reader(code int, r io.ReadCloser) { DefaultClient.Reader(code, r) }

func Reader200(r io.ReadCloser) { DefaultClient.Reader200(r) }

// Redirect301 永久重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect301(addr string) { DefaultClient.Redirect301(addr) }

// Redirect302 临时重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect302(addr string) { DefaultClient.Redirect302(addr) }

// Redirect307 Temporary Redirect临时重定向,不改变请求方式
func Redirect307(addr string) { DefaultClient.Redirect307(addr) }

// Redirect308 Permanent Redirect永久重定向,不改变请求方式
func Redirect308(addr string) { DefaultClient.Redirect308(addr) }

//=================================File=================================//

// FileLocal 返回本地文件
func FileLocal(name, filename string) { DefaultClient.FileLocal(name, filename) }

// FileReader 返回文件流
func FileReader(name string, r io.ReadCloser) { DefaultClient.FileReader(name, r) }

// FileBytes 返回文件,字节
func FileBytes(name string, bs []byte) { DefaultClient.FileBytes(name, bs) }

func Proxy(w http.ResponseWriter, r *http.Request, uri string) { DefaultClient.Proxy(w, r, uri) }

//=================================Injoy=================================//

// Succ 响应成功
func Succ(data any, count ...int64) { DefaultClient.Succ(data, count...) }

// Fail 响应失败
func Fail(data any) { DefaultClient.Fail(data) }

// Unauthorized 响应验证失败
func Unauthorized() { DefaultClient.Unauthorized() }

// Forbidden 响应没有权限
func Forbidden() { DefaultClient.Forbidden() }

// Err 退出,并校验错误
func Err(data any) { DefaultClient.Err(data) }

// CheckErr 检测错误(遇到错误结束)
func CheckErr(err error, failMsg ...string) { DefaultClient.CheckErr(err, failMsg...) }
