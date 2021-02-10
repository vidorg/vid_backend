package model

type Video struct {
	BaseModel
	Title       string  `json:"title" json:"title,omitempty"`
	Description *string `json:"description" json:"description,omitempty"`
	VodID       *string `json:"vod_id" json:"vod_id,omitempty"`
	URL         *string `json:"url" json:"url,omitempty"`
	Cover       *string `json:"cover" json:"cover,omitempty"`
	UserID      int64   `json:"user_id" json:"user_id,omitempty"`
	Author      User    `gorm:"foreignKey:UserID" json:"author" json:"author"`
	CategoryID  int64   `json:"category_id,omitempty"`
	ChannelID   *int64  `json:"channel_id"`
}
