package controller

import (
	"fmt"
	"net/http"
	"vid/app/controller/exception"
	po2 "vid/app/model/po"
	"vid/app/model/resp"
	"vid/app/util"

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
			Message: exception.JsonParamError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := util.CmnUtil.ImageExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exception.FileExtensionError.Error(),
		})
		return
	}

	filename = util.CmnUtil.CurrentTimeInt()
	// ./files/image/2/20190919170528.jpg
	filepath := fmt.Sprintf("./files/image/%d/%s%s", uid, filename, ext)
	if !util.CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: exception.ImageUploadError.Error(),
		})
	} else {
		// http://127.0.0.1:1234/raw/image/2/20190919170528.jpg
		c.JSON(http.StatusInternalServerError, resp.RawResp{
			Type: "Image",
			Url:  util.CmnUtil.GetImageUrl(uid, filename+ext),
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
			Message: exception.JsonParamError.Error(),
		})
		return
	}
	filename := header.Filename

	ok, ext := util.CmnUtil.VideoExt(filename)
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exception.FileExtensionError.Error(),
		})
		return
	}

	filename = util.CmnUtil.CurrentTimeInt()
	filepath := fmt.Sprintf("./files/video/%d/%s%s", authusr.(po2.User).Uid, filename, ext)
	if !util.CmnUtil.SaveFile(filepath, file) {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: exception.VideoUploadError.Error(),
		})
	} else {
		c.JSON(http.StatusInternalServerError, resp.RawResp{
			Type: "Video",
			Url:  util.CmnUtil.GetVideoUrl(uid, filename+ext),
		})
	}
}

// GET /raw/image/:user/:filename (Non-Admin)
func (r *rawCtrl) RawImage(c *gin.Context) {
	uid, ok := util.ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exception.RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := util.ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exception.RouteParamError.Error(), "filename"),
		})
		return
	}

	var filepath string
	if uid == -1 {
		filepath = fmt.Sprintf("./files/image/default/%s", filename)
	} else {
		filepath = fmt.Sprintf("./files/image/%d/%s", uid, filename)
	}

	if !util.CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exception.FileNotFoundError.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "image/png")
	c.File(filepath)
}

// GET /raw/video/:user/:filename (Non-Admin)
func (r *rawCtrl) RawVideo(c *gin.Context) {
	uid, ok := util.ReqUtil.GetIntParam(c.Params, "user")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exception.RouteParamError.Error(), "user"),
		})
		return
	}
	filename, ok := util.ReqUtil.GetStrParam(c.Params, "filename")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exception.RouteParamError.Error(), "filename"),
		})
		return
	}

	filepath := fmt.Sprintf("./files/video/%d/%s", uid, filename)
	if !util.CmnUtil.IsFileExist(filepath) {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exception.FileNotFoundError.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "video/mpeg4")
	c.File(filepath)
}
