package util

import (
	"os"
	"path"
	"strings"
	"time"
)

type commonUtil struct{}

var CommonUtil = new(commonUtil)

func (c *commonUtil) Capitalize(str string) string {
	return strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
}

// File UUID (20bit)
func (c *commonUtil) CurrentTimeUuid() string {
	return strings.Replace(time.Now().Format("20060102150405.000000"), ".", "", -1)
}

func (c *commonUtil) IsFileOrDirExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (c *commonUtil) CheckCreateDir(filename string) bool {
	dir := path.Dir(filename)
	if !c.IsFileOrDirExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm) // 0777
		return err == nil
	}
	return true
}
