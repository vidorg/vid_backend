package model

import (
	"github.com/vidorg/vid_backend/pkg/orm"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User user model
type User struct {
	ID       int     `gorm:"primaryKey;autoIncrement;comment:用户主键ID" json:"id"`
	UserName string  `gorm:"column:username;not null;unique;comment:用户名" json:"username"`
	Password string  `gorm:"not null;comment:用户密码" json:"password,omitempty"`
	Nickname string  `gorm:"not null;size:15;comment:用户昵称" json:"nickname"`
	Status   string  `gorm:"not null;default:active;comment:用户状态，active激活，incative未激活suspend被封禁" json:"status"`
	Avatar   string  `gorm:"size:1000;default:https://oss.seefs.cn/avatar.jpg;comment:用户头像" json:"avatar"`
	Email    *string `gorm:"column:email;comment:用户Email" json:"email"`
	Role     string  `gorm:"size:10;not null;comment:用户权限" json:"role"`
	Amount   int     `json:"amount"`
	//UpdatedAt int            `json:"updated_at"`
	//Created   int64          `gorm:"autoCreateTime" json:"created"`
	CreatedAt time.Time `json:"created_time"`
	UpdatedAt time.Time `json:"updated_time"`
	//DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_time,omitempty"`
}

const (
	PasswordCost = bcrypt.MinCost // password cost
	Active       = "active"       // active user
	Inactive     = "inactive"     // inactive user
	Suspend      = "suspend"      // banned user
)

// GetUser Get user by ID (for middleware GetUser)
func GetUser(ID interface{}) (*User, error) {
	user := &User{}
	rdb := orm.DB().First(user, ID)
	return user, rdb.Error
}

// SetPassword set user password
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// MatchPassword match password
func (user *User) MatchPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) // encrypted, unencrypted
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return false, err
	}
}
