package po

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/vidorg/vid_backend/src/common/constant"
	"time"
)

type User struct {
	Uid         int32              `gorm:"primary_key;auto_increment"`
	Username    string             `gorm:"not_null;type:varchar(30);uniquEWWe_index:idx_user_username_deleted_at_unique"` // 30
	Sex         constant.SexEnum   `gorm:"not_null;type:enum('unknown','male','female');default:'unknown'"`
	Profile     string             `gorm:"type:varchar(255)"`          // 255
	AvatarUrl   string             `gorm:"not_null;type:varchar(255)"` // 255
	Birthday    xdatetime.JsonDate `gorm:"not_null;type:date;default:'2000-01-01'"`
	Role        string             `gorm:"not_null;default:'normal'"`
	RegisterIP  string             `gorm:"type:varchar(15)"` // 15
	PhoneNumber string             `gorm:"type:varchar(11)"` // 11

	// tbl_subscribe
	Subscribings []*User `gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"` // up_uid -> subscriber_uid
	Subscribers  []*User `gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid

	xgorm.GormTimeWithoutDeletedAt
	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:00';unique_index:idx_user_username_deleted_at_unique"`
}
