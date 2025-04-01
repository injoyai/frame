package mux

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/injoyai/frame/middle/in"
	"io/fs"
	"net/http"
)

type (
	Handler func(c Ctx)
	Middle  = Handler

	//Option
	//Handler func(c Ctx)
	//func(s *Server)
	Option = any
)

type Grouper interface {
	http.Handler
	SetRespondent(respondent in.Respondent)
	Use(use ...Middle)
	Group(path string, handler func(g Grouper))
	ALL(path string, handler Handler)
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)
	DELETE(path string, handler Handler)
	OPTIONS(path string, handler Handler)
	HEAD(path string, handler Handler)
	CONNECT(path string, handler Handler)
	TRACE(path string, handler Handler)
	PATCH(path string, handler Handler)
}

type Group struct {
	*mux.Router          //路由实例
	prefix      string   //路由前缀
	middle      []Middle //中间件
	in.Respondent
}

func (this *Group) next(i int, r *Request) {
	if i >= len(this.middle) {
		this.Router.ServeHTTP(r.Writer(), r.request)
	} else {
		do := false
		r.next = func() { do = true; this.next(i+1, r) }
		this.middle[i](r)
		if !do {
			r.next()
		}
	}
}

func (this *Group) SetRespondent(r in.Respondent) {
	this.Respondent = r
}

func (this *Group) Use(f ...Middle) {
	this.middle = append(this.middle, f...)
}

func (this *Group) IgnoreMiddle(handler func(g Grouper)) Grouper {
	g := *this
	x := &g
	if handler != nil {
		handler(x)
	}
	return x
}

func (this *Group) Group(path string, handler func(g Grouper)) {
	g := &Group{
		Router:     this.Router,
		prefix:     this.prefix + path,
		middle:     this.middle,
		Respondent: this.Respondent,
	}
	if handler != nil {
		handler(g)
	}
}

// Static 放在最后执行,这个会占用Grouper下所有的路由
func (this *Group) Static(path string, dir string) {
	path = this.prefix + path
	s := http.StripPrefix(path, http.FileServer(http.Dir(dir)))
	this.Router.PathPrefix(path).Handler(s)
}

// StaticEmbed 放在最后执行,这个会占用Grouper下所有的路由
func (this *Group) StaticEmbed(path string, e embed.FS, dir string) error {
	web, err := fs.Sub(e, dir)
	if err != nil {
		return err
	}
	path = this.prefix + path
	s := http.StripPrefix(path, http.FileServer(http.FS(web)))
	this.Router.PathPrefix(path).Handler(s)
	return nil
}

// StaticFS 放在最后执行,这个会占用Grouper下所有的路由
func (this *Group) StaticFS(path string, fs fs.FS) error {
	path = this.prefix + path
	s := http.StripPrefix(path, http.FileServer(http.FS(fs)))
	this.Router.PathPrefix(path).Handler(s)
	return nil
}

func (this *Group) ALL(path string, handler Handler) {
	this.do([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodPatch, http.MethodConnect, http.MethodOptions, http.MethodTrace}, path, handler)
}

func (this *Group) GET(path string, handler Handler) {
	this.do([]string{http.MethodGet}, path, handler)
}

func (this *Group) POST(path string, handler Handler) {
	this.do([]string{http.MethodPost}, path, handler)
}

func (this *Group) PUT(path string, handler Handler) {
	this.do([]string{http.MethodPut}, path, handler)
}

func (this *Group) DELETE(path string, handler Handler) {
	this.do([]string{http.MethodDelete}, path, handler)
}

func (this *Group) OPTIONS(path string, handler Handler) {
	this.do([]string{http.MethodOptions}, path, handler)
}

func (this *Group) HEAD(path string, handler Handler) {
	this.do([]string{http.MethodHead}, path, handler)
}

func (this *Group) CONNECT(path string, handler Handler) {
	this.do([]string{http.MethodConnect}, path, handler)
}

func (this *Group) TRACE(path string, handler Handler) {
	this.do([]string{http.MethodTrace}, path, handler)
}

func (this *Group) PATCH(path string, handler Handler) {
	this.do([]string{http.MethodPatch}, path, handler)
}

func (this *Group) do(method []string, path string, handler Handler) {
	path = this.prefix + path
	middles := this.middle
	this.Router.Methods(method...).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c Ctx
		var ok bool
		if v := r.Context().Value("request"); v == nil {
			c = NewCtx(w, r, this.Respondent)
		} else if c, ok = v.(*Request); !ok {
			c = NewCtx(w, r, this.Respondent)
		}
		for _, f := range middles {
			f(c)
		}
		handler(c)
	})
}
