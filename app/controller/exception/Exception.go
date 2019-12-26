package exception

import (
	"errors"
)

// Request
var (
	RouteParamError = errors.New("request route param error")
	QueryParamError = errors.New("request query param error")
	FormParamError  = errors.New("request form data error")
	JsonParamError  = errors.New("request raw json error")

	FormatError           = errors.New("request format error")
	RequestSizeLargeError = errors.New("request body too large")
)

// Authorization
var (
	AuthorizationError = errors.New("authorization failed")
	TokenExpiredError  = errors.New("token has expired")

	PasswordError   = errors.New("password error")
	LoginError      = errors.New("login failed")
	RegisterError   = errors.New("register failed")
	UpdatePassError = errors.New("update password failed")

	NeedAdminError = errors.New("need admin authority")
)

// Model
var (
	// user
	UserUpdateError   = errors.New("user update failed") // U
	UserDeleteError   = errors.New("user delete failed") // D
	UserNotFoundError = errors.New("user not found")     // R

	UserNameUsedError  = errors.New("username has been used")
	SubscribeSelfError = errors.New("subscribe oneself invalid")
	SubscribeError     = errors.New("subscribe failed")
	UnSubscribeError   = errors.New("unsubscribe failed")

	// video
	VideoNotFoundError = errors.New("video not found")     // R
	VideoInsertError   = errors.New("video insert failed") // C
	VideoUpdateError   = errors.New("video update failed") // U
	VideoDeleteError   = errors.New("video delete failed") // D

	VideoExistError = errors.New("video resource has been used")

	// playlist
	PlaylistInsertError   = errors.New("playlist insert failed") // C
	PlaylistUpdateError   = errors.New("playlist update failed") // U
	PlaylistDeleteError   = errors.New("playlist delete failed") // D
	PlaylistNotFoundError = errors.New("playlist not found")     // R

	PlaylistNameUsedError    = errors.New("playlist name duplicated")
	PlaylistVideoDeleteError = errors.New("playlist video delete failed")
)

// File
var (
	ImageNotFoundError     = errors.New("image not found")
	ImageSaveError         = errors.New("image save failed")
	ImageNotSupportedError = errors.New("image type not supported")

	VideoUploadError   = errors.New("video upload failed")
	FileExtensionError = errors.New("extension not supported")
)
