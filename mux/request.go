package mux

import (
	"encoding/json"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame/middle/in"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Handler func(r *Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(NewRequest(w, r))
}

func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	req := &Request{
		Writer:    w,
		Request:   r,
		QueryForm: r.URL.Query(),
	}
	req.Extend = conv.NewExtend(req)

	//尝试获取中间件的cache
	if val := r.Context().Value("_cache"); val != nil {
		if cache, ok := val.(*maps.Safe); ok {
			req.cache = cache
		}
	}

	return req
}

type Request struct {
	Writer http.ResponseWriter
	*http.Request
	conv.Extend
	QueryForm url.Values             //解析后的query参数
	JsonFrom  map[string]interface{} //解析body后的json
	cache     *maps.Safe
	body      *[]byte
}

func (this *Request) Exit() {
	in.DefaultClient.Exit()
}

func (this *Request) GetBodyBytes() []byte {
	if this.body != nil {
		return *this.body
	}
	defer this.Body.Close()
	bs, err := io.ReadAll(this.Body)
	in.CheckErr(err)
	this.body = &bs
	return bs
}

func (this *Request) GetBodyString() string {
	return string(this.GetBodyBytes())
}

func (this *Request) GetRequest() *http.Request {
	return this.Request
}

func (this *Request) SetCache(key string, value interface{}) {
	if this.cache == nil {
		this.cache = maps.NewSafe()
	}
	this.cache.Set(key, value)
}

func (this *Request) GetCache(key string) *conv.Var {
	if this.cache == nil {
		return conv.Nil()
	}
	return this.cache.GetVar(key)
}

func (this *Request) Parse(ptr interface{}) {
	if this == nil || this.Request == nil {
		return
	}

	//multipart/form-data
	if strings.Contains(this.Header.Get("Content-Type"), "multipart/form-data") {
		//通过form-data解析
		if this.Request.Form == nil {
			if this.Request.ParseMultipartForm(1<<20) == nil {
				m := map[string]interface{}{}
				for k, v := range this.Request.Form {
					m[k] = v[0]
				}
				err := conv.Unmarshal(this.Request.Form, ptr)
				in.CheckErr(err)
			}
		}
	} else {
		//通过json解析
		defer this.Body.Close()
		bs, err := io.ReadAll(this.Body)
		in.CheckErr(err)
		err = conv.Unmarshal(bs, ptr)
		in.CheckErr(err)
	}

	//先尝试从header获取参数,也就是说改优先级最低
	if m := this.GetHeaderGMap(); len(m) > 0 {
		err := conv.Unmarshal(m, ptr)
		in.CheckErr(err)
	}

	//再尝试从url获取
	if m := this.GetQueryGMap(); len(m) > 0 {
		err := conv.Unmarshal(m, ptr)
		in.CheckErr(err)
	}

}

func (this *Request) GetVar(key string) *conv.Var {

	//先从query获取参数
	v := this.GetQueryVar(key)
	if !v.IsNil() {
		return v
	}

	//再从body获取参数
	v = this.GetBodyVar(key)
	if !v.IsNil() {
		return v
	}

	//再从header获取参数
	v = this.GetHeaderVar(key)
	if !v.IsNil() {
		return v
	}

	//最后尝试从cache获取参数
	return this.GetCache(key)
}

func (this *Request) GetQueryGMap() map[string]interface{} {
	if this == nil || this.Request == nil {
		return nil
	}
	m := map[string]interface{}{}
	for k, v := range this.QueryForm {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
	}
	return m
}

func (this *Request) GetQueryVar(key string) *conv.Var {
	if this == nil || this.Request == nil {
		return conv.Nil()
	}
	ls, ok := this.QueryForm[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
}

func (this *Request) parseJsonForm() error {
	if this.Body != nil {
		bs, err := io.ReadAll(this.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, &this.JsonFrom)
	}
	return nil
}

func (this *Request) GetBodyVar(key string) *conv.Var {
	if this == nil || this.Request == nil {
		return conv.Nil()
	}
	if strings.Contains(this.Header.Get("Content-Type"), "application/json") {
		if this.JsonFrom == nil {
			this.parseJsonForm()
		}
		if this.JsonFrom != nil {
			if val, ok := this.JsonFrom[key]; ok {
				return conv.New(val)
			}
		}
	}
	if this.Request.Form == nil {
		this.Request.ParseMultipartForm(1 << 20)
	}
	if this.Request.Form == nil {
		return conv.Nil()
	}
	ls, ok := this.Request.Form[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
}

func (this *Request) GetHeaderGMap() map[string]interface{} {
	if this == nil || this.Request == nil {
		return nil
	}
	m := map[string]interface{}{}
	for k, v := range this.Request.Header {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
	}
	return m
}

func (this *Request) GetHeaderVar(key string) *conv.Var {
	if this == nil || this.Request == nil || this.Request.Header == nil {
		return conv.Nil()
	}
	ls, ok := this.Request.Header[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
}

func (this *Request) GetHeader(key string) string {
	return this.Request.Header.Get(key)
}

func (this *Request) WriteTo(w io.Writer) error {
	return this.Request.Write(w)
}

func (this *Request) Write(p []byte) (int, error) {
	return this.Writer.Write(p)
}

func (this *Request) WriteJson(v interface{}) error {
	return json.NewEncoder(this.Writer).Encode(v)
}

func (this *Request) WriteAny(v interface{}) error {
	_, err := this.Writer.Write(conv.Bytes(v))
	return err
}

func (this *Request) SetHeader(k, v string) {
	this.Writer.Header().Set(k, v)
}
