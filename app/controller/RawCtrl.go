package controller

import (
	"fmt"
	"github.com/shomali11/util/xconditions"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/model/dto/common"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type rawCtrl struct{}

var RawCtrl = new(rawCtrl)

// @Router 				/raw/image/{uid}/{filename} [GET]
// @Summary 			获取图片
// @Description 		获取用户头像图片以及视频封面
// @Param 				uid path string true "用户id，或者default"
// @Param 				filename path string true "图片文件名，jpg后缀名"
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
			c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
			return
		}
	}
	filename := c.Param("filename")

	filePath := xconditions.IfThenElse(uid == -1, fmt.Sprintf("./usr/default/%s", filename), fmt.Sprintf("./usr/image/%d/%s", uid, filename)).(string)
	if !util.CommonUtil.IsFileOrDirExist(filePath) {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotFoundError.Error()))
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}
