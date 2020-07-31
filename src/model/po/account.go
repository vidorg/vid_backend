package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
)

type Account struct {
	Uid           int32  `gorm:"primary_key"`
	User          *User  `gorm:"foreignkey:To"`
	EncryptedPass string `gorm:"type:varchar(255);not null"`

	xgorm.GormTime
}
