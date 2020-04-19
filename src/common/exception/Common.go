package exception

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
)

// Request / Response
var (
	RequestParamError   = NewError(400, "request param error")    // 请求参数错误
	RequestFormatError  = NewError(400, "request format error")   // 请求格式错误
	RequestLargeError   = NewError(413, "request body too large") // 请求体过大
	ServerRecoveryError = NewError(500, "server unknown error")   // 服务器未知错误
)

// Authorization
var (
	UnAuthorizedError        = NewError(401, "unauthorized user")         // 未认证
	TokenExpiredError        = NewError(401, "token has expired")         // 令牌过期
	CheckUserRoleError       = NewError(500, "failed to check user role") // 检查用户角色错误
	RoleHasNoPermissionError = NewError(403, "role has no permission")    // 用户没有权限

	WrongPasswordError = NewError(401, "wrong password")         // 密码错误
	LoginError         = NewError(500, "login failed")           // 登录失败
	RegisterError      = NewError(500, "register failed")        // 注册失败
	UpdatePassError    = NewError(500, "update password failed") // 修改密码失败
	LogoutError        = NewError(500, "logout failed")          // 注销失败
)

// Model
var (
	// user
	UserNotFoundError  = NewError(404, "user not found")            // 用户未找到
	UserUpdateError    = NewError(500, "user update failed")        // 修改用户失败
	UserDeleteError    = NewError(500, "user delete failed")        // 删除用户失败
	UsernameUsedError  = NewError(409, "username has been used")    // 用户名已使用
	SubscribeSelfError = NewError(400, "subscribe oneself invalid") // 关注自我无效
	SubscribeError     = NewError(500, "subscribe failed")          // 关注失败
	UnSubscribeError   = NewError(500, "unsubscribe failed")        // 取消关注失败

	// video
	VideoNotFoundError = NewError(404, "video not found")         // 视频未找到
	VideoInsertError   = NewError(500, "video insert failed")     // 视频插入失败
	VideoUpdateError   = NewError(500, "video update failed")     // 视频更新失败
	VideoDeleteError   = NewError(500, "video delete failed")     // 视频删除失败
	VideoUrlExistError = NewError(409, "video url has been used") // 视频资源已使用

	// policy
	PolicySetRoleError  = NewError(403, "set root role failed") // 设置角色失败
	PolicyNotFountError = NewError(404, "policy not found")     // 策略未找到
	PolicyExistedError  = NewError(409, "policy has existed")   // 策略已存在
	PolicyInsertError   = NewError(500, "insert policy failed") // 策略新建失败
	PolicyDeleteError   = NewError(500, "delete policy failed") // 策略删除失败
)

// File
var (
	ImageNotFoundError     = NewError(404, "image not found")          // 图片未找到
	ImageNotSupportedError = NewError(400, "image type not supported") // 不支持的图片格式
	ImageSaveError         = NewError(500, "image save failed")        // 图片保存失败
)

func WrapValidationError(err error) *Error {
	isf := xgin.IsValidationFormatError(err)
	if isf {
		return RequestFormatError
	} else {
		return RequestParamError
	}
}
