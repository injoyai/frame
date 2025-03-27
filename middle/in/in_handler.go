package in

import (
	"github.com/injoyai/conv"
	"net/http"
)

// NewSuccWithCode 有些code为0是成功,有些ok是成功...
func (this *Client) NewSuccWithCode(code interface{}) func(data interface{}, count ...int64) {
	return func(data interface{}, count ...int64) {
		if len(count) > 0 {
			this.Json(http.StatusOK, &ResponseCount{
				Code:    code,
				Data:    data,
				Message: "成功",
				Count:   count[0],
			})
			return
		}
		this.Json(http.StatusOK, &Response{
			Code:    code,
			Data:    data,
			Message: "成功",
		})
	}
}

func (this *Client) NewFailWithCode(code interface{}) func(msg interface{}) {
	return func(msg interface{}) {
		this.Json(http.StatusOK, &Response{
			Code:    code,
			Message: conv.String(msg),
		})
	}
}

func (this *Client) NewUnauthorizedWithCode(code interface{}) func() {
	return func() {
		this.Json(http.StatusOK, &Response{
			Code:    code,
			Message: "验证失败",
		})
	}
}

func (this *Client) NewForbiddenWithCode(code interface{}) func() {
	return func() {
		this.Json(http.StatusOK, &Response{
			Code:    code,
			Message: "没有权限",
		})
	}
}
