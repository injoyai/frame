package swagger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

var Default = &Swagger{
	IndexPath:    "/swagger",
	JsonPath:     "/swagger/swagger.json",
	JsonFilename: "./docs/swagger.json",
	UI:           DefaultUI,
}

type Swagger struct {
	IndexPath    string //swagger的路由
	JsonPath     string //json路由
	JsonFilename string //json文件名称
	JsonBytes    []byte //文件流
	UI           string //ui界面
}

func (this *Swagger) Do(path string, f func(r io.Reader, contentType string)) (bool, error) {
	switch path {
	case this.IndexPath:
		r := strings.NewReader(fmt.Sprintf(this.UI, this.JsonPath))
		f(r, "text/html")

	case this.JsonPath:
		if this.JsonBytes == nil {
			file, err := os.Open(this.JsonFilename)
			if err != nil {
				return false, err
			}
			defer file.Close()
			if this.JsonBytes, err = io.ReadAll(file); err != nil {
				return false, err
			}
		}
		f(bytes.NewReader(this.JsonBytes), "application/json")

	default:
		return false, nil

	}

	return true, nil
}

var (
	DefaultUI = `<!DOCTYPE html>
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
