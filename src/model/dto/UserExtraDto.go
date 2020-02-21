package dto

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
