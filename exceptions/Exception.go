package exceptions

import "errors"

// Jwt
var (
	AuthorizationException = errors.New("Authorization failed")
	TokenExpiredException  = errors.New("Token has expired")
	TokenInvalidException  = errors.New("Token invalid")
)

// Db
var (
	UserExistException    = errors.New("User already existed")
	UserNotExistException = errors.New("User not found")

	InsertException = errors.New("User insert failed")
	UpdateException = errors.New("User update failed")
	DeleteException = errors.New("User delete failed")

	NotUpdateException        = errors.New("User information not updated")
	SubscribeOneSelfException = errors.New("Cound not subscribe to oneself")
)

// Ctrl
var (
	RequestBodyError = errors.New("Request body error")
	QueryParamError  = errors.New("Query param '%s' not found or error")
	RouteParamError  = errors.New("Route param '%s' not found or error")

	LoginFormatError    = errors.New("Login username or password format error")
	RegisterFormatError = errors.New("Register username or password format error")

	PasswordError = errors.New("User password error")
)
