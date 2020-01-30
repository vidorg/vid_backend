package param

// @Model      SubParam "关注请求参数"
// @Property   to integer true false "用户id" 1
type SubParam struct {
	Uid int32 `form:"to" json:"to" binding:"required,min=1"`
}
