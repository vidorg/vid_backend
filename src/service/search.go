package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type SearchService struct {
	db           *gorm.DB
	videoService *VideoService
}

func NewSearchService() *SearchService {
	return &SearchService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		videoService: xdi.GetByNameForce(sn.SVideoService).(*VideoService),
	}
}

func (s *SearchService) SearchUser(against string, page *param.PageParam) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := s.db.SearchHelper(&po.User{}, page.Limit, page.Page, "username, profile", against, &users)
	return users, total
}

func (s *SearchService) SearchVideo(against string, page *param.PageParam) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := s.db.SearchHelper(&po.Video{}, page.Limit, page.Page, "title, description", against, &videos)
	for _, video := range videos {
		s.videoService.WrapVideo(video)
	}
	return videos, total
}
