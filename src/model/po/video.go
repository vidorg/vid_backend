package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
)

type Video struct {
	Vid         uint64 `gorm:"primary_key; auto_increment"`                      // video id
	Title       string `gorm:"type:varchar(255);  not null"`                     // video title
	Description string `gorm:"type:varchar(1023); not null"`                     // video description
	VideoUrl    string `gorm:"type:varchar(255);  not null"`                     // video source url // TODO https://cloud.tencent.com/document/product/266/31766
	CoverUrl    string `gorm:"type:varchar(255);  not null"`                     // video cover url (oss)
	AuthorUid   uint64 `gorm:"                    not null"`                     // video author id
	Author      *User  `gorm:"foreignkey:AuthorUid; association_foreignkey:Uid"` // po.Video belongs to po.User (has many)

	Favoreds []*User `gorm:"many2many:favorite; foreignkey:Vid; association_foreignkey:Uid; jointable_foreignkey:vid; association_jointable_foreignkey:uid"` // tbl_favorite

	xgorm.GormTime
}
