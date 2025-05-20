package in

import (
	"bytes"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
)

var _ Respondent = (*client)(nil)

//=================================Response=================================//

func (this *client) Custom(httpCode int, data any, header http.Header) {
	this.NewWriter(httpCode, &CUSTOM{Data: data, Headers: header}).Exit()
}

func (this *client) Custom200(data any, header http.Header) {
	this.Custom(http.StatusOK, data, header)
}

func (this *client) Text(httpCode int, data any) {
	this.NewWriter(httpCode, &TEXT{Data: data}).Exit()
}

func (this *client) Text200(data any) {
	this.Text(http.StatusOK, data)
}

func (this *client) Html(httpCode int, data any) {
	this.NewWriter(httpCode, &HTML{Data: data}).Exit()
}

func (this *client) Html200(data any) {
	this.Html(http.StatusOK, data)
}

func (this *client) HtmlLocal(filename string) {
	f, err := os.Open(filename)
	this.CheckErr(err)
	this.Html(http.StatusOK, f)
}

func (this *client) HtmlReader(r io.ReadCloser) {
	this.Html(http.StatusOK, r)
}

// Json 返回json退出
func (this *client) Json(httpCode int, data any) {
	this.NewWriter(httpCode, &JSON{Data: data}).Exit()
}

func (this *client) Json200(data any) {
	this.Json(http.StatusOK, data)
}

func (this *client) Json400(data any) {
	this.Json(http.StatusBadRequest, data)
}

func (this *client) Json401(data any) {
	this.Json(http.StatusUnauthorized, data)
}

func (this *client) Json403(data any) {
	this.Json(http.StatusForbidden, data)
}

func (this *client) Json415(data any) {
	this.Json(http.StatusUnsupportedMediaType, data)
}

func (this *client) Json500(data any) {
	this.Json(http.StatusInternalServerError, data)
}

func (this *client) Reader(httpCode int, r io.ReadCloser) {
	this.NewWriter(httpCode, &READER{ReadCloser: r}).Exit()
}

func (this *client) Reader200(r io.ReadCloser) {
	this.Reader(http.StatusOK, r)
}

func (this *client) File(name string, size int64, r io.ReadCloser) {
	this.NewWriter(http.StatusOK, &FILE{
		Name:       name,
		Size:       size,
		ReadCloser: r}).Exit()
}

func (this *client) FileLocal(name, filename string) {
	f, err := os.Open(filename)
	this.CheckErr(err)
	i, err := f.Stat()
	this.CheckErr(err)
	this.File(name, i.Size(), f)
}

// FileReader 返回文件流
func (this *client) FileReader(name string, r io.ReadCloser) {
	this.File(name, -1, r)
}

func (this *client) FileBytes(name string, bs []byte) {
	this.File(name, int64(len(bs)), io.NopCloser(bytes.NewReader(bs)))
}

func (this *client) Redirect(httpCode int, url string) {
	w := this.NewWriter(httpCode, &TEXT{})
	w.SetHeader("Location", url)
	w.Exit()
}

// Redirect301 永久重定向,GET和HEAD自动请求,其他需要用户确认
func (this *client) Redirect301(addr string) {
	this.Redirect(http.StatusMovedPermanently, addr)
}

// Redirect302 临时重定向,GET和HEAD自动请求,其他需要用户确认
func (this *client) Redirect302(addr string) {
	this.Redirect(http.StatusFound, addr)
}

// Redirect307 Temporary Redirect临时重定向,不改变请求方式
func (this *client) Redirect307(addr string) {
	this.Redirect(http.StatusTemporaryRedirect, addr)
}

// Redirect308 Permanent Redirect永久重定向,不改变请求方式
func (this *client) Redirect308(addr string) {
	this.Redirect(http.StatusPermanentRedirect, addr)
}

// NewWriter 自定义退出
func (this *client) NewWriter(code int, i IMarshal) Writer {
	return NewWriter(code, i, this.writerOptions...)
}

// Exit 直接退出,设置的跨域啥的应该是无效的
func (this *client) Exit() {
	NewWriter(-1, nil).Exit()
}

//=============

// Succ 响应成功
func (this *client) Succ(data any, count ...int64) {
	this.succ(data, count...)
}

// Fail 响应失败
func (this *client) Fail(data any) {
	this.fail(data)
}

// Unauthorized 响应验证失败
func (this *client) Unauthorized() {
	this.unauthorized()
}

// Forbidden 响应没有权限
func (this *client) Forbidden() {
	this.forbidden()
}

// Err 退出,并校验错误
func (this *client) Err(data any) {
	if data == nil {
		this.Succ(nil)
	} else {
		this.Fail(data)
	}
}

// CheckErr 检测错误(遇到错误结束)
func (this *client) CheckErr(err error, failMsg ...string) {
	if err != nil {
		this.Err(conv.Default[string](err.Error(), failMsg...))
	}
}
