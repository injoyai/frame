package middle

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var DefaultSwagger = &Swagger{
	IndexPath: "/swagger",
	JsonPath:  "/swagger/swagger.json",
	Filename:  "./docs/swagger.json",
	UI:        DefaultSwaggerUI,
}

type Swagger struct {
	IndexPath string //swagger的路由
	JsonPath  string //json路由
	Filename  string //json文件名称
	UI        string //
}

func (this *Swagger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case this.IndexPath:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(this.UI, this.JsonPath)))
	case this.JsonPath:
		DefaultSwagger.WriteFile(w)
	}
}

func (this *Swagger) Use(w http.ResponseWriter, r *http.Request) bool {
	switch r.URL.Path {
	case this.IndexPath:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(this.UI, this.JsonPath)))
		return true
	case this.JsonPath:
		DefaultSwagger.WriteFile(w)
		return true
	}
	return false
}

func (this *Swagger) WriteFile(w http.ResponseWriter) {
	f, err := os.Open(this.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

var (
	DefaultSwaggerUI = `<!DOCTYPE html>
        <html>
          <head>
            <title>SwaggerUI</title>
            <!-- needed for adaptive design -->
            <meta charset="utf-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
            <style>
              body {
                margin: 0;
                padding: 0;
              }
            </style>
          </head>
          <body>
            <redoc spec-url='%s'></redoc>
            <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
          </body>
        </html>`
)
