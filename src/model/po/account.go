package po

import (
	"gorm.io/gorm"
)

type Account struct {
	Uid      uint64 `gorm:"type:bigint;       not null; primaryKey"` // user id (foreigner key)
	Password string `gorm:"type:varchar(255); not null"`             // encrypted password

	User *User `gorm:"foreignKey:Uid; references:Uid"` // po.Account belongs to po.User

	gorm.Model
}
