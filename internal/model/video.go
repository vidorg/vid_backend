package model

type Video struct {
	BaseModel
	Title       string  `json:"title"`
	Description *string `json:"description"`
	VodID       *string `json:"vod_id"`
	URL         *string `json:"url"`
	Cover       *string `json:"cover"`
	UserID      int64   `json:"user_id"`
	Author      User    `gorm:"foreignKey:UserID" json:"author"`
}
