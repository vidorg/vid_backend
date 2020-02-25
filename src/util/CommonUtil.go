package util

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"os"
	"path"
	"strings"
)

type commonUtil struct{}

var CommonUtil = new(commonUtil)

// Check directory or file exist
func (c *commonUtil) IsDirOrFileExist(filename string) bool {
	return xcondition.Second(os.Stat(filename)) == nil
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
func (c *commonUtil) GetFilenameFromUrl(url string, prefix string) string {
	if url == "" {
		return ""
	}
	if len(url) <= len(prefix) {
		return ""
	}
	filename := url[len(prefix):]
	idx := strings.Index(filename, "?")
	if idx != -1 {
		filename = filename[:idx]
	}
	return filename
}
