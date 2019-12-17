package dto

import (
	lhm "github.com/emirpasic/gods/maps/linkedhashmap"
	"net/http"
	"reflect"
	"strings"
)

type Result struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *lhm.Map `json:"data,omitempty"` // map[string]interface{}
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
	r.Data = lhm.New()

	elem := reflect.ValueOf(data).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		tag := relType.Field(i).Tag.Get("json")
		omitempty := strings.Index(tag, "omitempty") != -1

		field := strings.Split(tag, ",")[0] // elem.Field(i).Name
		value := elem.Field(i).Interface()

		if !omitempty || value != nil {
			r.Data.Put(field, value)
		}
	}
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = lhm.New()
	}
	r.Data.Put(field, data)
	return r
}

func (r *Result) SetPage(count int, page int, data interface{}) *Result {
	if r.Data == nil {
		r.Data = lhm.New()
	}
	r.Data.Put("count", count)
	r.Data.Put("page", page)
	r.Data.Put("data", data)
	return r
}
