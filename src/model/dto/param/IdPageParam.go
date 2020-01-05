package param

import (
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"strconv"
)

// make sure id is greater than 0
func BindRouteId(c *gin.Context, field string) (int, bool) {
	uid, err := strconv.Atoi(c.Param(field))
	if err != nil {
		return 0, false
	}
	if uid <= 0 {
		return 0, false // <<<
	} else {
		return uid, true
	}
}

// parse page
func BindQueryPage(c *gin.Context) (int, bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return 0, false
	}
	return xconditions.IfThenElse(page <= 0, 1, page).(int), true // <<<
}
