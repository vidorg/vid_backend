package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/model/common"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

type RawController struct {
	Config *config.ServerConfig  `di:"~"`
	Mapper *xmapper.EntityMapper `di:"~"`
}

func NewRawController(dic *xdi.DiContainer) *RawController {
	ctrl := &RawController{}
	if !dic.Inject(ctrl) {
		panic("Inject failed")
	}
	return ctrl
}

// @Router              /v1/raw/image [POST]
// @Template            Auth
// @Summary             上传图片
// @Description         上传公共图片，包括用户头像和视频封面
// @Tag                 Raw
// @Param               image formData file true "上传的图片，大小限制在2M，允许后缀名为 {.jpg, .jpeg, .png, .bmp, .gif}"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           400 image type not supported
// @ErrorCode           413 request body too large
// @ErrorCode           500 image save failed
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"url": "http://localhost:3344/v1/raw/image/20200110130323908439.jpg",
								"size": 381952
							}
 						} */
func (r *RawController) UploadImage(c *gin.Context) {
	imageFile, imageHeader, err := c.Request.FormFile("image")
	if err != nil || imageFile == nil {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	supported, ext := util.ImageUtil.CheckImageExt(imageHeader.Filename)
	if !supported {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()).JSON(c)
		return
	}

	filename := fmt.Sprintf("%s.jpg", xstring.CurrentTimeUuid(20))
	savePath := fmt.Sprintf("%s%s", r.Config.FileConfig.ImagePath, filename)
	if err := util.ImageUtil.SaveAsJpg(imageFile, ext, savePath); err != nil {
		common.Result{}.Error().SetMessage(exception.ImageSaveError.Error()).JSON(c)
		return
	}

	url := fmt.Sprintf("%s%s", r.Config.FileConfig.ImageUrlPrefix, filename)
	common.Result{}.Ok().PutData("url", url).PutData("size", imageHeader.Size).JSON(c)
}

// @Router              /v1/raw/image/{filename} [GET]
// @Summary             获取图片
// @Description         获取用户头像图片以及视频封面
// @Tag                 Raw
// @Param               filename path string true "图片文件名，jpg后缀名"
// @Accept              multipart/form-data
// @ErrorCode           404 image not found
/* @Response 200        {| "Content-Type": "image/jpeg" |} */
func (r *RawController) RawImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := fmt.Sprintf("%s%s", r.Config.FileConfig.ImagePath, filename)
	if !util.CommonUtil.IsDirOrFileExist(filePath) {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.ImageNotFoundError.Error()).JSON(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}
