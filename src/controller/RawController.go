package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
	"strings"
)

type rawController struct {
	config *config.ServerConfig
}

func RawController(config *config.ServerConfig) *rawController {
	return &rawController{
		config: config,
	}
}

// @Router				/v1/raw/image/{filename} [GET]
// @Summary				获取图片
// @Description			获取用户头像图片以及视频封面
// @Tag					Raw
// @Param				filename path string true "图片文件名，jpg后缀名"
// @Accept				multipart/form-data
// @ErrorCode			404 image not found
/* @Success 200			{ "Content-Type": "image/jpeg" } */
func (r *rawController) RawImage(c *gin.Context) {
	filename := c.Param("filename")
	isDefault := strings.Index(filename, "_") == -1

	filePath := xconditions.IfThenElse(isDefault, fmt.Sprintf("./usr/default/%s", filename), fmt.Sprintf("./usr/image/%s", filename)).(string)
	if !util.CommonUtil.IsDirOrFileExist(filePath) {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotFoundError.Error()).JSON(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}
