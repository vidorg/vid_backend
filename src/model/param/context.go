package param

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
)

var (
	ADPage  = goapidoc.NewPathParam("page", "integer#int32", true, "请求页")
	ADLimit = goapidoc.NewPathParam("limit", "integer#int32", true, "页大小")
	ADOrder = goapidoc.NewPathParam("order", "string", true, "排序")
)

func BindRouteId(c *gin.Context, field string) (int32, bool) {
	uid, err := xnumber.ParseInt32(c.Param(field), 10)
	if err != nil {
		return 0, false
	}
	if uid <= 0 {
		return 0, false // <<<
	}
	return uid, true
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
	page, err := xnumber.ParseInt32(c.DefaultQuery("page", "1"), 10)
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := xnumber.ParseInt32(c.DefaultQuery("limit", "0"), 10)
	if def := config.Meta.DefPageSize; err != nil || limit <= 0 {
		limit = def
	} else if max := config.Meta.MaxPageSize; limit > max {
		limit = max
	}
	return &PageParam{
		Page:  page,
		Limit: limit,
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
