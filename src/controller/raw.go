package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/util"
)

type RawController struct {
	Config *config.Config `di:"~"`
	Logger *logrus.Logger `di:"~"`
}

func NewRawController(dic *xdi.DiContainer) *RawController {
	ctrl := &RawController{}
	dic.MustInject(ctrl)
	return ctrl
}

// @Router              /v1/raw/image [POST]
// @Summary             上传图片
// @Tag                 Raw
// @Security            Jwt
// @Param               image formData file true "上传的图片，大小限制在2M，允许后缀名为 {.jpg, .jpeg, .png, .bmp, .gif}"
// @ResponseModel 200   #Result<ImageDto>
func (r *RawController) UploadImage(c *gin.Context) {
	imageFile, imageHeader, err := c.Request.FormFile("image")
	if err != nil || imageFile == nil {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	supported, ext := util.ImageUtil.CheckImageExt(imageHeader.Filename)
	if !supported {
		result.Error(exception.ImageNotSupportedError).JSON(c)
		return
	}

	filename := fmt.Sprintf("%s.jpg", xstring.CurrentTimeUuid(20))
	savePath := fmt.Sprintf("%s%s", r.Config.File.ImagePath, filename)
	if err := util.ImageUtil.SaveAsJpg(imageFile, ext, savePath); err != nil {
		result.Error(exception.ImageSaveError).JSON(c)
		return
	}

	url := fmt.Sprintf("%s%s", r.Config.File.ImageUrlPrefix, filename)
	result.Ok().PutData("url", url).PutData("size", imageHeader.Size).JSON(c)
}

// @Router               /v1/raw/image/{filename} [GET]
// @Summary              获取图片
// @Tag                  Raw
// @Param                filename path string true "图片文件名"
// @ResponseHeader 200   { "Content-Type": "image/jpeg" }
func (r *RawController) RawImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := fmt.Sprintf("%s%s", r.Config.File.ImagePath, filename)
	if !util.CommonUtil.IsDirOrFileExist(filePath) {
		result.Error(exception.ImageNotFoundError).JSON(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}
