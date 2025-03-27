package in

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type ExitOption func(e *Exit)

func NewExit(httpCode int, i IMarshal, op ...ExitOption) *Exit {
	buf := bytes.NewBuffer(nil)
	e := &Exit{
		Code:    httpCode,
		Headers: http.Header{},
		Body: struct {
			io.Reader
			io.Closer
		}{
			io.MultiReader(i, buf),
			i,
		},
		buf: buf,
	}
	for _, v := range op {
		v(e)
	}
	if i != nil {
		e.Headers = i.Header()
	}
	return e
}

var _ http.ResponseWriter = (*Exit)(nil)

type Exit struct {
	Code    int           `json:"code"`   //响应状态码
	Headers http.Header   `json:"header"` //响应请求头
	Body    io.ReadCloser `json:"body"`   //响应内容,body
	buf     *bytes.Buffer //补充的body
}

func (this *Exit) Header() http.Header {
	return this.Headers
}

// Write 写入数据,实现http.ResponseWriter
func (this *Exit) Write(bs []byte) (int, error) {
	if this.buf == nil {
		this.buf = bytes.NewBuffer(bs)
		return len(bs), nil
	}
	return this.buf.Write(bs)
}

// WriteHeader 设置状态码,实现http.ResponseWriter
func (this *Exit) WriteHeader(statusCode int) {
	this.Code = statusCode
}

func (this *Exit) AddHeader(i string, v ...string) *Exit {
	this.Headers[i] = append(this.Headers[i], v...)
	return this
}

// SetHeader 设置请求头
func (this *Exit) SetHeader(i string, v ...string) *Exit {
	this.Headers[i] = v
	return this
}

// SetContentType 设置请求头Content-Type
func (this *Exit) SetContentType(ct ...string) *Exit {
	return this.SetHeader("Content-Type", ct...)
}

// SetHeaderJson 设置请求头Content-Type
func (this *Exit) SetHeaderJson() *Exit {
	return this.SetHeader("Content-Type", "application/json;charset=utf-8")
}

// SetHeaderCORS 设置跨域
func (this *Exit) SetHeaderCORS() *Exit {
	this.SetHeader("Access-Control-Allow-Origin", "*")
	this.SetHeader("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
	this.SetHeader("Access-Control-Allow-Credentials", "true")
	this.SetHeader("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
	this.SetHeader("Access-Control-Allow-Max-Age", "3600")
	return this
}

// Exit 退出程序,中断执行,需要和recover配合使用
func (this *Exit) Exit() {
	panic(this)
}

// WriteTo 写入响应
// 这里要先设置Header,再设置Code,否则Header可能无效(例mux)
func (this *Exit) WriteTo(w http.ResponseWriter) {
	if ww, ok := w.(*Exit); ok {
		*ww = *this
		return
	}
	for i, v := range this.Headers {
		w.Header().Set(i, strings.Join(v, ","))
	}
	if this.Code >= 0 {
		w.WriteHeader(this.Code)
	}
	io.Copy(w, this.Body)
	this.Body.Close()
}
