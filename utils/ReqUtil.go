package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	. "vid/models"

	"github.com/gin-gonic/gin"
)

type reqUtil struct{}

var ReqUtil = new(reqUtil)

// 获得请求体
func (b *reqUtil) GetBody(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

// 检查请求体的 Json 是否合法，需要实现 IBaseModel 接口
func (b *reqUtil) CheckJsonValid(bodyJson string, objPtr IBaseModel) bool {
	err := json.Unmarshal([]byte(bodyJson), objPtr)
	if err != nil || !objPtr.CheckValid() {
		return false
	} else {
		return true
	}
}

// 获得 int 类型的路由参数
//
// @return (`int`, `ok`)
func (b *reqUtil) GetIntParam(param gin.Params, key string) (int, bool) {
	valStr, ok := param.Get(key)
	val, err := strconv.Atoi(valStr)
	if !ok || err != nil {
		return -1, false
	} else {
		return val, true
	}
}

// 获得 str 类型的路由参数
//
// @return (`str`, `ok`)
func (b *reqUtil) GetStrParam(param gin.Params, key string) (string, bool) {
	val, ok := param.Get(key)
	if !ok || val == "" {
		return "", false
	} else {
		return val, true
	}
}

// 获得 int 类型的查询参数
//
// @return (`int`, `ok`)
func (b *reqUtil) GetIntQuery(c *gin.Context, key string) (int, bool) {
	valStr, ok := c.GetQuery(key)
	val, err := strconv.Atoi(valStr)
	if !ok || valStr == "" || err != nil {
		return -1, false
	} else {
		return val, true
	}
}

// 获得 str 类型的查询参数
//
// @return (`str`, `ok`)
func (b *reqUtil) GetStrQuery(c *gin.Context, key string) (string, bool) {
	val, ok := c.GetQuery(key)
	if !ok || val == "" {
		return "", false
	} else {
		return val, true
	}
}
