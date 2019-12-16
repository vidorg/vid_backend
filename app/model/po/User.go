package po

import (
	"time"
	"vid/app/model/enum"

	// "vid/app/config"
)

type User struct {
	Uid       int           `json:"uid"      gorm:"primary_key;auto_increment"`
	Username  string        `json:"username" gorm:"type:varchar(30);unique;not_null"` // 30
	Sex       enum.SexType  `json:"sex"      gorm:"type:enum('unknown','male','female');default:'unknown'"`
	Profile   string        `json:"profile"`    // 255
	AvatarUrl string        `json:"avatar_url"` // 255
	BirthTime time.Time     `json:"birth_time" gorm:"default:'2000-01-01'"`
	Authority enum.AuthType `json:"authority" gorm:"type:enum('admin', 'normal');default:'normal';not_null"`

	// inner system
	RegisterIP  string `json:"-"`
	PhoneNumber string `json:"-" gorm:"type:varchar(11)"` // 11

	// tbl_subscribe
	Subscribers  []*User `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid
	Subscribings []*User `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"` // up_uid -> subscriber_uid

	TimePo
}
