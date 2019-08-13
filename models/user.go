package models

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username,omitempty" gorm:"not null"`
}

func (u *User) CheckValid() bool {
	return u.ID != 0 && u.Username != ""
}
