package model

type Channel struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Users       []User `gorm:"many2many:channel_author;"`
	Subscriber  []User `gorm:"many2many:channel_subscriber;"`
}
