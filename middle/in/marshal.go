package in

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/injoyai/conv"
)

type Marshal func(v any) ([]byte, error)

type CUSTOM struct {
	Data    any
	reader  io.Reader
	Headers http.Header
}

func (this *CUSTOM) Read(p []byte) (int, error) {
	if this.reader == nil {
		switch r := this.Data.(type) {
		case io.Reader:
			this.reader = r
		case []byte:
			this.reader = bytes.NewReader(r)
		default:
			this.reader = bytes.NewReader([]byte(conv.String(this.Data)))
		}
	}
	return this.reader.Read(p)
}

func (this *CUSTOM) Header() http.Header {
	return this.Headers
}

func (this *CUSTOM) Close() error {
	if this.reader != nil {
		if v, ok := this.reader.(io.Closer); ok {
			return v.Close()
		}
	}
	return nil
}

type IMarshal interface {
	io.ReadCloser
	Header() http.Header
}

type TEXT struct {
	Data   any
	reader io.Reader
}

func (this *TEXT) Read(p []byte) (int, error) {
	if this.reader == nil {
		switch r := this.Data.(type) {
		case io.Reader:
			this.reader = r
		case []byte:
			this.reader = bytes.NewReader(r)
		default:
			this.reader = bytes.NewReader([]byte(conv.String(this.Data)))
		}
	}
	return this.reader.Read(p)
}

func (this *TEXT) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	}
}

func (this *TEXT) Close() error {
	if this.reader != nil {
		if v, ok := this.reader.(io.Closer); ok {
			return v.Close()
		}
	}
	return nil
}

type JSON struct {
	Data   any
	reader io.Reader
}

func (this *JSON) Read(p []byte) (int, error) {
	if this.reader == nil {
		switch r := this.Data.(type) {
		case io.Reader:
			this.reader = r
		case []byte:
			this.reader = bytes.NewReader(r)
		default:
			bs, err := json.Marshal(this.Data)
			if err != nil {
				return 0, err
			}
			this.reader = bytes.NewReader(bs)
		}
	}
	return this.reader.Read(p)
}

func (this *JSON) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json; charset=utf-8"},
	}
}

func (this *JSON) Close() error {
	if this.reader != nil {
		if v, ok := this.reader.(io.Closer); ok {
			return v.Close()
		}
	}
	return nil
}

type HTML struct {
	Data   any
	reader io.Reader
}

func (this *HTML) Read(p []byte) (int, error) {
	if this.reader == nil {
		switch r := this.Data.(type) {
		case io.Reader:
			this.reader = r
		case []byte:
			this.reader = bytes.NewReader(r)
		default:
			this.reader = bytes.NewReader([]byte(conv.String(this.Data)))
		}
	}
	return this.reader.Read(p)
}

func (this *HTML) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"text/html; charset=utf-8"},
	}
}

func (this *HTML) Close() error {
	if this.reader != nil {
		if v, ok := this.reader.(io.Closer); ok {
			return v.Close()
		}
	}
	return nil
}

type FILE struct {
	Name string
	Size int64
	io.ReadCloser
}

func (this *FILE) Read(p []byte) (int, error) {
	if this.ReadCloser == nil {
		return 0, io.EOF
	}
	return this.ReadCloser.Read(p)
}

func (this *FILE) Header() http.Header {
	return http.Header{
		"Content-Type":        []string{"application/octet-stream"},
		"Content-Disposition": []string{"attachment; filename=" + this.Name},
		"Content-Length":      []string{strconv.FormatInt(this.Size, 10)},
	}
}

type READER struct {
	io.ReadCloser
}

func (this *READER) Read(p []byte) (int, error) {
	if this.ReadCloser == nil {
		return 0, io.EOF
	}
	return this.ReadCloser.Read(p)
}

func (this *READER) Header() http.Header {
	return http.Header{}
}

func (this *READER) Close() error {
	if this.ReadCloser != nil {
		return this.ReadCloser.Close()
	}
	return nil
}
