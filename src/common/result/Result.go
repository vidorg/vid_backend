package result

import (
	"github.com/vidorg/vid_backend/src/common/exception"
	"net/http"
	"strings"

	"github.com/Aoi-hosizora/ahlib/xlinkedhashmap"
	"github.com/gin-gonic/gin"
)

// @Model         Result
// @Description   返回统一响应结果
// @Property      code    integer true "响应码"
// @Property      message string  true "状态信息"
type Result struct {
	Code    int32                         `json:"code"`
	Message string                        `json:"message"`
	Data    *xlinkedhashmap.LinkedHashMap `json:"data,omitempty"`
}

func Status(code int32) *Result {
	message := http.StatusText(int(code))
	if code == 200 {
		message = "success"
	} else if message == "" {
		message = "unknown"
	}
	return &Result{
		Code:    code,
		Message: strings.ToLower(message),
	}
}

func Ok() *Result {
	return Status(http.StatusOK)
}

func Error(se *exception.ServerError) *Result {
	return Status(se.Code).SetMessage(se.Message)
}

func (r *Result) SetCode(code int32) *Result {
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

func (r *Result) SetPage(total int32, page int32, limit int32, data interface{}) *Result {
	if r.Data == nil {
		r.Data = xlinkedhashmap.NewLinkedHashMap()
	}
	r.Data.Set("total", total)
	r.Data.Set("page", page)
	r.Data.Set("limit", limit)
	r.Data.Set("data", data)
	return r
}

func (r *Result) JSON(c *gin.Context) {
	c.JSON(int(r.Code), r)
}

func (r *Result) XML(c *gin.Context) {
	c.XML(int(r.Code), r)
}
