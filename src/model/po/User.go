package po

import (
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/common/model"
	"time"
)

type User struct {
	Uid         int32          `gorm:"primary_key;auto_increment"`
	Username    string         `gorm:"not_null;type:varchar(30);unique_index:idx_user_username_deleted_at_unique"` // 30
	Sex         enum.SexType   `gorm:"not_null;type:enum('unknown','male','female');default:'unknown'"`
	Profile     string         `gorm:"type:varchar(255)"`          // 255
	AvatarUrl   string         `gorm:"not_null;type:varchar(255)"` // 255
	BirthTime   model.JsonDate `gorm:"not_null;type:datetime;default:'2000-01-01 00:00:00'"`
	Authority   enum.AuthType  `gorm:"not_null;type:enum('admin', 'normal');default:'normal'"`
	RegisterIP  string         `gorm:"type:varchar(15)"` // 15
	PhoneNumber string         `gorm:"type:varchar(11)"` // 11

	// tbl_subscribe
	Subscribings []*User `gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"` // up_uid -> subscriber_uid
	Subscribers  []*User `gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid

	model.GormTimeWithoutDeletedAt
	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:00';unique_index:idx_user_username_deleted_at_unique"`
}
