package fiber

import (
	"bufio"
	"github.com/injoyai/conv"
	"io"
	"strings"
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

	return this.Write([]byte(s))
}

func (this *sse) WriteStrings(s ...string) (int, error) {
	return this.Write([]byte(strings.Join(s, "\n")))
}

func (this *sse) WriteAny(v any) error {
	s := conv.String(v)
	_, err := this.Writer.Write(*(*[]byte)(unsafe.Pointer(&s)))
	return err
}
