package param

// https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags

type UserParam struct {
	Username    string  `form:"username"     json:"username"     binding:"required,min=5,max=30,name"`
	Profile     *string `form:"profile"      json:"profile"      binding:"required,min=0,max=255"`
	Sex         string  `form:"sex"          json:"sex"          binding:"required"`
	BirthTime   string  `form:"birth_time"   json:"birth_time"   binding:"required,date"`
	PhoneNumber string  `form:"phone_number" json:"phone_number" binding:"required,phone"`
	AvatarUrl   string  `form:"avatar_url"   json:"avatar_url"   binding:"required,url"`
}
