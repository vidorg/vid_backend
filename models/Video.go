package models

import (
	"encoding/json"
	"strings"
	"time"

	"vid/config"
)

type Video struct {
	Vid         int       `json:"vid" gorm:"primary_key;auto_increment"`
	Title       string    `json:"title" gorm:"type:varchar(100);not_null"` // 100
	Description string    `json:"description"`                             // 255
	VideoUrl    string    `json:"video_url" gorm:"not_null;unique"`        // 255
	UploadTime  time.Time `json:"upload_time" gorm:"type:datetime;default:'2000-01-01'"`
	AuthorUid   int       `json:"-"`
	Author      *User     `json:"author" gorm:"-"` // omitempty
}

func (v *Video) Equals(obj *Video) bool {
	return v.Vid == obj.Vid && v.Title == obj.Title && v.Description == obj.Description &&
		v.VideoUrl == obj.VideoUrl && v.AuthorUid == obj.AuthorUid
}

func (v *Video) Unmarshal(jsonBody string, isNewVideo bool) bool {
	err := json.Unmarshal([]byte(jsonBody), v)
	if err != nil ||
		(isNewVideo && (v.Title == "" || v.VideoUrl == "")) ||
		(!isNewVideo && (v.Vid == 0)) {
		return false
	}
	// No description Field
	if strings.Index(jsonBody, "\"description\": \"") == -1 {
		v.Description = config.AppCfg.MagicToken
	}
	return true
}
