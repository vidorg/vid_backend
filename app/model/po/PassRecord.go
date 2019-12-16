package po

type PassRecord struct {
	Uid           int    `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:char(48);not null"`

	User *User `gorm:"foreignkey:Uid"`

	TimePo
}

func (PassRecord) TableName() string {
	return "tbl_password"
}
