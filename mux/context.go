package mux

import (
	"context"
	"encoding/json"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame/middle/in"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	_nil = conv.Nil()
)

type Ctx interface {
	conv.Extend
	in.Respondent
	Request() *http.Request
	Writer() http.ResponseWriter
	Context() context.Context
	SetContext(context.Context)
	Next()
	GetBodyBytes() []byte
	GetBodyString() string
	GetHeader(key string) string
	Parse(ptr any)
}

func NewCtx(w http.ResponseWriter, r *http.Request, res in.Respondent) Ctx {
	req := &Request{
		writer:     w,
		request:    r,
		Respondent: res,
		queryForm:  r.URL.Query(),
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
	writer  http.ResponseWriter
	request *http.Request
	conv.Extend
	in.Respondent

	queryForm url.Values     //解析后的query参数
	jsonFrom  map[string]any //解析body后的json
	cache     *maps.Safe     //缓存其他数据
	body      *[]byte        //缓存body数据
	next      func()
}

func (this *Request) Next() {
	if this.next != nil {
		this.next()
	}
}

func (this *Request) Request() *http.Request {
	return this.request
}

func (this *Request) Writer() http.ResponseWriter {
	return this.writer
}

func (this *Request) Context() context.Context {
	return this.request.Context()
}

func (this *Request) SetContext(ctx context.Context) {
	this.request.WithContext(ctx)
}

func (this *Request) GetBodyBytes() []byte {
	if this.body != nil {
		return *this.body
	}
	defer this.Request().Body.Close()
	bs, err := io.ReadAll(this.Request().Body)
	in.CheckErr(err)
	this.body = &bs
	return bs
}

func (this *Request) GetBodyString() string {
	return string(this.GetBodyBytes())
}

func (this *Request) Cache() *maps.Safe {
	if this.cache == nil {
		this.cache = maps.NewSafe()
	}
	return this.cache
}

func (this *Request) GetHeader(key string) string {
	return this.Request().Header.Get(key)
}

func (this *Request) Parse(ptr any) {
	if this == nil || this.Request() == nil {
		return
	}

	//multipart/form-data
	if strings.Contains(this.GetHeader("Content-Type"), "multipart/form-data") {
		//通过form-data解析
		if this.Request().Form == nil {
			if this.Request().ParseMultipartForm(1<<20) == nil {
				m := map[string]any{}
				for k, v := range this.Request().Form {
					m[k] = v[0]
				}
				err := conv.Unmarshal(this.Request().Form, ptr)
				in.CheckErr(err)
			}
		}
	} else {
		//通过json解析
		defer this.Request().Body.Close()
		bs, err := io.ReadAll(this.Request().Body)
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
	if this.cache != nil {
		return this.cache.GetVar(key)
	}

	return _nil
}

func (this *Request) GetQueryGMap() map[string]any {
	if this == nil || this.Request == nil {
		return nil
	}
	m := map[string]any{}
	for k, v := range this.queryForm {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
	}
	return m
}

func (this *Request) GetQueryVar(key string) *conv.Var {
	if this == nil || this.Request == nil {
		return _nil
	}
	ls, ok := this.queryForm[key]
	if !ok || len(ls) == 0 {
		return _nil
	}
	return conv.New(ls[0])
}

func (this *Request) parseJsonForm() error {
	if this.Request().Body != nil {
		bs, err := io.ReadAll(this.Request().Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, &this.jsonFrom)
	}
	return nil
}

func (this *Request) GetBodyVar(key string) *conv.Var {
	if this == nil || this.Request == nil {
		return _nil
	}
	if strings.Contains(this.GetHeader("Content-Type"), "application/json") {
		if this.jsonFrom == nil {
			this.parseJsonForm()
		}
		if this.jsonFrom != nil {
			if val, ok := this.jsonFrom[key]; ok {
				return conv.New(val)
			}
		}
	}
	if this.Request().Form == nil {
		this.Request().ParseMultipartForm(1 << 20)
	}
	if this.Request().Form == nil {
		return _nil
	}
	ls, ok := this.Request().Form[key]
	if !ok || len(ls) == 0 {
		return _nil
	}
	return conv.New(ls[0])
}

func (this *Request) GetHeaderGMap() map[string]any {
	if this == nil || this.Request() == nil {
		return nil
	}
	m := map[string]any{}
	for k, v := range this.Request().Header {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
	}
	return m
}

func (this *Request) GetHeaderVar(key string) *conv.Var {
	if this == nil || this.Request() == nil || this.Request().Header == nil {
		return conv.Nil()
	}
	ls, ok := this.Request().Header[key]
	if !ok || len(ls) == 0 {
		return _nil
	}
	return conv.New(ls[0])
}
