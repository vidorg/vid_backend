package po

type PassRecord struct {
	Uid           int32  `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:varchar(255);not null"`

	User *User `gorm:"foreignkey:Uid"`

	GormTime
}

func (PassRecord) TableName() string {
	return "tbl_password"
}
