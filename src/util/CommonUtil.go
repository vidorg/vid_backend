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

// Time UUID (20bit)
func (c *commonUtil) CurrentTimeUuid() string {
	return strings.Replace(time.Now().Format("20060102150405.000000"), ".", "", -1)
}

// Time UUID (14bit)
func (c *commonUtil) CurrentTimeString() string {
	return time.Now().Format("20060102150405")
}

// Check directory or file exist
func (c *commonUtil) IsDirOrFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Check exist and create it if not exist
func (c *commonUtil) CheckCreateDir(filename string) error {
	dir := path.Dir(filename)
	if !c.IsDirOrFileExist(dir) {
		return os.MkdirAll(dir, os.ModePerm) // 0777
	}
	return nil
}

// Get image / video filename from url: split prefix and query
func (c *commonUtil) GetFilenameFromUrl(url string, prefix string) (filename string, ok bool) {
	if len(url) <= len(prefix) {
		return url, false
	}
	filename = (url)[len(prefix):]
	idx := strings.Index(filename, "?")
	if idx != -1 {
		filename = filename[:idx]
	}
	return filename, true
}
