package po

import (
	"time"
)

// Subscribe: tbl_subscribe, po.User <-> po.Channel.
type Subscribe struct {
	Uid       uint64
	Cid       uint64
	CreatedAt time.Time
}
