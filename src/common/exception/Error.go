package exception

type Error struct {
	Code    int32
	Message string
}

func NewError(code int32, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
