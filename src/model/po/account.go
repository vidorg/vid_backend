package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
)

type Account struct {
	User     *User  `gorm:"foreignkey:Uid; association_foreignkey:Uid"` // po.Account belongs to po.User
	Uid      uint64 `gorm:"primary_key;       not null"`                // user id (foreigner key)
	Password string `gorm:"type:varchar(255); not null"`                // encrypted password

	xgorm.GormTime
}
