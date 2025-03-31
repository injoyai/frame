package fiber

import (
	"context"
	"errors"
	"github.com/fasthttp/websocket"
	"github.com/injoyai/conv"
	"github.com/valyala/fasthttp"
)

const (
	WSText   = websocket.TextMessage
	WSBinary = websocket.BinaryMessage
	WSClose  = websocket.CloseMessage
	WSPing   = websocket.PingMessage
	WSPong   = websocket.PongMessage
)

var DefaultUpgrader = websocket.FastHTTPUpgrader{
	HandshakeTimeout:  0,
	Subprotocols:      nil,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
	WriteBufferPool:   nil,
	CheckOrigin:       func(ctx *fasthttp.RequestCtx) bool { return true },
}

type Websocket struct {
	*websocket.Conn
	ctx context.Context
}

// ReadMessage 实现ios.MReader接口
func (this *Websocket) ReadMessage() ([]byte, error) {
	_, p, err := this.Conn.ReadMessage()
	return p, err
}

// DiscardRead 丢弃读取数据,但是还是需要读取,才能监听到客户端关闭信号
// 例如浏览器,需要监听才能响应给客户端,浏览器才能正常关闭ws链接
func (this *Websocket) DiscardRead() {
	go func() {
		for {
			select {
			case <-this.ctx.Done():
				return
			default:
				_, _, err := this.Conn.ReadMessage()
				if err != nil {
					return
				}
			}
		}
	}()
}

// Write 实现io.Writer接口
func (this *Websocket) Write(p []byte) (int, error) {
	err := this.WriteMessage(WSBinary, p)
	return len(p), err
}

func (this *Websocket) WriteText(s string) (int, error) {
	err := this.WriteMessage(WSText, []byte(s))
	return len(s), err
}

// WriteChan 从chan中读取,并写入到ws
func (this *Websocket) WriteChan(c chan any, messageType ...int) error {
	mt := WSText
	if len(messageType) > 0 {
		mt = messageType[0]
	}
	for {
		select {
		case <-this.ctx.Done():
			return this.ctx.Err()
		case data, ok := <-c:
			if !ok {
				return errors.New("chan closed")
			}
			if err := this.WriteMessage(mt, conv.Bytes(data)); err != nil {
				return err
			}
		}
	}
}
