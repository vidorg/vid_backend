package model

type Channel struct {
	BaseModel
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Videos      []Video `json:"videos"`
	Users       []User  `gorm:"many2many:channel_author;"`
	Subscribers []User  `gorm:"many2many:channel_subscriber;"`
}
