package po

import (
	"time"
)

// Follow: tbl_follow, po.User <->  po.User.
type Follow struct {
	FromUid   uint64
	ToUid     uint64
	CreatedAt time.Time
}
