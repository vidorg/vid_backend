package exception

type ServerError struct {
	Code    int
	Message string
}

func NewError(code int, message string) *ServerError {
	return &ServerError{
		Code:    code,
		Message: message,
	}
}
