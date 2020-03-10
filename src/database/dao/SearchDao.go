package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

type SearchDao struct {
	Db       *helper.GormHelper `di:"~"`
	VideoDao *VideoDao          `di:"~"`
}

func NewSearchDao(dic *xdi.DiContainer) *SearchDao {
	repo := &SearchDao{}
	if !dic.Inject(repo) {
		log.Fatalln("Inject failed")
	}
	return repo
}

func (s *SearchDao) SearchUser(against string, page int32, limit int32) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := s.Db.SearchHelper(&po.User{}, limit, page, "username, profile", against, &users)
	return users, total
}

func (s *SearchDao) SearchVideo(against string, page int32, limit int32) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := s.Db.SearchHelper(&po.Video{}, limit, page, "title, description", against, &videos)
	for _, video := range videos {
		s.VideoDao.WrapVideo(video)
	}
	return videos, total
}
