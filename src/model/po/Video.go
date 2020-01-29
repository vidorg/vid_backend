package po

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
)

type Video struct {
	Vid         int32                  `gorm:"primary_key;auto_increment"`
	Title       string                 `gorm:"not_null;type:varchar(100)"` // 100
	Description string                 `gorm:"type:varchar(1024)"`         // 1024
	VideoUrl    string                 `gorm:"not_null;type:varchar(255)"` // 255
	CoverUrl    string                 `gorm:"not_null;type:varchar(255)"` // 255
	UploadTime  xdatetime.JsonDateTime `gorm:"not_null;type:datetime;default:CURRENT_TIMESTAMP"`
	AuthorUid   int32

	Author *User `gorm:"foreignkey:AuthorUid"`

	xgorm.GormTime
}
