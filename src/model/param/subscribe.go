package param

// @Model         SubscribeParam
// @Description   关注请求参数
// @Property      to integer true "用户id" 1
type SubscribeParam struct {
	Uid int32 `form:"to" json:"to" binding:"required,min=1"`
}
