package exception

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
)

// Request
var (
	RequestParamError  = errors.New("request param error")
	RequestFormatError = errors.New("request format error")
	RequestSizeError   = errors.New("request body too large")
)

// Authorization
var (
	UnAuthorizedError   = errors.New("unauthorized user")
	TokenExpiredError   = errors.New("token has expired")
	AuthorizedUserError = errors.New("authorized user not found")

	PasswordError   = errors.New("password error")
	LoginError      = errors.New("login failed")
	RegisterError   = errors.New("register failed") // C
	LogoutError     = errors.New("logout failed")
	UpdatePassError = errors.New("update password failed")

	NeedAdminError = errors.New("need admin authority")
)

// Model
var (
	// user
	UserUpdateError   = errors.New("user update failed") // U
	UserDeleteError   = errors.New("user delete failed") // D
	UserNotFoundError = errors.New("user not found")     // R

	UsernameUsedError  = errors.New("username has been used")
	SubscribeSelfError = errors.New("subscribe oneself invalid")
	SubscribeError     = errors.New("subscribe failed")
	UnSubscribeError   = errors.New("unsubscribe failed")

	// video
	VideoNotFoundError = errors.New("video not found")     // R
	VideoInsertError   = errors.New("video insert failed") // C
	VideoUpdateError   = errors.New("video update failed") // U
	VideoDeleteError   = errors.New("video delete failed") // D

	VideoExistError = errors.New("video has been updated")
)

// File
var (
	ImageNotFoundError     = errors.New("image not found")
	ImageNotSupportedError = errors.New("image type not supported")
	ImageSaveError         = errors.New("image save failed")
)

func WrapValidationError(err error) error {
	isf := xgin.IsValidationFormatError(err)
	if isf {
		return RequestFormatError
	} else {
		return RequestParamError
	}
}
