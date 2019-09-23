package models

type VideoInList struct {
	Gid int `json:"-" gorm:"primary_key;auto_increment:false"`
	Vid int `json:"-" gorm:"primary_key;auto_increment:false"`
}

func (VideoInList) TableName() string {
	return "tbl_videoinlist"
}
