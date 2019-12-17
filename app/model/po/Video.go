package po

import (
	"encoding/json"
	"strings"
	"time"
	"vid/app/util"

	// "vid/app/config"
)

type Video struct {
	Vid         int       `json:"vid" gorm:"primary_key;auto_increment"`
	Title       string    `json:"title" gorm:"type:varchar(100);not_null"` // 100
	Description string    `json:"description"`                             // 255
	VideoUrl    string    `json:"video_url" gorm:"not_null;unique"`        // 255
	CoverUrl    string    `json:"cover_url"`
	UploadTime  time.Time `json:"upload_time" gorm:"type:datetime;default:'2000-01-01'"`
	AuthorUid   int       `json:"-"`
	Author      *User     `json:"author" gorm:"-"` // omitempty

	GormTime
}

func (v *Video) Equals(obj *Video) bool {
	return v.Vid == obj.Vid && v.Title == obj.Title && v.Description == obj.Description &&
		v.VideoUrl == obj.VideoUrl && v.AuthorUid == obj.AuthorUid && v.CoverUrl == obj.CoverUrl
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
		// v.Description = config.AppConfig.MagicToken
	}
	return true
}

// Server -> DB
func (v *Video) ToDB() {
	if strings.HasPrefix(v.CoverUrl, "http") {
		sp := strings.Split(v.CoverUrl, "/")
		if sp[len(sp)-2] == "-1" {
			// http://127.0.0.1:3344/raw/image/-1/avatar.jpg -> /avatar.jpg
			v.CoverUrl = "/" + sp[len(sp)-1]
		} else {
			// http://127.0.0.1:3344/raw/image/233/avatar.jpg -> avatar.jpg
			v.CoverUrl = sp[len(sp)-1]
		}
	}
}

// DB -> Server
func (v *Video) ToServer() {
	if !strings.HasPrefix(v.CoverUrl, "http") {
		if strings.Index(v.CoverUrl, "/") != -1 {
			// /avatar.jpg -> http://127.0.0.1:3344/raw/image/-1/avatar.jpg
			v.CoverUrl = util.CmnUtil.GetImageUrl(-1, strings.Trim(v.CoverUrl, "/"))
		} else {
			// avatar.jpg -> http://127.0.0.1:3344/raw/image/233/avatar.jpg
			v.CoverUrl = util.CmnUtil.GetImageUrl(v.AuthorUid, v.CoverUrl)
		}
	}
}
