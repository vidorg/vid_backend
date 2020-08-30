package param

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"strings"
)

var (
	ADPage  = goapidoc.NewPathParam("page", "integer#int32", true, "请求页")
	ADLimit = goapidoc.NewPathParam("limit", "integer#int32", true, "页大小")
	ADOrder = goapidoc.NewPathParam("order", "string", true, "排序")
)

type PageParam struct {
	Page  int32
	Limit int32
}

type PageOrderParam struct {
	Page  int32
	Limit int32
	Order string
}

// Bind ?page&limit
func BindPage(c *gin.Context, config *config.Config) *PageParam {
	page, err := xnumber.Atoi32(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := xnumber.Atoi32(c.DefaultQuery("limit", "0"))
	if def := config.Meta.DefPageSize; err != nil || limit <= 0 {
		limit = def
	} else if max := config.Meta.MaxPageSize; limit > max {
		limit = max
	}

	return &PageParam{Page: page, Limit: limit}
}

// Bind ?page&limit&order
func BindPageOrder(c *gin.Context, config *config.Config) *PageOrderParam {
	page := BindPage(c, config)
	order := c.DefaultQuery("order", "")
	return &PageOrderParam{Page: page.Page, Limit: page.Limit, Order: order}
}

// Bind :xid
func BindRouteId(c *gin.Context, field string) (uint64, error) {
	uid, err := xnumber.Atou64(c.Param(field))
	if err != nil {
		return 0, err
	}
	if uid <= 0 {
		return 0, fmt.Errorf("id shoule larger than 0")
	}

	return uid, nil
}

// Bind ?xxx=true|false
func BindQueryBool(c *gin.Context, field string) bool {
	val := c.DefaultQuery(field, "false")
	if strings.ToLower(val) == "true" {
		return true
	}
	return false
}
