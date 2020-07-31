package exception

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
)

// Request / Response
var (
	RequestParamError   = New(400, 40000, "request param error")  // 请求参数错误
	RequestFormatError  = New(400, 40001, "request format error") // 请求格式错误
	ServerRecoveryError = New(500, 50000, "server unknown error") // 服务器未知错误
)

// Authorization
var (
	UnAuthorizedError        = New(401, 40002, "unauthorized user")         // 未认证
	InvalidTokenError        = New(401, 40003, "token is invalid")          // 令牌无效
	TokenExpiredError        = New(401, 40004, "token is expired")          // 令牌过期
	CheckUserRoleError       = New(500, 50001, "failed to check user role") // 检查用户角色错误
	RoleHasNoPermissionError = New(403, 40005, "role has no permission")    // 用户没有权限

	WrongPasswordError = New(401, 40006, "wrong password")         // 密码错误
	LoginError         = New(500, 50002, "login failed")           // 登录失败
	RegisterError      = New(500, 50003, "register failed")        // 注册失败
	UpdatePassError    = New(500, 50004, "update password failed") // 修改密码失败
	LogoutError        = New(500, 50005, "logout failed")          // 注销失败
)

// Model
var (
	// user
	UserNotFoundError  = New(404, 40007, "user not found")            // 用户未找到
	UserUpdateError    = New(500, 50006, "user update failed")        // 修改用户失败
	UserDeleteError    = New(500, 50007, "user delete failed")        // 删除用户失败
	UsernameUsedError  = New(409, 40008, "username has been used")    // 用户名已使用
	SubscribeSelfError = New(400, 40009, "subscribe oneself invalid") // 关注自我无效
	SubscribeError     = New(500, 50008, "subscribe failed")          // 关注失败
	UnSubscribeError   = New(500, 50009, "unsubscribe failed")        // 取消关注失败

	// video
	VideoNotFoundError = New(404, 40010, "video not found")         // 视频未找到
	VideoInsertError   = New(500, 50010, "video insert failed")     // 视频插入失败
	VideoUpdateError   = New(500, 50011, "video update failed")     // 视频更新失败
	VideoDeleteError   = New(500, 50012, "video delete failed")     // 视频删除失败
	VideoUrlExistError = New(409, 40011, "video url has been used") // 视频资源已使用

	// policy
	PolicySetRoleError  = New(403, 40012, "set root role failed") // 设置角色失败
	PolicyNotFountError = New(404, 40013, "policy not found")     // 策略未找到
	PolicyExistedError  = New(409, 40014, "policy has existed")   // 策略已存在
	PolicyInsertError   = New(500, 50013, "insert policy failed") // 策略新建失败
	PolicyDeleteError   = New(500, 50014, "delete policy failed") // 策略删除失败
)

func WrapValidationError(err error) *Error {
	isf := xgin.IsValidationFormatError(err)
	if isf {
		return RequestFormatError
	} else {
		return RequestParamError
	}
}
