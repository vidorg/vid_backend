package po

import (
	"github.com/vidorg/vid_backend/src/model/common"
)

type Video struct {
	Vid         int32               `gorm:"primary_key;auto_increment"`
	Title       string              `gorm:"not_null;type:varchar(100)"` // 100
	Description string              `gorm:"type:varchar(1024)"`         // 1024
	VideoUrl    string              `gorm:"not_null;type:varchar(255)"` // 255
	CoverUrl    string              `gorm:"not_null;type:varchar(255)"` // 255
	UploadTime  common.JsonDateTime `gorm:"not_null;type:datetime;default:CURRENT_TIMESTAMP"`
	AuthorUid   int32

	Author *User `gorm:"foreignkey:AuthorUid"`

	GormTime
}
