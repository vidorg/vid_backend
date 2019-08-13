package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	. "vid/exceptions"
	. "vid/models"

	"github.com/gin-gonic/gin"
)

type ReqUtil struct{}

func (b *ReqUtil) GetBody(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func (b *ReqUtil) CheckJsonValid(bodyJson string, objPtr BaseModel, key string) {
	err := json.Unmarshal([]byte(bodyJson), objPtr)
	if err != nil || !(objPtr).CheckValid() {
		var param = make([]string, 1, 1)
		param = append(param, key)
		panic(NewParamError(param, false))
	}
}

func (b *ReqUtil) GetIntParam(param gin.Params, key string) int {
	valStr, ok := param.Get(key)

	val, err := strconv.Atoi(valStr)

	if !ok || err != nil {
		var param = make([]string, 1, 1)
		param = append(param, key)
		panic(NewParamError(param, true))
	} else {
		return val
	}
}

func (b *ReqUtil) GetStrParam(param gin.Params, key string) string {
	val, ok := param.Get(key)
	if !ok {
		var param = make([]string, 1, 1)
		param = append(param, key)
		panic(NewParamError(param, true))
	} else {
		return val
	}
}
