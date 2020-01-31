package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
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
// @Security            Jwt
// @Template            Auth ParamA
// @Summary             上传图片
// @Tag                 Raw
// @Param               image formData file true false "上传的图片，大小限制在2M，允许后缀名为 {.jpg, .jpeg, .png, .bmp, .gif}"
// @ResponseDesc 400    "image type not supported"
// @ResponseDesc 413    "request body too large"
// @ResponseDesc 500    "image save failed"
// @ResponseModel 200   #Result<ImageDto>
// @Response 200        ${resp_upload_image}
func (r *RawController) UploadImage(c *gin.Context) {
	imageFile, imageHeader, err := c.Request.FormFile("image")
	if err != nil || imageFile == nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	supported, ext := util.ImageUtil.CheckImageExt(imageHeader.Filename)
	if !supported {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()).JSON(c)
		return
	}

	filename := fmt.Sprintf("%s.jpg", xstring.CurrentTimeUuid(20))
	savePath := fmt.Sprintf("%s%s", r.Config.FileConfig.ImagePath, filename)
	if err := util.ImageUtil.SaveAsJpg(imageFile, ext, savePath); err != nil {
		result.Result{}.Error().SetMessage(exception.ImageSaveError.Error()).JSON(c)
		return
	}

	url := fmt.Sprintf("%s%s", r.Config.FileConfig.ImageUrlPrefix, filename)
	result.Result{}.Ok().PutData("url", url).PutData("size", imageHeader.Size).JSON(c)
}

// @Router               /v1/raw/image/{filename} [GET]
// @Summary              获取图片
// @Tag                  Raw
// @Param                filename path string true false "图片文件名"
// @ResponseDesc 404     "image not found"
// @ResponseHeader 200   { "Content-Type": "image/jpeg" }
func (r *RawController) RawImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := fmt.Sprintf("%s%s", r.Config.FileConfig.ImagePath, filename)
	if !util.CommonUtil.IsDirOrFileExist(filePath) {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.ImageNotFoundError.Error()).JSON(c)
		return
	}

	c.Writer.Header().Add("Content-Type", "image/jpeg")
	c.File(filePath)
}
