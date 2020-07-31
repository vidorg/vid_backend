package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
)

type Video struct {
	Vid         int32  `gorm:"primary_key;auto_increment"`
	Title       string `gorm:"not_null;type:varchar(100)"` // 100
	Description string `gorm:"type:varchar(1024)"`         // 1024
	VideoUrl    string `gorm:"not_null;type:varchar(255)"` // 255 // TODO url
	CoverUrl    string `gorm:"not_null;type:varchar(255)"` // 255 // TODO url
	Author      *User  `gorm:"foreignkey:AuthorUid"`
	AuthorUid   int32

	xgorm.GormTime
}
