package result

import (
	"net/http"
	"strings"

	"github.com/Aoi-hosizora/ahlib/xlinkedhashmap"
	"github.com/gin-gonic/gin"
)

// @Model         Result
// @Description   返回统一响应结果
// @Property      code    integer true false "响应码"
// @Property      message string  true false "状态信息"
type Result struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    *xlinkedhashmap.LinkedHashMap `json:"data,omitempty"`
}

func Status(code int) *Result {
	message := http.StatusText(code)
	if message == "" {
		message = "Unknown status"
	}
	return &Result{
		Code:    code,
		Message: strings.ToLower(message),
	}
}

func Ok() *Result {
	return Status(http.StatusOK)
}

func Error() *Result {
	return Status(http.StatusInternalServerError)
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
	r.Data = xlinkedhashmap.ObjectToLinkedHashMap(data)
	return r
}

func (r *Result) PutData(field string, data interface{}) *Result {
	if r.Data == nil {
		r.Data = xlinkedhashmap.NewLinkedHashMap()
	}
	r.Data.Set(field, data)
	return r
}

func (r *Result) SetPage(count int32, page int32, data interface{}) *Result {
	if r.Data == nil {
		r.Data = xlinkedhashmap.NewLinkedHashMap()
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
