package param

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"strconv"
)

// make sure id is greater than 0
func BindRouteId(c *gin.Context, field string) (int32, bool) {
	uid, err := strconv.Atoi(c.Param(field))
	if err != nil {
		return 0, false
	}
	if uid <= 0 {
		return 0, false // <<<
	} else {
		return int32(uid), true
	}
}

type PageParam struct {
	Page  int32
	Limit int32
}

type PageOrderParam struct {
	*PageParam
	Order string
}

func BindPage(c *gin.Context, config *config.Config) *PageParam {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil {
		if limit <= 0 {
			limit = int(config.Meta.DefPageSize)
		} else if limit > int(config.Meta.MaxPageSize) {
			limit = int(config.Meta.MaxPageSize)
		}
	}
	return &PageParam{
		Page:  int32(page),
		Limit: int32(limit),
	}
}

func BindPageOrder(c *gin.Context, config *config.Config) *PageOrderParam {
	page := BindPage(c, config)
	order := c.DefaultQuery("order", "")
	return &PageOrderParam{
		PageParam: page,
		Order:     order,
	}
}
