package in

import (
	"io"
	"net/http"
	"strings"

	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/codec"
)

func GetVar(r *http.Request, key string) *conv.Var {
	//尝试从query中获取
	if v := GetQueryVar(r, key); !v.IsNil() {
		return v
	}
	//优先从body获取
	if v := GetBodyVar(r, key); !v.IsNil() {
		return v
	}
	//尝试从header中获取
	if v := GetHeaderVar(r, key); !v.IsNil() {
		return v
	}
	return conv.Nil()
}

func GetQueryVar(r *http.Request, key string) *conv.Var {
	ls, ok := r.URL.Query()[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
}

func GetBodyVar(r *http.Request, key string) *conv.Var {

	//通过json解析
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, _ := r.GetBody()
		defer body.Close()
		return conv.NewMap(body, codec.Json).GetVar(key)
	}

	if r.Form == nil {
		r.ParseForm()
	}

	//尝试从from中获取
	if r.Form != nil {
		if ls, ok := r.Form[key]; ok && len(ls) > 0 {
			return conv.New(ls[0])
		}
	}

	return conv.Nil()
}

func GetBodyMap(r *http.Request, codec ...codec.Interface) *conv.Map {
	body, _ := r.GetBody()
	defer body.Close()
	return conv.NewMap(body, codec...)
}

func GetHeaderVar(r *http.Request, key string) *conv.Var {
	ls, ok := r.Header[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
}

func GetPageNum(r *http.Request, key string, _default ...int) int {
	if v, ok := r.URL.Query()[key]; ok && len(v) > 0 {
		return conv.Int(v[0]) - 1
	}
	if len(_default) > 0 {
		return _default[0]
	}
	return 0
}

func GetPageSize(r *http.Request, key string, _default ...int) int {
	if v, ok := r.URL.Query()[key]; ok && len(v) > 0 {
		return conv.Int(v)
	}
	if len(_default) > 0 {
		return _default[0]
	}
	return 10
}

/*



 */

func NewRequest(r *http.Request) *Request {
	x := &Request{
		Request: r,
	}
	x.Extend = conv.NewExtend(x)
	return x
}

type Request struct {
	*http.Request
	bodyMap *conv.Map  //解析后的json body
	parsed  bool       //是否解析了body
	cache   *maps.Safe //自定义缓存
	conv.Extend
}

func (this *Request) GetQueryVar(key string) *conv.Var {
	if this == nil || this.Request == nil || this.Request.URL == nil {
		return conv.Nil()
	}
	ls, ok := this.Request.URL.Query()[key]
	if !ok || len(ls) == 0 {
		return conv.Nil()
	}
	return conv.New(ls[0])
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

func (this *Request) GetBodyVar(key string) *conv.Var {
	if !this.parsed {
		//通过json解析
		if strings.Contains(this.Header.Get("Content-Type"), "application/json") {
			bs, _ := io.ReadAll(this.Body)
			this.Body.Close()
			this.bodyMap = conv.NewMap(bs)
		} else if this.Request.Form == nil {
			this.Request.ParseMultipartForm(1 << 20)
		}
		this.parsed = true
	}

	//优先从body获取
	if this.bodyMap != nil {
		if val := this.bodyMap.GetVar(key); !val.IsNil() {
			return val
		}
	}

	//尝试从from中获取
	if this.Request.Form != nil {
		if ls, ok := this.Request.Form[key]; ok && len(ls) > 0 {
			return conv.New(ls[0])
		}
	}
	return conv.Nil()
}

func (this *Request) GetCacheVar(key string) *conv.Var {
	if this.cache == nil {
		return conv.Nil()
	}
	return this.cache.GetVar(key)
}

func (this *Request) GetVar(key string) *conv.Var {
	//尝试从query中获取
	if v := this.GetQueryVar(key); !v.IsNil() {
		return v
	}
	//优先从body获取
	if v := this.GetBodyVar(key); !v.IsNil() {
		return v
	}
	//尝试从header中获取
	if v := this.GetHeaderVar(key); !v.IsNil() {
		return v
	}
	return this.GetCacheVar(key)
}

// GetFile 获取上传的文件流
func (this *Request) GetFile(name string) io.ReadCloser {
	f, _, err := this.Request.FormFile(name)
	CheckErr(err)
	return f
}

func (this *Request) ParseJsonBody(ptr any) {
	defer this.Body.Close()
	bs, err := io.ReadAll(this.Body)
	CheckErr(err)
	if err = conv.Unmarshal(bs, ptr); err != nil {
		Json415(err)
	}
}

func (this *Request) ParseBody(ptr any) {
	if strings.Contains(this.Header.Get("Content-Type"), "application/json") {
		defer this.Body.Close()
		bs, err := io.ReadAll(this.Body)
		CheckErr(err)
		if err = conv.Unmarshal(bs, ptr); err != nil {
			Json415(err)
		}
		return
	}
	//如果不是json,则使用自带的form解析
	if this.Request.Form == nil {
		this.Request.ParseMultipartForm(1 << 20)
	}
	if this.Request.Form != nil {
		if err := conv.Unmarshal(this.Request.Form, ptr); err != nil {
			Json415(err)
		}
	}
}

func (this *Request) Parse(ptr any) {
	if this == nil || this.Request == nil {
		return
	}

	//先尝试从header获取参数,也就是说改优先级最低
	if this.Request.Header != nil {
		if err := conv.Unmarshal(this.Request.Header, ptr); err != nil {
			Json415(err)
		}
	}

	//再尝试从url获取
	if this.Request.URL != nil {
		if err := conv.Unmarshal(this.URL.Query(), ptr); err != nil {
			Json415(err)
		}
	}

	//通过json解析
	if strings.Contains(this.Header.Get("Content-Type"), "application/json") {
		defer this.Body.Close()
		bs, err := io.ReadAll(this.Body)
		CheckErr(err)
		if err = conv.Unmarshal(bs, ptr); err != nil {
			Json415(err)
		}
		return
	}

	//如果不是json,则使用自带的form解析
	if this.Request.Form == nil {
		this.Request.ParseMultipartForm(1 << 20)
	}
	if this.Request.Form != nil {
		if err := conv.Unmarshal(this.Request.Form, ptr); err != nil {
			Json415(err)
		}
	}
}
