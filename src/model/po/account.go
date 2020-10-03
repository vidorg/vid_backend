package po

import (
	"github.com/vidorg/vid_backend/lib/xgorm"
)

// Account: tbl_account.
type Account struct {
	Uid      uint64 `gorm:"                   not null; primaryKey"` // user id (foreigner key)
	Password string `gorm:"type:varchar(255); not null"`             // encrypted password

	User *User `gorm:"foreignKey:Uid; references:Uid"` // po.Account belongs to po.User

	xgorm.Model
}
