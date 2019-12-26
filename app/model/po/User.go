package po

import (
	"fmt"
	"reflect"
	"vid/app/model/dto"
	"vid/app/model/enum"
	"vid/app/model/vo"
)

type User struct {
	Uid       int           `json:"uid"         gorm:"primary_key;auto_increment"`
	Username  string        `json:"username"    gorm:"type:varchar(30);unique;not_null"` // 30
	Sex       enum.SexType  `json:"sex"         gorm:"type:enum('unknown','male','female');default:'unknown'"`
	Profile   string        `json:"profile"`                     // 255
	AvatarUrl string        `json:"avatar_url"  gorm:"not_null"` // 255
	BirthTime vo.JsonDate   `json:"birth_time"  gorm:"not_null;default:'2000-01-01 00:00:00'"`
	Authority enum.AuthType `json:"authority"   gorm:"type:enum('admin', 'normal');default:'normal';not_null"`

	// inner system
	RegisterIP  string `json:"-"`
	PhoneNumber string `json:"-" gorm:"type:varchar(11)"` // 11

	// tbl_subscribe
	Subscribings []*User `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"` // up_uid -> subscriber_uid
	Subscribers  []*User `json:"-" gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid

	GormTime `json:"-"`
}

func (User) AvatarUrlConverter() dto.Converter {
	return dto.Converter{
		FieldType: reflect.TypeOf(&User{}),
		Converter: func(obj interface{}) {
			fmt.Println(reflect.TypeOf(obj))
			if content, ok := obj.(*User); ok {
				if content.AvatarUrl == "" {
					content.AvatarUrl = "http://localhost:3344/raw/image/default/avatar.jpg"
				} else {
					content.AvatarUrl = fmt.Sprintf("http://localhost:3344/raw/image/%d/%s", content.Uid, content.AvatarUrl)
				}
			}
		},
	}
}
