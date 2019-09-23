package models

type Videoinlist struct {
	Gid int `json:"-" gorm:"primary_key;auto_increment:false"`
	Vid int `json:"-" gorm:"primary_key;auto_increment:false"`
}