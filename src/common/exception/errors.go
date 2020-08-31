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
	RequestParamError   = New(400, cerr, "request param error")  // 请求参数错误
	RequestFormatError  = New(400, ce(), "request format error") // 请求参数格式错误
	ServerRecoveryError = New(500, serr, "server unknown error") // 服务器未知错误
)

// auth mw exceptions
var (
	CheckAuthorizeError = New(500, se(), "check authorize failed") // 检查身份失败
	InvalidTokenError   = New(401, ce(), "invalid token")          // 无效的令牌
	UnAuthorizedError   = New(401, ce(), "unauthorized")           // 未授权
	TokenExpiredError   = New(401, ce(), "token expired")          // 令牌过期
	CheckRoleError      = New(500, se(), "failed to check role")   // 检查权限失败
	NoPermissionError   = New(401, ce(), "no permission")          // 无权限
)

// auth exceptions
var (
	RegisterError        = New(500, se(), "register failed")                        // 注册失败
	EmailRegisteredError = New(409, ce(), "email has been registered")              // 邮箱已被注册
	LoginError           = New(500, se(), "login failed")                           // 登录失败
	LoginParameterError  = New(401, ce(), "email, username, uid or password wrong") // 用户名密码错误
	LogoutError          = New(500, se(), "logout failed")                          // 登出失败
	UpdatePasswordError  = New(500, se(), "update password failed")                 // 更新密码失败
	WrongPasswordError   = New(401, ce(), "password is wrong")                      // 密码不一致

	SendActivateEmailError = New(500, se(), "send email failed")                // 发送邮件失败
	AlreadyActivatedError  = New(400, ce(), "you have been activated")          // 已经被激活
	ActivateSuspendError   = New(400, ce(), "suspend user can not be activate") // 激活封禁的用户
	ActivateUserError      = New(500, se(), "activate user error")              // 用户已激活
	InvalidSpecError       = New(400, ce(), "invalid spec code")                // 无效的认证码
)

// user exceptions
var (
	QueryUserError    = New(500, se(), "query user failed") // 查找用户失败
	UserNotFoundError = New(404, ce(), "user not found")    // 用户未找到

	UpdateUserError   = New(500, se(), "update user failed")     // 更新用户失败
	UsernameUsedError = New(409, ce(), "username has been used") // 用户名已被使用
	DeleteUserError   = New(500, se(), "delete user failed")     // 删除用户失败
)

// subscribe exception
var (
	GetSubscriberListError  = New(500, se(), "get follower list failed")  // 获取粉丝列表失败
	GetSubscribingListError = New(500, se(), "get following list failed") // 获取关注列表失败

	SubscribeError          = New(500, se(), "follow failed")              // 关注用户失败
	SubscribeSelfError      = New(400, ce(), "could not follow self")      // 无法关注自己
	AlreadySubscribingError = New(409, ce(), "user has been followed")     // 已经关注的用户
	UnSubscribeError        = New(500, se(), "unfollow failed")            // 取消关注用户失败
	NotSubscribeYetError    = New(409, ce(), "user has not been followed") // 还没有关注的用户
)

// video exception
var (
	QueryVideoError    = New(500, se(), "query video failed") // 查找视频错误
	VideoNotFoundError = New(404, ce(), "video not found")    // 视频未找到

	InsertVideoError     = New(500, se(), "video insert failed")      // 插入视频失败
	UpdateVideoError     = New(500, se(), "video update failed")      // 更新视频失败
	VideoPermissionError = New(400, ce(), "no permission with video") // 无权限操作视频
	DeleteVideoError     = New(500, se(), "video delete failed")      // 删除视频失败
)

// rbac rule exceptions
var (
	QueryRbacRuleError  = New(500, se(), "query rbac rule failed")       // 查询规则失败
	ChangeRoleError     = New(500, se(), "change user role failed")      // 更改用户角色失败
	ChangeSelfRoleError = New(400, ce(), "could not change self's role") // 更改自己的角色

	RbacSubjectInsertError   = New(500, se(), "insert rbac subject failed") // 角色插入失败
	RbacSubjectExistedError  = New(409, ce(), "rbac subject exists")        // 角色已存在
	RbacSubjectDeleteError   = New(500, se(), "delete rbac subject failed") // 角色删除失败
	RbacSubjectNotFoundError = New(404, ce(), "rbac subject not found")     // 角色未找到

	RbacPolicyInsertError   = New(500, se(), "insert rbac policy failed") // 策略插入失败
	RbacPolicyExistedError  = New(409, ce(), "rbac policy exists")        // 策略已存在
	RbacPolicyDeleteError   = New(500, se(), "delete rbac policy failed") // 策略删除失败
	RbacPolicyNotFoundError = New(404, ce(), "rbac policy not found")     // 策略未找到
)

func WrapValidationError(err error) *Error {
	if xvalidator.ValidationRequiredError(err) {
		return RequestParamError
	}
	return RequestFormatError
}
