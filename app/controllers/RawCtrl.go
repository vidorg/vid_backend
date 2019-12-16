package controllers

import (
	"fmt"
	"net/http"
	"vid/app/controllers/exceptions"
	po2 "vid/app/models/po"
	"vid/app/models/resp"
	"vid/app/utils"

	"github.com/gin-gonic/gin"
)

type rawCtrl struct{}

var RawCtrl = new(rawCtrl)

// POST /raw/upload/img (Admin)
func (r *rawCtrl) UploadImage(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(po2.User).Uid

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.RequestBodyError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := utils.CmnUtil.ImageExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.FileExtException.Error(),
		})
		return
	}

	filename = utils.CmnUtil.CurrentTimeInt()
	// ./files/image/2/20190919170528.jpg
	filepath := fmt.Sprintf("./files/image/%d/%s%s", uid, filename, ext)
	if !utils.CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: exceptions.ImageUploadException.Error(),
		})
	} else {
		// http://127.0.0.1:1234/raw/image/2/20190919170528.jpg
		c.JSON(http.StatusInternalServerError, resp.RawResp{
			Type: "Image",
			Url:  utils.CmnUtil.GetImageUrl(uid, filename+ext),
		})
	}
}

// POST /raw/upload/video (Admin)
func (r *rawCtrl) UploadVideo(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(po2.User).Uid

	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.RequestBodyError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := utils.CmnUtil.VideoExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.FileExtException.Error(),
		})
		return
	}

	filename = utils.CmnUtil.CurrentTimeInt()
	filepath := fmt.Sprintf("./files/video/%d/%s%s", authusr.(po2.User).Uid, filename, ext)
	if !utils.CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: exceptions.VideoUploadException.Error(),
		})
	} else {
		c.JSON(http.StatusInternalServerError, resp.RawResp{
			Type: "Video",
			Url:  utils.CmnUtil.GetVideoUrl(uid, filename+ext),
		})
	}
}

// GET /raw/image/:user/:filename (Non-Admin)
func (r *rawCtrl) RawImage(c *gin.Context) {
	uid, ok := utils.ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := utils.ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "filename"),
		})
		return
	}

	var filepath string
	if uid == -1 {
		filepath = fmt.Sprintf("./files/image/default/%s", filename)
	} else {
		filepath = fmt.Sprintf("./files/image/%d/%s", uid, filename)
	}

	if !utils.CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exceptions.FileNotExistException.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "image/png")
	c.File(filepath)
}

// GET /raw/video/:user/:filename (Non-Admin)
func (r *rawCtrl) RawVideo(c *gin.Context) {
	uid, ok := utils.ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := utils.ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "filename"),
		})
		return
	}

	filepath := fmt.Sprintf("./files/video/%d/%s", uid, filename)
	if !utils.CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exceptions.FileNotExistException.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "video/mpeg4")
	c.File(filepath)
}
