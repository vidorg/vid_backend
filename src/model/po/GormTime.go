package po

import "time"

const (
	DefaultDeleteAtTimeStamp = "2000-01-01 00:00:01"
)

type GormTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"default:'2000-01-01 00:00:01'"`
}

type GormTimeWithoutDeletedAt struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
