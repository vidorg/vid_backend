package utils

import (
	"strings"
)

type CmnCtrl struct{}

// 字符串首字母大写
func (c *CmnCtrl) Capitalize(str string) string {
	return strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
}
