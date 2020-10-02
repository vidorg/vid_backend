package po

import (
	"github.com/vidorg/vid_backend/lib/xgorm"
)

type Channel struct {
	Cid         uint64
	Name        string
	Description string
	AuthorUid   uint64

	Author *User

	xgorm.Model
}
