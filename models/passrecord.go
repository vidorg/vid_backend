package models

type Passrecord struct {
	Uid      int    `gorm:"primary_key"`
	HashPass string `gorm:"type:char(128);not null"`
}
