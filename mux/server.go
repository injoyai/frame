package mux

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/injoyai/frame"
	"github.com/injoyai/frame/middle"
	"github.com/injoyai/frame/middle/in"
	"io/fs"
	"net/http"
	"sync"
	"time"
)

type Option func(s *Server)

// WithCORS 设置跨域
func WithCORS() Option {
	return func(s *Server) {
		s.Use(func(r *Request, next func()) {
			next()
			middle.WithCORS(r.Writer)
		})
	}
}

// WithSwagger 设置swagger
func WithSwagger(swag *middle.Swagger) Option {
	return func(s *Server) {
		s.Use(func(r *Request, next func()) {
			if swag.Use(r.Writer, r.Request) {
				r.Exit()
			}
		})
	}
}

// WithLog 设置日志
func WithLog() Option {
	return func(s *Server) {
		s.Use(func(r *Request, next func()) {
			start := time.Now()
			defer func() { frame.Log.Printf("%-7s %s  耗时: %s\n", r.Method, r.URL, time.Now().Sub(start)) }()
			next()
		})
	}
}

func WithPing(content ...any) Option {
	return func(s *Server) {
		s.GET("/ping", func(r *Request) {
			if len(content) > 0 {
				in.Text200(content[0])
			}
			in.Succ(nil)
		})
	}
}

// WithPort 设置端口
func WithPort(port ...int) Option {
	return func(s *Server) { s.SetPort(port...) }
}

// WithPrefix 设置全局前缀,注意使用
func WithPrefix(prefix string) Option {
	return func(s *Server) { s.Grouper.Prefix = prefix }
}

func New(op ...Option) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		Port:    []int{frame.DefaultPort},
		ctx:     ctx,
		cancel:  cancel,
		Grouper: &Grouper{Router: mux.NewRouter()},
	}
	for _, v := range op {
		v(s)
	}
	return s
}

type Server struct {
	Port   []int
	use    []func(r *Request, next func())
	ctx    context.Context
	cancel context.CancelFunc
	*Grouper
}

func (this *Server) Close() error {
	this.cancel()
	return nil
}

// ServeHTTP 实现http.Handler接口
func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := NewRequest(w, r)
	r = r.WithContext(context.WithValue(r.Context(), "request", req))
	this.next(0, req)
}

// next 递归执行中间件,不用忘记调用next函数
func (this *Server) next(i int, r *Request) {
	if i >= len(this.use) {
		this.Router.ServeHTTP(r.Writer, r.Request)
	} else {
		do := false
		next := func() { do = true; this.next(i+1, r) }
		this.use[i](r, next)
		if !do {
			next()
		}
	}
}

// SetPort 设置端口
func (this *Server) SetPort(port ...int) *Server {
	this.Port = port
	return this
}

func (this *Server) Use(f ...func(r *Request, next func())) {
	this.use = append(this.use, f...)
}

func (this *Server) Run() (err error) {

	f := func(port int) (err error) {
		defer func() {
			frame.Log.Println("HTTP服务结束监听:", err)
		}()
		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: in.Recover(this),
		}
		go func() {
			<-this.ctx.Done()
			s.Close()
		}()
		frame.Log.Printf("[%s] 开启HTTP服务成功...\n", s.Addr)
		return s.ListenAndServe()
	}

	if len(this.Port) == 0 {
		this.Port = []int{80}
	}

	wg := sync.WaitGroup{}
	for _, port := range this.Port {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			err = f(port)
		}(port)
	}
	wg.Wait()
	return
}

type Grouper struct {
	Router *mux.Router        //路由实例
	Prefix string             //路由前缀
	middle []func(r *Request) //中间件
}

func (this *Grouper) Middle(f ...func(r *Request)) {
	this.middle = append(this.middle, f...)
}

func (this *Grouper) IgnoreMiddle(handler func(g *Grouper)) *Grouper {
	g := *this
	x := &g
	if handler != nil {
		handler(x)
	}
	return x
}

func (this *Grouper) Group(path string, handler func(g *Grouper)) *Grouper {
	prefix := this.Prefix + path
	g := &Grouper{
		Router: this.Router,
		Prefix: prefix,
		middle: this.middle,
	}
	if handler != nil {
		handler(g)
	}
	return g
}

// Static 放在最后执行,这个会占用Grouper下所有的路由
func (this *Grouper) Static(path string, dir string) *Grouper {
	path = this.Prefix + path
	s := http.StripPrefix(path, http.FileServer(http.Dir(dir)))
	this.Router.PathPrefix(path).Handler(s)
	return this
}

// StaticEmbed 放在最后执行,这个会占用Grouper下所有的路由
func (this *Grouper) StaticEmbed(path string, e embed.FS, dir string) error {
	web, err := fs.Sub(e, dir)
	if err != nil {
		return err
	}
	path = this.Prefix + path
	s := http.StripPrefix(path, http.FileServer(http.FS(web)))
	this.Router.PathPrefix(path).Handler(s)
	return nil
}

// StaticFS 放在最后执行,这个会占用Grouper下所有的路由
func (this *Grouper) StaticFS(path string, fs fs.FS) error {
	path = this.Prefix + path
	s := http.StripPrefix(path, http.FileServer(http.FS(fs)))
	this.Router.PathPrefix(path).Handler(s)
	return nil
}

func (this *Grouper) ALL(path string, handler func(r *Request)) *Grouper {
	return this.do([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodPatch, http.MethodConnect, http.MethodOptions, http.MethodTrace}, path, handler)
}

func (this *Grouper) GET(path string, handler func(r *Request)) *Grouper {
	return this.do([]string{http.MethodGet}, path, handler)
}

func (this *Grouper) POST(path string, handler func(r *Request)) *Grouper {
	return this.do([]string{http.MethodPost}, path, handler)
}

func (this *Grouper) PUT(path string, handler func(r *Request)) *Grouper {
	return this.do([]string{http.MethodPut}, path, handler)
}

func (this *Grouper) DELETE(path string, handler func(r *Request)) *Grouper {
	return this.do([]string{http.MethodDelete}, path, handler)
}

func (this *Grouper) do(method []string, path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middles := this.middle
	this.Router.Methods(method...).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req *Request
		var ok bool

		if v := r.Context().Value("request"); v == nil {
			req = NewRequest(w, r)
		} else if req, ok = v.(*Request); !ok {
			req = NewRequest(w, r)
		}

		for _, f := range middles {
			f(req)
		}

		handler(req)
	})
	return this
}
