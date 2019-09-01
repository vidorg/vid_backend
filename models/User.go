package models

import (
	"strings"
	"time"
	"vid/config"
)

// http://gorm.io/docs/many_to_many.html#Self-Referencing
type User struct {
	Uid          int       `json:"uid" gorm:"primary_key;AUTO_INCREMENT"`
	Username     string    `json:"username" gorm:"type:varchar(50);unique"`
	Profile      string    `json:"profile" gorm:"type:varchar(120)"`
	RegisterTime time.Time `json:"register_time" gorm:"type:datetime"`
	Subscribers  []*User   `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"`
	Subscribings []*User   `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"`
}

// @override
func (u *User) CheckValid() bool {
	return u.Username != "" && strings.Index(u.Username, " ") == -1
}

func (u *User) Equals(obj *User) bool {
	return u.Uid == obj.Uid && u.Username == obj.Username && u.Profile == obj.Profile
}

func (u *User) CheckFormat() bool {
	cfg := config.AppCfg
	return len(u.Username) >= cfg.FormatConfig.MinLen_Username &&
		len(u.Username) <= cfg.FormatConfig.MaxLen_Username
}
