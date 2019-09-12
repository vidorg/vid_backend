package models

import (
	"encoding/json"
	"time"
)

type Video struct {
	Vid         int       `json:"vid" gorm:"primary_key;AUTO_INCREMENT"`
	Title       string    `json:"title" gorm:"type:varchar(100);not_null"` // 100
	Description string    `json:"description"`                             // 255
	VideoUrl    string    `json:"video_url"`                               // 255
	UploadTime  time.Time `json:"upload_time" gorm:"type:datetime;default:'2000-01-01'"`
	AuthorUid   int       `json:"-"`
	Author      *User     `json:"author" gorm:"-"` // omitempty
}

func (v *Video) Equals(obj *Video) bool {
	return v.Vid == obj.Vid && v.Title == obj.Title && v.Description == obj.Description &&
		v.VideoUrl == obj.VideoUrl && v.AuthorUid == obj.AuthorUid
}

func (v *Video) Unmarshal(jsonBody string, needVid bool, needUrl bool) bool {
	err := json.Unmarshal([]byte(jsonBody), v)
	if err != nil ||
		(needVid && v.Vid == 0) ||
		(needUrl && v.VideoUrl == "") {
		return false
	}
	return true
}
