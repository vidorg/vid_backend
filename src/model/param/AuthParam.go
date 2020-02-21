package param

// @Model         LoginParam
// @Description   登录请求参数
// @Property      username string  true  "用户名"
// @Property      password string  true  "密码"
// @Property      expire   integer false "登录有效期，默认为七天" 604800
type LoginParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Expire   int64  `form:"expire"   json:"expire"`
}

// @Model         RegisterParam
// @Description   注册请求参数
// @Property      username string true "用户名，长度在 [5, 30] 之间"
// @Property      password string true "密码，长度在 [8, 30] 之间"
type RegisterParam struct {
	Username string `form:"username" json:"username" binding:"required,min=5,max=30,name"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=30,pwd"`
}

// @Model         PassParam
// @Description   修改密码请求参数
// @Property      password string true "密码，长度在 [8, 30] 之间"
type PassParam struct {
	Password string `form:"password" json:"password" binding:"required,min=8,max=30,pwd"`
}
