package utils

import (
	"strings"
)

type cmnCtrl struct{}

var CmnCtrl = new(cmnCtrl)

// 字符串首字母大写
func (c *cmnCtrl) Capitalize(str string) string {
	return strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
}
