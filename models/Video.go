package models

import "time"

type Video struct {
	Vid         int       `json:"vid" gorm:"primary_key;AUTO_INCREMENT"`
	Title       string    `json:"title" gorm:"type:varchar(80);not_null"`
	Description string    `json:"description" gorm:"type:varchar(250)"`
	VideoUrl    string    `json:"video_url"`
	Author      *User     `json:"author,omitempty" gorm:"-"`
	AuthorUid   int       `json:"-"`
	UploadTime  time.Time `json:"upload_time"`
}
