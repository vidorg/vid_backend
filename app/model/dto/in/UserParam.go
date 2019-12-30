package in

import "time"

type UserParam struct {
	Username    string    `form:"username"     json:"username"     binding:"required"`
	Profile     string    `form:"profile"      json:"profile"      binding:"required"`
	Sex         string    `form:"sex"          json:"sex"          binding:"required"`
	BirthTime   time.Time `form:"birth_time"   json:"birth_time"   binding:"required" time_format:"2006-01-02"`
	PhoneNumber string    `form:"phone_number" json:"phone_number" binding:"required"`
}
