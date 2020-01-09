package po

import (
	"github.com/vidorg/vid_backend/src/model/dto/common"
)

type Video struct {
	Vid         int                 `gorm:"primary_key;auto_increment"`
	Title       string              `gorm:"not_null;type:varchar(100)"` // 100
	Description string              `gorm:"type:varchar(255)"`          // 255
	VideoUrl    string              `gorm:"not_null;type:varchar(255)"` // 255
	CoverUrl    string              `gorm:"not_null;type:varchar(255)"` // 255
	UploadTime  common.JsonDateTime `gorm:"not_null;type:datetime;default:CURRENT_TIMESTAMP"`
	AuthorUid   int

	Author *User `gorm:"foreignkey:AuthorUid"`

	GormTime
}
