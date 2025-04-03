package fiber

import (
	"bufio"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/injoyai/conv"
	"github.com/injoyai/frame/middle/in"
	"sync"
	"time"
)

type (
	Bind         = fiber.Bind
	Writer       = in.Writer
	WriterOption = in.WriterOption
	Respondent   = in.Respondent
)

var (
	_nil = conv.Nil()
)

var ctxPoll = sync.Pool{New: func() any {
	return &ctx{}
}}

type Ctx interface {
	fiber.Ctx
	Requester
	in.Respondent

	SetHeader(k, v string)
	SetContentType(contentType string)

	// Parse 解析body数据到ptr,需要指针类型
	Parse(ptr any)

	// Websocket websocket
	Websocket(handler func(conn *Websocket))

	SSE(handler func(w SSE))

	// Stream stream
	Stream(handler func(w *bufio.Writer))

	// free 释放内存
	free()
}

func NewCtx(c fiber.Ctx, r Respondent) Ctx {
	cc := ctxPoll.Get().(*ctx)
	cc.Ctx = c
	cc.Extend = conv.NewExtend(cc)
	cc.Respondent = r
	return cc
}

type ctx struct {
	fiber.Ctx
	conv.Extend
	Respondent
	requestHeader map[string][]string
	bodyMap       *conv.Map
}

func (this *ctx) Get(key string, defaultValue ...string) string {
	return this.Ctx.Get(key, defaultValue...)
}

func (this *ctx) GetVar(key string) *conv.Var {
	//尝试从query中获取数据
	if val := this.Ctx.RequestCtx().QueryArgs().Peek(key); val != nil {
		return conv.New(val)
	}

	//尝试从path中获取数据
	if val := this.Ctx.Params(key); val != "" {
		return conv.New(val)
	}

	//尝试从header中获取数据
	if this.requestHeader == nil {
		this.requestHeader = this.Ctx.GetReqHeaders()
	}
	if vals, ok := this.requestHeader[key]; ok && len(vals) > 0 {
		return conv.New(vals[0])
	}

	//尝试从body中获取数据
	contentType := this.requestHeader[fiber.HeaderContentType]
	if len(contentType) == 0 {
		contentType = []string{fiber.MIMEApplicationJSON}
	}
	switch contentType[0] {
	case fiber.MIMEApplicationJSON:
		if this.bodyMap == nil {
			this.bodyMap = conv.NewMap(this.Ctx.Body())
		}
		if val := this.bodyMap.GetVar(key); !val.IsNil() {
			return val
		}
	case fiber.MIMEApplicationForm:
		if val := this.Ctx.RequestCtx().FormValue(key); val != nil {
			return conv.New(val)
		}
	}

	return _nil
}

func (this *ctx) SetHeader(k, v string) {
	this.Ctx.Set(k, v)
}

func (this *ctx) SetContentType(contentType string) {
	this.Ctx.Type(fiber.HeaderContentType, contentType)
}

func (this *ctx) Parse(ptr any) {
	err := conv.Unmarshal(this.Ctx.Body(), ptr)
	this.CheckErr(err)
}

func (this *ctx) Websocket(handler func(ws *Websocket)) {
	ctx := this.Ctx.Context()
	err := DefaultUpgrader.Upgrade(this.Ctx.RequestCtx(), func(ws *websocket.Conn) {
		defer ws.Close()
		handler(&Websocket{
			Conn: ws,
			ctx:  ctx,
		})
	})
	this.CheckErr(err)
}

func (this *ctx) Stream(handler func(w *bufio.Writer)) {
	this.SetHeader(fiber.HeaderContentType, "text/event-stream")
	this.SetHeader(fiber.HeaderCacheControl, "no-cache")
	this.SetHeader(fiber.HeaderConnection, "keep-alive")
	err := this.Ctx.SendStreamWriter(handler)
	this.CheckErr(err)
}

func (this *ctx) SSE(handler func(sse SSE)) {
	this.Stream(func(writer *bufio.Writer) { handler(&sse{Writer: writer}) })
}

// free 手动释放内存
func (this *ctx) free() {
	this.Ctx = nil
	this.Respondent = nil
	this.Extend = nil
	this.requestHeader = nil
	this.bodyMap = nil
	ctxPoll.Put(this)
}

func (this *ctx) WriteHeader(statusCode int) {
	this.Ctx.Status(statusCode)
}

/*



 */

type Requester interface {
	IsDefault(key string) bool
	IsNumber(key string) bool
	IsString(key string) bool
	IsInt(key string) bool
	IsFloat(key string) bool
	IsBool(key string) bool
	IsTime(key string) bool
	IsArray(key string) bool
	IsPointer(key string) bool
	IsNil(key string) bool
	GetInterface(key string, def ...interface{}) interface{}
	GetInterfaces(key string, def ...[]interface{}) []interface{}
	GetByte(key string, def ...byte) byte
	GetBytes(key string, def ...[]byte) []byte
	GetRune(key string, def ...rune) rune
	GetRunes(key string, def ...[]rune) []rune
	GetString(key string, def ...string) string
	GetStrings(key string, def ...[]string) []string
	GetInt(key string, def ...int) int
	GetInts(key string, def ...[]int) []int
	GetInt8(key string, def ...int8) int8
	GetInt16(key string, def ...int16) int16
	GetInt32(key string, def ...int32) int32
	GetInt64(key string, def ...int64) int64
	GetInt64s(key string, def ...[]int64) []int64
	GetUint(key string, def ...uint) uint
	GetUint8(key string, def ...uint8) uint8
	GetUint16(key string, def ...uint16) uint16
	GetUint32(key string, def ...uint32) uint32
	GetUint64(key string, def ...uint64) uint64
	GetFloat32(key string, def ...float32) float32
	GetFloat64(key string, def ...float64) float64
	GetBool(key string, def ...bool) bool
	GetDuration(key string, def ...time.Duration) time.Duration
	GetMicrosecond(key string, def ...time.Duration) time.Duration
	GetMillisecond(key string, def ...time.Duration) time.Duration
	GetSecond(key string, def ...time.Duration) time.Duration
	GetMinute(key string, def ...time.Duration) time.Duration
	GetHour(key string, def ...time.Duration) time.Duration
	GetErr(key string, def ...error) error
	GetGMap(key string, def ...map[string]interface{}) map[string]interface{}
	GetSMap(key string, def ...map[string]string) map[string]string
	GetIMap(key string, def ...map[interface{}]interface{}) map[interface{}]interface{}
	GetDMap(key string, def ...interface{}) *conv.Map
}
