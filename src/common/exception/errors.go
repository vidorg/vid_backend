package exception

import (
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
)

var (
	cerr = int32(40000) // client error code
	serr = int32(50000) // server error code
)

// Return ++ cerr.
func ce() int32 { cerr++; return cerr }

// Return ++ serr.
func se() int32 { serr++; return serr }

// global exceptions
var (
	RequestParamError   = New(400, cerr, "request param error")
	RequestFormatError  = New(400, ce(), "request format error")
	ServerRecoveryError = New(500, serr, "server unknown error")
)

// auth mw exceptions
var (
	CheckAuthorizeError = New(500, se(), "check authorize failed")
	InvalidTokenError   = New(401, ce(), "invalid token")
	UnAuthorizedError   = New(401, ce(), "unauthorized")
	TokenExpiredError   = New(401, ce(), "token expired")
	CheckRoleError      = New(500, se(), "failed to check role")
	NoPermissionError   = New(401, ce(), "no permission")
)

// auth exceptions
var (
	RegisterError        = New(500, se(), "register failed")
	EmailRegisteredError = New(409, ce(), "email has been registered")
	LoginError           = New(500, se(), "login failed")
	LoginParameterError  = New(401, ce(), "email, username, uid or password wrong")
	LogoutError          = New(500, se(), "logout failed")

	UpdatePasswordError = New(500, se(), "update password failed")
	WrongPasswordError  = New(401, ce(), "password is wrong")
)

// user exceptions
var (
	QueryUserError = New(500, se(), "query user failed")

	UserNotFoundError = New(404, ce(), "user not found")
	UpdateUserError   = New(500, se(), "update user failed")
	UsernameUsedError = New(409, ce(), "username has been used")
	DeleteUserError   = New(500, se(), "delete user failed")
)

// subscribe exception
var (
	GetSubscriberListError  = New(500, se(), "get follower list failed")
	GetSubscribingListError = New(500, se(), "get following list failed")

	SubscribeError          = New(500, se(), "follow failed")
	SubscribeSelfError      = New(400, ce(), "could not follow self")
	AlreadySubscribingError = New(409, ce(), "user has been followed")
	UnSubscribeError        = New(500, se(), "unfollow failed")
	UnSubscribeSelfError    = New(400, ce(), "could not unfollow self")
	NotSubscribeYetError    = New(409, ce(), "user has not been followed")
)

// video exception
var (
	VideoNotFoundError = New(404, 40010, "video not found")         // 视频未找到
	VideoInsertError   = New(500, 50010, "video insert failed")     // 视频插入失败
	VideoUpdateError   = New(500, 50011, "video update failed")     // 视频更新失败
	VideoDeleteError   = New(500, 50012, "video delete failed")     // 视频删除失败
	VideoUrlExistError = New(409, 40011, "video url has been used") // 视频资源已使用
)

func WrapValidationError(err error) *Error {
	if xvalidator.ValidationRequiredError(err) {
		return RequestParamError
	}
	return RequestFormatError
}
