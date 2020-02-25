package param

import (
	"github.com/gin-gonic/gin"
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

// parse page
func BindQueryPage(c *gin.Context) int32 {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return 1
	}
	return int32(page)
}

// parse order
func BindQueryOrder(c *gin.Context) string {
	return c.DefaultQuery("order", "")
}
