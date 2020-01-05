package param

type PasswordParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"username" binding:"required"`
	Expire   int64  `form:"expire"   json:"expire"`
}
