package po

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/model/constant"
)

// User: tbl_user.
type User struct {
	Uid      uint64             `gorm:"                   not null; primaryKey; autoIncrement"` // user id
	Username string             `gorm:"type:varchar(127); not null; uniqueIndex:uk_username"`   // username, unique
	Email    string             `gorm:"type:varchar(255); not null; uniqueIndex:uk_email"`      // user email, unique
	Nickname string             `gorm:"type:varchar(127); not null"`                            // user nickname
	Gender   constant.Gender    `gorm:"type:tinyint;      not null; default:0"`                 // user gender (0X, 1M, 2D)
	Profile  string             `gorm:"type:varchar(255); not null"`                            // user profile, allowempty
	Avatar   string             `gorm:"type:varchar(255); not null"`                            // user avatar url, using oss
	Birthday xtime.JsonDate     `gorm:"type:date;         not null; default:'2000-01-01'"`      // user birthday
	Role     string             `gorm:"type:varchar(255); not null; default:'normal'"`          // user role, used in casbin
	State    constant.UserState `gorm:"type:tinyint;      not null; default:0"`                 // user state (0|1|2)

	Followings   []*User    `gorm:"many2many:tbl_follow;    foreignKey:Uid; references:Uid; joinForeignKey:from_uid; joinReferences:to_uid"`   // tbl_follow
	Followers    []*User    `gorm:"many2many:tbl_follow;    foreignKey:Uid; references:Uid; joinForeignKey:to_uid;   joinReferences:from_uid"` // tbl_follow
	Subscribings []*Channel `gorm:"many2many:tbl_subscribe; foreignKey:Uid; references:Cid; joinForeignKey:uid;      joinReferences:cid"`      // tbl_subscribe
	Favorites    []*Video   `gorm:"many2many:tbl_favorite;  foreignKey:Uid; references:Vid; joinForeignKey:uid;      joinReferences:vid"`      // tbl_favorite

	xgorm.Model
}
