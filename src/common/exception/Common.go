package exception

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
)

// Request
var (
	RequestParamError  = NewError(400, "request param error")
	RequestFormatError = NewError(400, "request format error")
	RequestLargeError  = NewError(413, "request body too large")
)

// Authorization
var (
	UnAuthorizedError           = NewError(401, "unauthorized user")
	TokenExpiredError           = NewError(401, "token has expired")
	AuthorizedUserNotFoundError = NewError(401, "authorized user not found")
	NeedAdminError              = NewError(403, "need admin authority")

	PasswordError   = NewError(401, "password error")
	LoginError      = NewError(500, "login failed")           // R
	RegisterError   = NewError(500, "register failed")        // C
	UpdatePassError = NewError(500, "update password failed") // U
	LogoutError     = NewError(500, "logout failed")
)

// Model
var (
	// user
	UserNotFoundError = NewError(404, "user not found")     // R
	UserUpdateError   = NewError(500, "user update failed") // U
	UserDeleteError   = NewError(500, "user delete failed") // D

	UsernameUsedError  = NewError(400, "username has been used")
	SubscribeSelfError = NewError(400, "subscribe oneself invalid")
	SubscribeError     = NewError(500, "subscribe failed")
	UnSubscribeError   = NewError(500, "unsubscribe failed")

	// video
	VideoNotFoundError = NewError(404, "video not found")     // R
	VideoInsertError   = NewError(500, "video insert failed") // C
	VideoUpdateError   = NewError(500, "video update failed") // U
	VideoDeleteError   = NewError(500, "video delete failed") // D

	VideoUrlExistError = NewError(400, "video url has been used")
)

// File
var (
	ImageNotFoundError     = NewError(404, "image not found")
	ImageNotSupportedError = NewError(400, "image type not supported")
	ImageSaveError         = NewError(500, "image save failed")
)

func WrapValidationError(err error) *Error {
	isf := xgin.IsValidationFormatError(err)
	if isf {
		return RequestFormatError
	} else {
		return RequestParamError
	}
}
