package po

import (
	"fmt"
	"reflect"
	"strings"
	"vid/app/model/dto"
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

func (Video) UrlConverter() dto.Converter {
	return dto.Converter{
		FieldType: reflect.TypeOf(&Video{}),
		Converter: func(obj interface{}) {
			if content, ok := obj.(*Video); ok {
				if !strings.HasPrefix(content.CoverUrl, "http") {
					if content.CoverUrl == "" {
						content.CoverUrl = "http://localhost:3344/raw/image/default/cover.jpg"
					} else {
						content.CoverUrl = fmt.Sprintf("http://localhost:3344/raw/image/%d/%s", content.AuthorUid, content.CoverUrl)
					}
				}
				if !strings.HasPrefix(content.VideoUrl, "http") {
					if content.VideoUrl != "" {
						content.VideoUrl = fmt.Sprintf("http://localhost:3344/raw/video/%d/%s", content.AuthorUid, content.VideoUrl)
					}
				}
			}
		},
	}
}
