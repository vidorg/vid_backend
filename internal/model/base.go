package model

import "gorm.io/gorm"

type BaseModel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UpdatedAt int64 `json:"updated_at"`
	Created   int64 `gorm:"autoCreateTime" json:"created"`
	//CreatedAt time.Time `json:"created_time"`
	//UpdatedAt time.Time `json:"updated_time"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_time,omitempty"`
}
