package po

type PassRecord struct {
	TimePo

	EncryptedPass string `gorm:"type:char(48);not null"`
	User          *User  `gorm:"foreignkey:Uid"`
	Uid           int    `gorm:"primary_key"`
}

func (PassRecord) TableName() string {
	return "tbl_password"
}
