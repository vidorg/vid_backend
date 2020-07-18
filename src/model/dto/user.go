package dto

// @Model         _LoginDto
// @Description   登录信息
// @Property      user   object(#_UserDto) true "用户信息"
// @Property      token  string            true "登录令牌"
// @Property      expire integer           true "登录有效期，单位为秒"

// @Model         _UserDto
// @Description   用户信息
// @Property      uid           integer                          true "用户id"
// @Property      username      string                           true "用户名"
// @Property      sex           string(enum:male,female,unknown) true "用户性别"
// @Property      profile       string                           true "用户简介"
// @Property      avatar_url    string                           true "用户头像"
// @Property      birthday      string(format:date)              true "用户生日"
// @Property      role          string                           true "用户角色"
// @Property      phone_number  string                           true "用户手机号码，部分接口可见"
// @Property      register_time string(format:datetime)          true "用户注册时间"
type UserDto struct {
	Uid          int32  `json:"uid"`
	Username     string `json:"username"`
	Sex          string `json:"sex"`
	Profile      string `json:"profile"`
	AvatarUrl    string `json:"avatar_url"`
	Birthday     string `json:"birthday"`
	Role         string `json:"role"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	RegisterTime string `json:"register_time"`
}

// @Model         _UserAndExtraDto
// @Description   用户与数量信息
// @Property      user  object(#_UserDto)      true "用户信息"
// @Property      extra object(#_UserExtraDto) true "用户额外信息"

// @Model         _UserExtraDto
// @Description   用户额外信息
// @Property      subscribing_cnt integer true "用户关注数量"
// @Property      subscriber_cnt  integer true "用户粉丝数量"
// @Property      video_cnt       integer true "用户视频数量"
type UserExtraDto struct {
	SubscribingCount int32 `json:"subscribing_cnt"`
	SubscriberCount  int32 `json:"subscriber_cnt"`
	VideoCount       int32 `json:"video_cnt"`
}

type UserDetailDto struct {
	User  *UserDto      `json:"user"`
	Extra *UserExtraDto `json:"extra"`
}

type UserLoginDto struct {
	User   *UserDto `json:"user"`
	Token  string   `json:"token"`
	Expire int64    `json:"expire"`
}
