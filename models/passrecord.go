package models

type Passrecord struct {
	Uid           int    `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:char(128);not null"`
}

// @override
func (p *Passrecord) CheckValid() bool {
	return p.Uid != 0 && p.EncryptedPass != ""
}
