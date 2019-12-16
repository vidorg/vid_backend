package dto

import (
	"net/http"
	"reflect"
)

type Result struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
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

func (r *Result) SetData(data map[string]interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) SetObject(data interface{}) *Result {
	dataMap := make(map[string]interface{})

	elem := reflect.ValueOf(data).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		dataMap[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	r.Data = dataMap
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = make(map[string]interface{})
	}
	r.Data[field] = data
	return r
}
