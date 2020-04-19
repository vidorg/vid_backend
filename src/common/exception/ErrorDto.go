package exception

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xstack"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"strings"
	"time"
)

type ErrorDto struct {
	Time    string   `json:"time"`
	Type    string   `json:"type"`
	Detail  string   `json:"detail"`
	Request []string `json:"request"`

	Filename string `json:"filename,omitempty"`
	Function string `json:"function,omitempty"`
	Line     int    `json:"line,omitempty"`
	Content  string `json:"content,omitempty"`
}

type Stack struct {
	Index    int
	Filename string
	Function string
	Pc       uintptr
	Line     int
	Content  string
}

func (s *Stack) String() string {
	return fmt.Sprintf("%s:%d (0x%x)\n\t%s: %s", s.Filename, s.Line, s.Pc, s.Function, s.Content)
}

func NewErrorDto(err interface{}, skip int, c *gin.Context, print bool) *ErrorDto {
	now := time.Now().Format(time.RFC3339)
	errType := fmt.Sprintf("%T", err)
	errDetail := fmt.Sprintf("%v", err)
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	}
	dto := &ErrorDto{
		Time:   now,
		Type:   errType,
		Detail: errDetail,
	}

	// request
	if c != nil {
		requestBytes, _ := httputil.DumpRequest(c.Request, false)
		requestParams := strings.Split(string(requestBytes), "\r\n")
		request := make([]string, 0)
		for _, param := range requestParams {
			if strings.HasPrefix(param, "Authorization:") { // Authorization header
				request = append(request, "Authorization: *")
			} else if param != "" { // other param
				request = append(request, param)
			}
		}
		dto.Request = request
	}

	// runtime
	if skip >= 0 {
		stacks := xstack.GetStack(skip)
		filename := stacks[0].Filename
		function := stacks[0].Function
		line := stacks[0].Line
		content := stacks[0].Content
		if print {
			fmt.Println()
			fmt.Println(xcolor.Yellow.Paint("[Panic Stack]"))
			xstack.PrintStacks(stacks)
			fmt.Println(xcolor.Yellow.Paint("[Panic Stack End]"))
			fmt.Println()
		}

		dto.Filename = filename
		dto.Function = function
		dto.Line = line
		dto.Content = content
	}

	return dto
}
