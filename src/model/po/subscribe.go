package po

import (
	"time"
)

type Subscribe struct {
	FromUid   uint64
	ToUid     uint64
	CreatedAt time.Time
}
