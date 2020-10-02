package po

import (
	"time"
)

type Follow struct {
	FromUid   uint64
	ToUid     uint64
	CreatedAt time.Time
}
