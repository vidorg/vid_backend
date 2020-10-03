package po

import (
	"github.com/vidorg/vid_backend/lib/xgorm"
)

// Channel: tbl_channel.
type Channel struct {
	Cid         uint64 `gorm:"                    not null; primaryKey; autoIncrement"` // channel id
	Name        string `gorm:"type:varchar(255);  not null; uniqueIndex:uk_name"`       // channel name
	Description string `gorm:"type:varchar(1023); not null"`                            // channel description
	CoverUrl    string `gorm:"type:varchar(255);  not null"`                            // channel cover url
	AuthorUid   uint64 `gorm:"type:bigint;        not null"`                            // channel author uid

	Author      *User   `gorm:"                         foreignKey:AuthorUid; references:Uid"`                                         // po.Channel belongs to po.User
	Subscribers []*User `gorm:"many2many:tbl_subscribe; foreignKey:Cid;       references:Uid; joinForeignKey:cid; joinReferences:uid"` // tbl_subscribe

	xgorm.Model
}
