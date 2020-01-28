package po

import "github.com/vidorg/vid_backend/src/common/model"

type Account struct {
	Uid           int32  `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:varchar(255);not null"`

	User *User `gorm:"foreignkey:Uid"`

	model.GormTime
}
