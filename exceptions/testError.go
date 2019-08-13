package exceptions

import (
	"time"

	"vid/models"
)

type TestError struct {
	msg    string
	detail string
}

func NewTestError(msg string, detail string) *TestError {
	return &TestError{
		msg:    msg,
		detail: detail,
	}
}

func (e *TestError) Error() string {
	return e.msg
}

func (e *TestError) Info() *models.ErrorInfo {
	return &models.ErrorInfo{
		Message: e.msg,
		Detail:  "TestError: " + e.detail,
		Time:    time.Now().String(),
	}
}
