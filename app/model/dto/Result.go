package dto

import (
	"net/http"
	"vid/app/model/vo"
)

type Result struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Data      *vo.OrderMap `json:"data,omitempty"` // map[string]interface{}
	Converter []Converter  `json:"-"`
}

func (r *Result) AddConverter(converter Converter) *Result {
	r.Converter = append(r.Converter, converter)
	return r
}

func (Result) Ok() *Result {
	return &Result{
		Code:    200,
		Message: "success",
	}
}

func (Result) Error(code int) *Result {
	var message string
	switch code {
	case http.StatusBadRequest:
		message = "bad Request"
	case http.StatusUnauthorized:
		message = "unauthorized"
	case http.StatusNotFound:
		message = "not found"
	case http.StatusInternalServerError:
		message = "internal server error"
	default:
		message = "unknown error"
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
	r.convert(data)
	r.Data = vo.OrderMap{}.FromObject(data)
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = vo.NewOrderMap()
	}
	r.convert(data)
	r.Data.Put(field, data)
	return r
}

func (r *Result) SetPage(count int, page int, data interface{}) *Result {
	if r.Data == nil {
		r.Data = vo.NewOrderMap()
	}
	r.convert(data)
	r.Data.Put("count", count)
	r.Data.Put("page", page)
	r.Data.Put("data", data)
	return r
}
