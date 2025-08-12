package fiber

import (
	"bufio"
	"github.com/injoyai/conv"
	"io"
	"unsafe"
)

type SSE interface {
	io.Writer
	io.StringWriter
	WriteAny(v any) error
}

type sse struct {
	*bufio.Writer
}

func (this *sse) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := this.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	_, err = this.Writer.Write([]byte("\n\n"))
	if err != nil {
		return 0, err
	}
	return n, this.Writer.Flush()
}

func (this *sse) WriteString(s string) (int, error) {
	return this.Write(*(*[]byte)(unsafe.Pointer(&s)))
}

func (this *sse) WriteAny(v any) error {
	if bs, ok := v.([]byte); ok {
		_, err := this.Writer.Write(bs)
		return err
	}
	s := conv.String(v)
	_, err := this.Writer.Write(*(*[]byte)(unsafe.Pointer(&s)))
	return err
}
