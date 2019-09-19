package controllers

import (
	"net/http"
	"fmt"

	. "vid/exceptions"
	. "vid/utils"
	. "vid/models/resp"
	. "vid/models"

	"github.com/gin-gonic/gin"
)

type rawCtrl struct{}

var RawCtrl = new(rawCtrl)

// POST /raw/upload/img
func (r *rawCtrl) UploadImage(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := CmnUtil.ImageExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: FileExtException.Error(),
		})
		return
	}

	filename = CmnUtil.CurrentTimeInt()
	// ./files/image/2/20190919170528.jpg
	filepath := fmt.Sprintf("./files/image/%d/%s%s", uid, filename, ext)
	if !CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, Message{
			Message: ImageUploadException.Error(),
		})
	} else {
		// http://127.0.0.1:1234/raw/image/2/20190919170528.jpg
		c.JSON(http.StatusInternalServerError, RawResp{
			Type: "Image",
			Url: CmnUtil.GetServerUrl(fmt.Sprintf("raw/image/%d/%s%s", uid, filename, ext)),
		})
	}
}

// POST /raw/upload/video
func (r *rawCtrl) UploadVideo(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid

	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := CmnUtil.VideoExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: FileExtException.Error(),
		})
		return
	}
	
	filename = CmnUtil.CurrentTimeInt()
	filepath := fmt.Sprintf("./files/video/%d/%s%s", authusr.(User).Uid, filename, ext)
	if !CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, Message{
			Message: VideoUploadException.Error(),
		})
	} else {
		c.JSON(http.StatusInternalServerError, RawResp{
			Type: "Video",
			Url: CmnUtil.GetServerUrl(fmt.Sprintf("raw/video/%d/%s%s", uid, filename, ext)),
		})
	}
}

// GET /raw/image/:user/:filename
func (r *rawCtrl) RawImage(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "filename"),
		})
		return
	}
	
	filepath := fmt.Sprintf("./files/image/%d/%s", uid, filename)
	if !CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, Message{
			Message: FileNotExistException.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "image/png")
	c.File(filepath)
}

// GET /raw/video/:user/:filename
func (r *rawCtrl) RawVideo(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "filename"),
		})
		return
	}

	filepath := fmt.Sprintf("./files/video/%d/%s", uid, filename)
	if !CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, Message{
			Message: FileNotExistException.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "video/mpeg4")
	c.File(filepath)
}
