package common

import "github.com/goccy/go-json"

type HTTPResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data"`
}

func NewSuccessHTTPResponse(data any) *HTTPResponse {
	return &HTTPResponse{Success: true, Data: data}
}

func NewErrorHTTPResponse(message string) *HTTPResponse {
	return &HTTPResponse{Success: false, Message: message}
}

func NewErrorCodeHTTPResponse(code int) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: Code2Message[code]}
}

func NewErrorWithCodeHTTPResponse(code int, message string) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: message}
}

func NewErrorWithDataHTTPResponse(code int, message string, data any) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: message, Data: data}
}

func (h HTTPResponse) ToJSON() string {
	v, _ := json.Marshal(h)
	return string(v)
}
