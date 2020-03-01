package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/common/seg"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"log"
	"strings"
)

type SearchController struct {
	Config         *config.ServerConfig  `di:"~"`
	SearchDao      *dao.SearchDao        `di:"~"`
	Mapper         *xmapper.EntityMapper `di:"~"`
	SegmentService *seg.SegmentService   `di:"~"`
}

func NewSearchController(dic *xdi.DiContainer) *SearchController {
	ctrl := &SearchController{}
	if !dic.Inject(ctrl) {
		log.Fatalln("Inject failed")
	}
	return ctrl
}

// @Router              /v1/search/user [GET]
// @Template            Page ParamA
// @Summary             搜索用户
// @Tag                 Search
// @Param               key query string true "搜索关键词"
// @ResponseModel 200   #Result<Page<UserDto>>
// @ResponseEx 200      ${resp_page_users}
func (s *SearchController) SearchUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	page := param.BindQueryPage(c)

	keys := s.SegmentService.Seg(key)
	against := s.SegmentService.CatAgainst(keys)
	users, total := s.SearchDao.SearchUser(against, page)

	retDto := xcondition.First(s.Mapper.MapSlice(xslice.Sti(users), &dto.UserDto{})).([]*dto.UserDto)
	result.Ok().SetPage(total, page, retDto).JSON(c)
}

// @Router              /v1/search/video [GET]
// @Template            Page ParamA
// @Summary             搜索视频
// @Tag                 Search
// @Param               key query string true "搜索关键词"
// @ResponseModel 200   #Result<Page<VideoDto>>
// @ResponseEx 200      ${resp_page_videos}
func (s *SearchController) SearchVideo(c *gin.Context) {
	key := strings.TrimSpace(c.DefaultQuery("key", ""))
	if key == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	page := param.BindQueryPage(c)

	keys := s.SegmentService.Seg(key)
	against := s.SegmentService.CatAgainst(keys)
	videos, total := s.SearchDao.SearchVideo(against, page)

	retDto := xcondition.First(s.Mapper.MapSlice(xslice.Sti(videos), &dto.VideoDto{})).([]*dto.VideoDto)
	result.Ok().SetPage(total, page, retDto).JSON(c)
}
