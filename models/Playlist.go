package models

import (
	"encoding/json"
	"strings"
	"time"

	"vid/config"
)

type Playlist struct {
	Gid         int       `json:"gid" gorm:"primary_key;auto_increment"`
	Groupname   string    `json:"groupname" gorm:"type:varchar(50);not_null"` // 50
	Description string    `json:"description"`
	CreateTime  time.Time `json:"create_time" gorm:"default:'2000-01-01'"`
	AuthorUid   int       `json:"-"`
	Author      *User     `json:"author" gorm:"-"`
	Videos      []*Video  `json:"videos" gorm:"-"`
}

func (p *Playlist) Equals(obj *Playlist) bool {
	return p.Gid == obj.Gid && p.Groupname == obj.Groupname && p.Description == obj.Description
}

func (p *Playlist) Unmarshal(jsonBody string, isNewPlaylist bool) bool {
	err := json.Unmarshal([]byte(jsonBody), p)
	if err != nil ||
		(!isNewPlaylist && p.Gid == 0) ||
		(isNewPlaylist && (p.Groupname == "")) {
		return false
	}
	// No description Field
	if strings.Index(jsonBody, "\"description\": \"") == -1 {
		p.Description = config.AppCfg.MagicToken
	}
	return true
}
