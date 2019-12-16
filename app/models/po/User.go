package po

import (
	"encoding/json"
	"strings"
	"time"
	"vid/app/utils"

	// "vid/app/config"
)

type User struct {
	TimePo

	Uid          int         `json:"uid" gorm:"primary_key;auto_increment"`
	Username     string      `json:"username" gorm:"type:varchar(50);unique;not_null"` // 50
	Profile      string      `json:"profile"`                                          // 255
	Sex          string      `json:"sex" gorm:"type:char(5);default:'X'"`              // 1
	AvatarUrl    string      `json:"avatar_url"`                                       // 255
	BirthTime    time.Time   `json:"birth_time" gorm:"type:datetime;default:'2000-01-01'"`
	RegisterTime time.Time   `json:"register_time" gorm:"type:datetime;default:'2000-01-01'"`
	Authority    AuthType    `json:"authority" gorm:"type:ENUM('admin', 'normal');default:'normal';not_null"`
	Subscribers  []*User     `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"`
	Subscribings []*User     `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"`

	// inner system
	RegisterIP  string `json:"-" gorm:"type:char(16)"`
	PhoneNumber int    `json:"-" gorm:"default:-1"`
}

func (u *User) Unmarshal(jsonBody string, needUid bool) bool {
	err := json.Unmarshal([]byte(jsonBody), u)
	if err != nil || (needUid && u.Uid == 0) {
		return false
	}
	if strings.Index(u.Username, " ") != -1 {
		return false
	}
	// No Profile Field
	if strings.Index(jsonBody, "\"profile\": \"") == -1 {
		// u.Profile = config.AppConfig.MagicToken
	}
	return true
}

func (u *User) Equals(obj *User) bool {
	return u.Uid == obj.Uid &&
		u.Username == obj.Username &&
		u.Profile == obj.Profile &&
		u.Sex == obj.Sex &&
		u.AvatarUrl == obj.AvatarUrl &&
		u.BirthTime.Equal(obj.BirthTime)
}

func (u *User) CheckFormat() bool {
	// cfg := config.AppConfig
	// return len(u.Username) >= cfg.FormatConfig.MinLen_Username &&
	// 	len(u.Username) <= cfg.FormatConfig.MaxLen_Username &&
	// 	(u.Sex == "M" || u.Sex == "F" || u.Sex == "X" || strings.Trim(u.Sex, " ") == "")
	return true
}

// Server -> DB
func (u *User) ToDB() {
	if strings.HasPrefix(u.AvatarUrl, "http") {
		sp := strings.Split(u.AvatarUrl, "/")
		if sp[len(sp)-2] == "-1" {
			// http://127.0.0.1:3344/raw/image/-1/avatar.jpg -> /avatar.jpg
			u.AvatarUrl = "/" + sp[len(sp)-1]
		} else {
			// http://127.0.0.1:3344/raw/image/233/avatar.jpg -> avatar.jpg
			u.AvatarUrl = sp[len(sp)-1]
		}
	}
}

// DB -> Server
func (u *User) ToServer() {
	if !strings.HasPrefix(u.AvatarUrl, "http") {
		if strings.Index(u.AvatarUrl, "/") != -1 {
			// /avatar.jpg -> http://127.0.0.1:3344/raw/image/-1/avatar.jpg
			u.AvatarUrl = utils.CmnUtil.GetImageUrl(-1, strings.TrimLeft(u.AvatarUrl, "/"))
		} else {
			// avatar.jpg -> http://127.0.0.1:3344/raw/image/233/avatar.jpg
			u.AvatarUrl = utils.CmnUtil.GetImageUrl(u.Uid, u.AvatarUrl)
		}
	}
}

// Authority

type AuthType string

const (
	AuthAdmin  AuthType = "admin"
	AuthNormal AuthType = "normal"
)
