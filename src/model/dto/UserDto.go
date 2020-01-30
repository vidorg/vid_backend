package dto

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/model/po"
)

// @Model      UserDtoResult "返回用户信息"
// @Property   code    integer           true false "返回响应码"
// @Property   message string            true false "返回信息"
// @Property   data    object(#_UserDto) true false "返回数据"

// @Model      UserDtoPageResult "返回用户分页信息"
// @Property   code      integer               true false "返回响应码"
// @Property   message   string                true false "返回信息"
// @Property   data      object(#_UserDtoPage) true false "返回数据"

// @Model      _UserDtoPage "用户分页信息"
// @Property   total   integer          true false "数据总量"
// @Property   message string           true false "当前页"
// @Property   data    array(#_UserDto) true false "返回数据"

// @Model      LoginDtoResult "登录信息"
// @Property   code    integer            true false "返回响应码"
// @Property   message string             true false "返回信息"
// @Property   data    object(#_LoginDto) true false "返回数据"

// @Model      _LoginDto "登录信息"
// @Property   user   object(#_UserDto) true false "用户信息"
// @Property   token  string            true false "登录令牌"
// @Property   expire integer           true false "登录有效期，单位为秒"

// @Model      _UserDto "用户信息"
// @Property   uid           integer true false "用户id"
// @Property   username      string  true false "用户名"
// @Property   sex           string  true false "用户性别，枚举 {male, female, unknown}"
// @Property   profile       string  true true  "用户简介"
// @Property   avatar_url    string  true false "用户头像"
// @Property   birth_time    string  true false "用户生日，固定格式为 2000-01-01"
// @Property   authority     string  true false "用户权限，枚举 {normal, admin}"
// @Property   phone_number  string  true false "用户手机号码，部分接口可见"
// @Property   register_time string  true false "用户注册时间，固定格式为 2000-01-01 00:00:00"

type UserDto struct {
	Uid          int32  `json:"uid"`
	Username     string `json:"username"`
	Sex          string `json:"sex"`
	Profile      string `json:"profile"`
	AvatarUrl    string `json:"avatar_url"`
	BirthTime    string `json:"birth_time"`
	Authority    string `json:"authority"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	RegisterTime string `json:"register_time"`
}

// show all info
// Only used in QueryAllUsers()
func UserDtoAdminMapOption() *xmapper.DisposableMapOption {
	return xmapper.NewMapOption(&po.User{}, &UserDto{}, func(i interface{}, j interface{}) interface{} {
		user := i.(po.User)
		userDto := j.(UserDto)
		userDto.PhoneNumber = user.PhoneNumber
		return userDto
	})
}

// show info dependent on authUser
// Only used in QueryUser()
func UserDtoExtraMapOption(authUser *po.User) *xmapper.DisposableMapOption {
	return xmapper.NewMapOption(&po.User{}, &UserDto{}, func(i interface{}, j interface{}) interface{} {
		if authUser == nil { // not login, nothing (default)
			return j
		}
		user := i.(po.User)
		userDto := j.(UserDto)
		if authUser.Authority == enum.AuthAdmin { // admin, all info
			userDto.PhoneNumber = user.PhoneNumber // add phone number
			return userDto
		} else { // normal, only me
			if user.Uid == authUser.Uid {
				userDto.PhoneNumber = user.PhoneNumber
			}
			return userDto
		}
	})
}
