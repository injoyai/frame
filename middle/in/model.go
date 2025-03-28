package in

type ResponseCount struct {
	Code    interface{} `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
	Count   int64       `json:"count"`
}

type Response struct {
	Code    interface{} `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
}
