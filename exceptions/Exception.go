package exceptions

import "errors"

// Authorization
var (
	AuthorizationException = errors.New("Authorization failed")
	TokenExpiredException  = errors.New("Token has expired")
	TokenInvalidException  = errors.New("Token invalid")
	PasswordException      = errors.New("User password error")
)

// Db
var (
	// exist
	UserExistException     = errors.New("User already existed")
	UserNotExistException  = errors.New("User not found")
	VideoNotExistException = errors.New("Video not found")

	// user action failed
	ModifyPassException = errors.New("User password modify failed")
	InsertUserException = errors.New("User insert failed")
	UpdateUserException = errors.New("User update failed")
	DeleteUserException = errors.New("User delete failed")

	// user other exception
	UserNameUsedException     = errors.New("Username has been used")
	NotUpdateUserException    = errors.New("User information not updated")
	UpdateInvalidException    = errors.New("User information invalid")
	SubscribeOneSelfException = errors.New("Cound not subscribe to oneself")

	// video exception
)

// Ctrl
var (
	// request error
	RequestBodyError = errors.New("Request body error")
	QueryParamError  = errors.New("Query param '%s' not found or error")
	RouteParamError  = errors.New("Route param '%s' not found or error")

	// format error
	LoginFormatError    = errors.New("Login username or password format error")
	RegisterFormatError = errors.New("Register username or password format error")
)
