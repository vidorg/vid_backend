package util

import (
	"fmt"
	"io"
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
		err := os.MkdirAll(dir, os.ModePerm)
		return err == nil
	}
	return true
}

// 获得服务器根网址
//
//  `str` `xx/xx/`
//
// @return `http://xx:xx/`
func (c *commonUtil) GetServerUrl(str string) string {
	// return fmt.Sprintf("http://%s:%d/%s", config.AppConfig.HTTPServer, config.AppConfig.HTTPPort, str)
	return ""
}

// 获得默认头像
func (c *commonUtil) GetDefaultAvatarUrl() string {
	return c.GetServerUrl("raw/image/-1/avatar.jpg")
}

// 获得默认视频封面
func (c *commonUtil) GetDefaultVideoCoverUrl() string {
	return c.GetServerUrl("raw/image/-1/videocover.jpg")
}

// 获得图片地址
func (c *commonUtil) GetImageUrl(uid int, resource string) string {
	return c.GetServerUrl(fmt.Sprintf("raw/image/%d/%s", uid, resource))
}

// 获得视频地址
func (c *commonUtil) GetVideoUrl(uid int, resource string) string {
	return c.GetServerUrl(fmt.Sprintf("raw/video/%d/%s", uid, resource))
}

// ////////////////////////////////////////////////////////////

// 判断是否是受支持视频格式
// mp4
//
// @return `ok` `ext`
func (c *commonUtil) VideoExt(filename string) (bool, string) {
	ext := path.Ext(filename)
	return ext == ".mp4",
		ext
}

// 保存文件，并且覆盖已存在文件
func (c *commonUtil) SaveFile(filename string, file io.Reader) bool {

	// File path
	dir := path.Dir(filename)
	if !c.IsFileOrDirExist(dir) {
		os.MkdirAll(dir, 0777)
		if !c.IsFileOrDirExist(dir) {
			return false
		}
	}

	// Delete file
	if c.IsFileOrDirExist(filename) {
		err := os.Remove(filename)
		if err != nil {
			return false
		}
	}

	// Create file
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return false
	}

	// File exist
	if !c.IsFileOrDirExist(filename) {
		return false
	}

	// Save file
	_, err = io.Copy(f, file)
	if err != nil {
		return false
	}

	return true
}
