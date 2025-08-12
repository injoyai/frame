package in

import (
	"bytes"
	"encoding/json"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

type Writer interface {
	io.ReadCloser
	http.ResponseWriter
	StatusCode() int
	Header() http.Header
	Write(bs []byte) (int, error)
	WriteAny(v any) error
	WriteJson(v any) error
	WriteHeader(code int)
	AddHeader(i string, v ...string)
	SetHeader(i string, v ...string)
	SetCORS()
	SetContentType(ct ...string)
	SetContentTypeJson()
	WriteTo(w http.ResponseWriter)
	Exit()
}

type WriterOption func(w Writer)

func NewWriter(code int, i IMarshal, op ...WriterOption) Writer {
	buf := bytes.NewBuffer(nil)
	w := &writer{
		code:   code,
		header: http.Header{},
		ReadCloser: struct {
			io.Reader
			io.Closer
		}{
			io.MultiReader(i, buf),
			i,
		},
		Writer: buf,
	}
	if i != nil && i.Header() != nil {
		w.header = i.Header()
	}
	for _, v := range op {
		v(w)
	}
	return w
}

type writer struct {
	code          int         //响应状态码
	header        http.Header //响应请求头
	io.ReadCloser             //响应内容,body
	io.Writer                 //写入body
}

func (this *writer) StatusCode() int {
	return this.code
}

func (this *writer) Header() http.Header {
	return this.header
}

func (this *writer) WriteHeader(code int) {
	this.code = code
}

func (this *writer) AddHeader(i string, v ...string) {
	this.header[i] = append(this.header[i], v...)
}

func (this *writer) SetHeader(i string, v ...string) {
	this.header[i] = v
}

func (this *writer) SetCORS() {
	this.SetHeader("Access-Control-Allow-Origin", "*")
	this.SetHeader("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
	this.SetHeader("Access-Control-Allow-Credentials", "true")
	this.SetHeader("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
	this.SetHeader("Access-Control-Allow-Max-Age", "3600")
}

func (this *writer) SetContentType(ct ...string) {
	this.SetHeader("Content-Type", ct...)
}

func (this *writer) SetContentTypeJson() {
	this.SetHeader("Content-Type", "application/json;charset=utf-8")
}

func (this *writer) WriteAny(v any) error {
	s := conv.String(v)
	_, err := this.Writer.Write(*(*[]byte)(unsafe.Pointer(&s)))
	return err
}

func (this *writer) WriteJson(v any) error {
	return json.NewEncoder(this.Writer).Encode(v)
}

func (this *writer) WriteTo(w http.ResponseWriter) {
	if ww, ok := w.(*writer); ok {
		*ww = *this
		return
	}
	if this.header != nil {
		for i, v := range this.header {
			w.Header().Set(i, strings.Join(v, ","))
		}
	}
	if this.code > 0 {
		w.WriteHeader(this.code)
	}
	io.Copy(w, this)
	this.Close()
}

func (this *writer) Exit() {
	panic(this)
}
