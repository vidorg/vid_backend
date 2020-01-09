package param

type LoginParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"username" binding:"required"`
	Expire   int64  `form:"expire"   json:"expire"`
}

type RegisterParam struct {
	Username string `form:"username" json:"username" binding:"required,min=5,max=30,name"`
	Password string `form:"password" json:"username" binding:"required,min=8,max=30,pwd"`
}

type PassParam struct {
	Password string `form:"password" json:"username" binding:"required,min=8,max=30,pwd"`
}
