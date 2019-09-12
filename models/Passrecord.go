package models

type PassRecord struct {
	Uid           int    `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:char(48);not null"`
}

// Rename Table
func (PassRecord) TableName() string {
	// Default: pass_record
	return "tbl_passrecord"
}
