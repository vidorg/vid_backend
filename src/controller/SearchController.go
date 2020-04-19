package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/service"
	"strings"
)

type SearchController struct {
	Config         *config.ServerConfig    `di:"~"`
	Logger         *logrus.Logger          `di:"~"`
	Mappers        *xentity.EntityMappers  `di:"~"`
	SearchService  *service.SearchService  `di:"~"`
	SegmentService *service.SegmentService `di:"~"`
}

func NewSearchController(dic *xdi.DiContainer) *SearchController {
	ctrl := &SearchController{}
	dic.MustInject(ctrl)
	return ctrl
}

// @Router              /v1/search/user [GET]
// @Summary             搜索用户
// @Tag                 Search
// @Template            Page
// @Param               key query string true "搜索关键词"
// @ResponseModel 200   #Result<Page<UserDto>>
func (s *SearchController) SearchUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	page := param.BindPage(c, s.Config)

	keys := s.SegmentService.Seg(key)
	against := s.SegmentService.Cat(keys)
	users, total := s.SearchService.SearchUser(against, page)

	retDto := xcondition.First(s.Mappers.MapSlice(xslice.Sti(users), &dto.UserDto{})).([]*dto.UserDto)
	result.Ok().SetPage(total, page.Page, page.Limit, retDto).JSON(c)
}

// @Router              /v1/search/video [GET]
// @Summary             搜索视频
// @Tag                 Search
// @Template            Page
// @Param               key query string true "搜索关键词"
// @ResponseModel 200   #Result<Page<VideoDto>>
func (s *SearchController) SearchVideo(c *gin.Context) {
	key := strings.TrimSpace(c.DefaultQuery("key", ""))
	if key == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	page := param.BindPage(c, s.Config)

	keys := s.SegmentService.Seg(key)
	against := s.SegmentService.Cat(keys)
	videos, total := s.SearchService.SearchVideo(against, page)

	retDto := xcondition.First(s.Mappers.MapSlice(xslice.Sti(videos), &dto.VideoDto{})).([]*dto.VideoDto)
	result.Ok().SetPage(total, page.Page, page.Limit, retDto).JSON(c)
}
