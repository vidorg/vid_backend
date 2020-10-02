package po

import (
	"time"
)

type Favorite struct {
	Uid       uint64
	Vid       uint64
	CreatedAt time.Time
}
