package result

import (
	"net/http"

	"github.com/Aoi-hosizora/ahlib/xhashmap"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    *xhashmap.LinkedHashMap `json:"data,omitempty"`
}

func (Result) Result(code int) *Result {
	var message string
	switch code {
	case http.StatusOK: // 200
		message = "success"
	case http.StatusCreated: // 201
		message = "created"
	case http.StatusBadRequest: // 400
		message = "bad Request"
	case http.StatusUnauthorized: // 401
		message = "unauthorized"
	case http.StatusForbidden: // 403
		message = "forbidden"
	case http.StatusNotFound: // 404
		message = "not found"
	case http.StatusMethodNotAllowed: // 405
		message = "method not allowed"
	case http.StatusNotAcceptable: // 406
		message = "not acceptable"
	case http.StatusRequestEntityTooLarge: // 413
		message = "request entity too large"
	case http.StatusUnsupportedMediaType: // 415
		message = "unsupported media type"
	case http.StatusInternalServerError: // 500
		message = "internal server error"
	default:
		message = "unknown status"
	}

	return &Result{
		Code:    code,
		Message: message,
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
	r.Message = message
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
