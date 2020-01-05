package exception

import (
	"gopkg.in/go-playground/validator.v9"
)

func WrapValidationError(err error) error {
	// validator.ValidationErrors
	// *errors.errorString
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return RequestParamError
	}

	for _, field := range errs {
		if field.Tag() == "required" {
			return RequestParamError // weight
		}
	}
	return RequestFormatError
}
