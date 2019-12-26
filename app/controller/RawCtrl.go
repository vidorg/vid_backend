package controller

import (
	"fmt"
	"github.com/shomali11/util/xconditions"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/model/dto"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type rawCtrl struct{}

var RawCtrl = new(rawCtrl)

// @Router 				/raw/image/{uid}/{filename} [GET]
// @Summary 			获取图片
// @Description 		获取用户头像图片
// @Param 				uid path string true "用户id，或者default"
// @Param 				filename path string true "图片文件名"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request route param error
// @ErrorCode 			404 image not found
/* @Success 200 		| Key | Value |
						| --- | --- |
 						| Content-Type | image/jpeg | */
func (r *rawCtrl) RawImage(c *gin.Context) {
	uidString := c.Param("uid")
	uid := -1
	if uidString != "default" {
		var err error
		uid, err = strconv.Atoi(uidString)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
			return
		}
	}
	filename := c.Param("filename")

	filePath := xconditions.IfThenElse(uid == -1, fmt.Sprintf("./usr/image/default/%s", filename), fmt.Sprintf("./usr/image/%d/%s", uid, filename)).(string)
	if !util.CommonUtil.IsFileOrDirExist(filePath) {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotFoundError.Error()))
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}

// // GET /raw/video/:user/:filename (Non-Admin)
// func (r *rawCtrl) RawVideo(c *gin.Context) {
// 	uid, ok := util.ReqUtil.GetIntParam(c.Params, "user")
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: fmt.Sprintf(exception.RouteParamError.Error(), "user"),
// 		})
// 		return
// 	}
// 	filename, ok := util.ReqUtil.GetStrParam(c.Params, "filename")
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: fmt.Sprintf(exception.RouteParamError.Error(), "filename"),
// 		})
// 		return
// 	}
//
// 	filepath := fmt.Sprintf("./usr/video/%d/%s", uid, filename)
// 	if !util.CommonUtil.IsFileOrDirExist(filepath) {
// 		c.JSON(http.StatusNotFound, resp.Message{
// 			Message: exception.FileNotFoundError.Error(),
// 		})
// 		return
// 	}
// 	c.Writer.Header().Add("Content-Type", "video/mpeg4")
// 	c.File(filepath)
// }
