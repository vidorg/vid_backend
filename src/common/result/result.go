package result

import (
	"github.com/Aoi-hosizora/ahlib-web/xdto"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"net/http"
	"strings"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("Result", "global response").
			Properties(
				goapidoc.NewProperty("code", "integer#int32", true, "status code"),
				goapidoc.NewProperty("message", "string", true, "status message"),
			),

		goapidoc.NewDefinition("_Result", "global response").
			Generics("T").
			Properties(
				goapidoc.NewProperty("code", "integer#int32", true, "status code"),
				goapidoc.NewProperty("message", "string", true, "status message"),
				goapidoc.NewProperty("data", "T", true, "response data"),
			),
	)
}

type Result struct {
	Status  int32          `json:"-"`
	Code    int32          `json:"code"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *xdto.ErrorDto `json:"error,omitempty"`
}

func Status(status int32) *Result {
	message := http.StatusText(int(status))
	if status == 200 || status == 201 {
		message = "success"
	} else if message == "" {
		message = "unknown"
	}
	return &Result{
		Status:  status,
		Code:    status,
		Message: strings.ToLower(message),
	}
}

func Ok() *Result {
	return Status(http.StatusOK)
}

func Created() *Result {
	return Status(http.StatusCreated)
}

func Error(e *exception.Error) *Result {
	return Status(e.Status).SetCode(e.Code).SetMessage(e.Message)
}

func (r *Result) SetStatus(status int32) *Result {
	r.Status = status
	return r
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
	r.Data = data
	return r
}

func (r *Result) SetPage(page int32, limit int32, total int32, data interface{}) *Result {
	r.Data = NewPage(page, limit, total, data)
	return r
}

func (r *Result) SetError(err error, c *gin.Context) *Result {
	if gin.Mode() == gin.DebugMode && err != nil {
		r.Error = xgin.BuildBasicErrorDto(err, c, map[string]interface{}{
			"request_ip": c.ClientIP(),
			"request_id": c.Writer.Header().Get("X-Request-Id"),
		})
	}
	return r
}

func (r *Result) JSON(c *gin.Context) {
	c.JSON(int(r.Status), r)
}

func (r *Result) XML(c *gin.Context) {
	c.XML(int(r.Status), r)
}

// Simplify controller's functions.
func J(fn func(c *gin.Context) *Result) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		result := fn(c)
		if result != nil {
			result.JSON(c)
		}
	}
}
