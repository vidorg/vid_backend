package dto

import (
	"net/http"
	"vid/app/model/vo"
)

type Result struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *vo.OrderMap `json:"data,omitempty"` // map[string]interface{}
}

func (Result) Ok() *Result {
	return &Result{
		Code:    200,
		Message: "Success",
	}
}

func (Result) Error(code int) *Result {
	var message string
	switch code {
	case http.StatusBadRequest:
		message = "Bad Request"
	case http.StatusUnauthorized:
		message = "Unauthorized"
	case http.StatusNotFound:
		message = "Not Found"
	case http.StatusInternalServerError:
		message = "Internal Server Error"
	default:
		message = "Unknown Error"
	}

	return &Result{
		Code:    code,
		Message: message,
	}
}

func (r *Result) SetCode(code int) *Result {
	r.Code = code
	return r
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = vo.OrderMap{}.FromObject(data)
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = vo.NewOrderMap()
	}
	r.Data.Put(field, data)
	return r
}

func (r *Result) SetPage(count int, page int, data interface{}) *Result {
	if r.Data == nil {
		r.Data = vo.NewOrderMap()
	}
	r.Data.Put("count", count)
	r.Data.Put("page", page)
	r.Data.Put("data", data)
	return r
}
