package param

type SubParam struct {
	Uid int `form:"to" json:"to" binding:"required,min=1"`
}
