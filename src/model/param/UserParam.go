package param

// https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags

// @Model      UserParam "用户请求参数"
// @Property   username     string true false "用户名，长度在 [5, 30] 之间"
// @Property   profile      string true true  "用户简介，长度在 [0, 255] 之间"
// @Property   sex          string true false "用户性别，允许值为 {male, female, unknown}"
// @Property   birth_time   string true false "用户生日，固定格式为 2000-01-01"
// @Property   phone_number string true false "用户手机号码，长度为 11，仅限中国大陆手机号码"
// @Property   avatar_url   string true false "用户头像链接"
type UserParam struct {
	Username    string  `form:"username"     json:"username"     binding:"required,min=5,max=30,name"`
	Profile     *string `form:"profile"      json:"profile"      binding:"required,min=0,max=255"`
	Sex         string  `form:"sex"          json:"sex"          binding:"required"`
	BirthTime   string  `form:"birth_time"   json:"birth_time"   binding:"required,date"`
	PhoneNumber string  `form:"phone_number" json:"phone_number" binding:"required,phone"`
	AvatarUrl   string  `form:"avatar_url"   json:"avatar_url"   binding:"required,url"`
}
