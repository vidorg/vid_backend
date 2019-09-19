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
func (r *rawCtrl) UploadImg(c *gin.Context) {
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

	if !CmnCtrl.IsImageExt(filename) {
		c.JSON(http.StatusBadRequest, Message{
			Message: FileExtException.Error(),
		})
		return
	}

	filepath := fmt.Sprintf("./files/image/%d/%s", uid, filename)
	if !CmnCtrl.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, Message{
			Message: ImageUploadException.Error(),
		})
	} else {
		c.JSON(http.StatusInternalServerError, RawResp{
			Type: "Image",
			Url: CmnCtrl.GetServerUrl(fmt.Sprintf("raw/img/%d/%s", uid, filename)),
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

	if !CmnCtrl.IsVideoExt(filename) {
		c.JSON(http.StatusBadRequest, Message{
			Message: FileExtException.Error(),
		})
		return
	}
	
	filepath := fmt.Sprintf("./files/video/%d/%s", authusr.(User).Uid, filename)
	if !CmnCtrl.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, Message{
			Message: VideoUploadException.Error(),
		})
	} else {
		c.JSON(http.StatusInternalServerError, RawResp{
			Type: "Video",
			Url: CmnCtrl.GetServerUrl(fmt.Sprintf("raw/video/%d/%s", uid, filename)),
		})
	}
}

// GET /raw/img/:user/:filename
func (r *rawCtrl) RawImg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Raw Img",
	})
}

// GET /raw/video/:user/:filename
func (r *rawCtrl) RawVideo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Raw Video",
	})
}
