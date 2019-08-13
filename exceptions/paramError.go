package exceptions

import (
	"fmt"
	"time"

	"vid/models"
)

type ParamError struct {
	params  []string
	isQuery bool
}

func NewParamError(params []string, isQuery bool) *ParamError {
	return &ParamError{
		params:  params,
		isQuery: isQuery,
	}
}

func (e *ParamError) Info() *models.ErrorInfo {
	var detail string
	if e.isQuery {
		detail = "Query parameters"
	} else {
		detail = "Request body json"
	}
	return &models.ErrorInfo{
		Message: "Parameter Error",
		Detail:  fmt.Sprintf("%s %s not found or error", detail, e.params),
		Time:    time.Now().String(),
	}
}
