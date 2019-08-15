package models

type Passrecord struct {
	// User     User   `gorm:"foreignkey:uid;association_foreignkey:uid"`
	Uid      int    `gorm:"primary_key"`
	HashPass string `gorm:"type:char(128);not null"`
}
