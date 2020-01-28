package param

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
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
func BindQueryPage(c *gin.Context) (int32, bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return 0, false
	}
	return xcondition.IfThenElse(page <= 0, 1, int32(page)).(int32), true // <<<
}
