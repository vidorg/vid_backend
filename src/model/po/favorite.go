package po

import (
	"time"
)

// Favorite: tbl_Favorite, po.User <-> po.Video.
type Favorite struct {
	Uid       uint64
	Vid       uint64
	CreatedAt time.Time
}
