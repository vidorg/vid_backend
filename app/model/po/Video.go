package po

import (
	"vid/app/model/vo"
)

type Video struct {
	Vid         int         `json:"vid"           gorm:"primary_key;auto_increment"`
	Title       string      `json:"title"         gorm:"type:varchar(100);not_null"` // 100
	Description string      `json:"description"`                                     // 255
	VideoUrl    string      `json:"video_url"     gorm:"not_null;unique"`            // 255
	CoverUrl    string      `json:"cover_url"     gorm:"not_null"`                   // 255
	UploadTime  vo.JsonDate `json:"upload_time"   gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	AuthorUid   int         `json:"-"`

	Author *User `json:"author" gorm:"foreignkey:AuthorUid"`

	GormTime `json:"-"`
}
