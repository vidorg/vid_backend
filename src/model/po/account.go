package po

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
)

type Account struct {
	Uid           int32  `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:varchar(255);not null"`

	User *User `gorm:"foreignkey:Uid"`

	xgorm.GormTime
}
