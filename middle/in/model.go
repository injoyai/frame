package in

type ResponseCount struct {
	Code    any    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"msg"`
	Count   int64  `json:"count"`
}

type Response struct {
	Code    any    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"msg"`
}
