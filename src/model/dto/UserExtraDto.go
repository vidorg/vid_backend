package dto

// @Model      UserExtraDtoResult "返回用户与数量信息"
// @Property   code     integer                   true false "返回响应码"
// @Property   message  string                    true false "返回信息"
// @Property   data     object(#_UserAndExtraDto) true false "返回数据"

// @Model      _UserAndExtraDto "用户与数量信息"
// @Property   user  object(#_UserDto)      true false "用户信息"
// @Property   extra object(#_UserExtraDto) true false "用户额外信息"

// @Model      ImageDtoResult "返回上传图片信息"
// @Property   code    integer            true false "返回响应码"
// @Property   message string             true false "返回信息"
// @Property   data    object(#_ImageDto) true false "返回数据"

// @Model      _ImageDto "上传图片信息"
// @Property   url  string  true false "图片链接"
// @Property   size integer true false "图片大小，单位为字节"

// @Model      _UserExtraDto "用户额外信息"
// @Property   subscribing_cnt integer true false "用户关注数量"
// @Property   subscriber_cnt  integer true false "用户粉丝数量"
// @Property   video_cnt       integer true false "用户视频数量"

type UserExtraDto struct {
	SubscribingCount int32 `json:"subscribing_cnt"`
	SubscriberCount  int32 `json:"subscriber_cnt"`
	VideoCount       int32 `json:"video_cnt"`
}
