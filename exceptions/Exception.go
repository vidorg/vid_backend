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
	UserExistException     = errors.New("User already existed")
	UserNotExistException  = errors.New("User not found")
	VideoNotExistException = errors.New("Video not found")

	InsertException = errors.New("User insert failed")
	UpdateException = errors.New("User update failed")
	DeleteException = errors.New("User delete failed")

	UserNameUsedException     = errors.New("Username has been used")
	NotUpdateException        = errors.New("User information not updated")
	UpdateInvalidException    = errors.New("User information invalid")
	SubscribeOneSelfException = errors.New("Cound not subscribe to oneself")
)

// Ctrl
var (
	RequestBodyError = errors.New("Request body error")
	QueryParamError  = errors.New("Query param '%s' not found or error")
	RouteParamError  = errors.New("Route param '%s' not found or error")

	LoginFormatError    = errors.New("Login username or password format error")
	RegisterFormatError = errors.New("Register username or password format error")

	NoAuthorizationException = errors.New("No authorization of this user")
	PasswordError            = errors.New("User password error")
)
