package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SearchService struct {
	Db       *helper.GormHelper `di:"~"`
	Logger   *logrus.Logger     `di:"~"`
	VideoDao *VideoService      `di:"~"`
}

func NewSearchService(dic *xdi.DiContainer) *SearchService {
	repo := &SearchService{}
	dic.MustInject(repo)
	return repo
}

func (s *SearchService) SearchUser(against string, page *param.PageParam) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := s.Db.SearchHelper(&po.User{}, page.Limit, page.Page, "username, profile", against, &users)
	return users, total
}

func (s *SearchService) SearchVideo(against string, page *param.PageParam) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := s.Db.SearchHelper(&po.Video{}, page.Limit, page.Page, "title, description", against, &videos)
	for _, video := range videos {
		s.VideoDao.WrapVideo(video)
	}
	return videos, total
}
