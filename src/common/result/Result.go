package result

import (
	"net/http"
	"strings"

	"github.com/Aoi-hosizora/ahlib/xhashmap"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    *xhashmap.LinkedHashMap `json:"data,omitempty"`
}

func (Result) Result(code int) *Result {
	message := http.StatusText(code)
	if message == "" {
		message = "Unknown status"
	}
	return &Result{
		Code:    code,
		Message: strings.ToLower(message),
	}
}

func (Result) Ok() *Result {
	return Result{}.Result(http.StatusOK)
}

func (Result) Error() *Result {
	return Result{}.Result(http.StatusInternalServerError)
}

func (r *Result) SetCode(code int) *Result {
	r.Code = code
	return r
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = strings.ToLower(message)
	return r
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = xhashmap.ObjectToLinkedHashMap(data)
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = new(xhashmap.LinkedHashMap)
	}
	r.Data.Set(field, data)
	return r
}

func (r *Result) SetPage(count int32, page int32, data interface{}) *Result {
	if r.Data == nil {
		r.Data = new(xhashmap.LinkedHashMap)
	}
	r.Data.Set("total", count)
	r.Data.Set("page", page)
	r.Data.Set("data", data)
	return r
}

func (r *Result) JSON(c *gin.Context) {
	c.JSON(r.Code, r)
}

func (r *Result) XML(c *gin.Context) {
	c.XML(r.Code, r)
}
