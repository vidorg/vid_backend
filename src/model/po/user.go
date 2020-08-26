package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/vidorg/vid_backend/src/model/constant"
	"time"
)

type User struct {
	Uid         int32            `gorm:"primary_key;auto_increment"`
	Username    string           `gorm:"not_null;type:varchar(30);unique_index:idx_user_username_deleted_at_unique"` // 30
	Gender      constant.SexEnum `gorm:"not_null;type:tinyint;default:0"`                                            // 0X 1M 2F
	Profile     string           `gorm:"type:varchar(255)"`                                                          // 255
	AvatarUrl   string           `gorm:"not_null;type:varchar(255)"`                                                 // 255 // TODO url
	Birthday    xtime.JsonDate   `gorm:"not_null;type:date;default:'2000-01-01'"`                                    // RFC3339
	Role        string           `gorm:"not_null;type:varchar(10);default:'normal'"`                                 // normal / admin / root
	RegisterIP  string           `gorm:"type:varchar(15)"`                                                           // 15
	PhoneNumber string           `gorm:"type:varchar(11)"`                                                           // 11

	// tbl_subscribe
	Subscribings []*User `gorm:"many2many:subscribe;jointable_foreignkey:subscriber_uid;association_jointable_foreignkey:up_uid"` // up_uid -> subscriber_uid
	Subscribers  []*User `gorm:"many2many:subscribe;jointable_foreignkey:up_uid;association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid

	xgorm.GormTimeWithoutDeletedAt
	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:00';unique_index:idx_user_username_deleted_at_unique"`
}
