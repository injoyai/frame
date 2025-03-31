package in

import (
	"github.com/injoyai/base/maps"
	"io"
	"net/http"
)

type Client interface {
	Recover(h http.Handler) http.Handler
	RecoverFunc(f http.HandlerFunc) http.Handler
	RecoverWriter(e any, w http.ResponseWriter)
	Bind(path string, handler http.HandlerFunc)
	SetWriterOption(f ...WriterOption)
	SetHandlerWithCode(succ, fail, unauthorized, forbidden any)
	Proxy(w http.ResponseWriter, r *http.Request, uri string)
	Cache() *maps.Safe
	NewSuccWithCode(code any) func(data any, count ...int64)
	NewFailWithCode(code any) func(msg any)
	NewUnauthorizedWithCode(code any) func()
	NewForbiddenWithCode(code any) func()

	Respondent
}

type Respondent interface {
	Custom(code int, data any, header http.Header)
	Custom200(data any, h http.Header)

	Text(code int, data any)
	Text200(data any)

	Html(code int, data any)
	Html200(data any)
	HtmlLocal(filename string)
	HtmlReader(r io.ReadCloser)

	Json(code int, data any)
	Json200(data any)
	Json400(data any)
	Json401(data any)
	Json403(data any)
	Json415(data any)
	Json500(data any)

	Reader(code int, r io.ReadCloser)
	Reader200(r io.ReadCloser)

	File(name string, size int64, r io.ReadCloser)
	FileLocal(name, filename string)
	FileReader(name string, r io.ReadCloser)
	FileBytes(name string, bs []byte)

	Redirect301(addr string)
	Redirect302(addr string)
	Redirect307(addr string)
	Redirect308(addr string)

	Succ(data any, count ...int64)
	Fail(data any)
	Unauthorized()
	Forbidden()
	Err(data any)
	CheckErr(err error, failMsg ...string)
	NewWriter(code int, i IMarshal) Writer
	Exit()
}
