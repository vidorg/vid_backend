package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/vidorg/vid_backend/src/model/constant"
	"time"
)

type User struct {
	Uid      uint64             `gorm:"primary_key; auto_increment"`                           // user id
	Username string             `gorm:"type:varchar(127); not null; unique_index:uk_username"` // username, unique
	Email    string             `gorm:"type:varchar(255); not null; unique_index:uk_email"`    // user email, unique
	Nickname string             `gorm:"type:varchar(127); not null"`                           // user nickname
	Gender   constant.Gender    `gorm:"type:tinyint;      not null; default:0"`                // user gender (0X, 1M, 2D)
	Profile  string             `gorm:"type:varchar(255); not null"`                           // user profile, allowempty
	Avatar   string             `gorm:"type:varchar(255); not null"`                           // user avatar url, using oss
	Birthday xtime.JsonDate     `gorm:"type:date;         not null; default:'2000-01-01'"`     // user birthday
	Role     string             `gorm:"type:varchar(255); not null; default:'normal'"`         // user role, used in casbin
	State    constant.UserState `gorm:"type:tinyint;      not null; default:0"`                // user state (0|1|2)
	Phone    string             `gorm:"type:varchar(127); not null"`                           // user phone number

	// tbl_subscribe
	Subscribings []*User `gorm:"many2many:subscribe; jointable_foreignkey:subscriber_uid; association_jointable_foreignkey:up_uid"`         // up_uid -> subscriber_uid
	Subscribers  []*User `gorm:"many2many:subscribe; jointable_foreignkey:up_uid;         association_jointable_foreignkey:subscriber_uid"` // subscriber_uid -> up_uid

	xgorm.GormCUTime
	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:00'; unique_index:uk_username,uk_email"`
}
