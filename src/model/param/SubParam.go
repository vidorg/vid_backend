package param

type SubParam struct {
	Uid int32 `form:"to" json:"to" binding:"required,min=1"`
}
