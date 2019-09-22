package exceptions

import (
	"errors"
)

// Authorization
var (
	AuthorizationException = errors.New("Authorization failed")
	TokenExpiredException  = errors.New("Token has expired")
	TokenInvalidException  = errors.New("Token invalid")

	PasswordException  = errors.New("User password error")
	NeedAdminException = errors.New("Action need admin authority")
)

// Db
var (
	// exist
	UserExistException     = errors.New("User already existed")
	UserNotExistException  = errors.New("User not found")
	VideoNotExistException = errors.New("Video not found")

	// user crud failed
	ModifyPassException    = errors.New("User password modify failed")
	InsertUserException    = errors.New("User insert failed")
	NotUpdateUserException = errors.New("User information not updated")
	DeleteUserException    = errors.New("User delete failed")

	// video crud failed
	CreateVideoException         = errors.New("Video insert failed")
	NotUpdateVideoException      = errors.New("Video information not updated")
	DeleteVideoException         = errors.New("Video delete failed")
	NoAuthToActionVideoException = errors.New("Have no authorization to action video")

	// user other exception
	UserNameUsedException     = errors.New("Username has been used")
	UsetInfoException         = errors.New("User information invalid")
	SubscribeOneSelfException = errors.New("Cound not subscribe to oneself")

	// video exception
	VideoUrlUsedException = errors.New("Video resource url has been used")

	// raw exception
	ImageUploadException  = errors.New("Image upload failed")
	VideoUploadException  = errors.New("Video upload failed")
	FileExtException      = errors.New("Extension not supported")
	FileNotExistException = errors.New("File not exist")
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
