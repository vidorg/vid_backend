package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SearchDao struct {
	Config   *config.ServerConfig `di:"~"`
	Db       *gorm.DB             `di:"~"`
	VideoDao *VideoDao            `di:"~"`

	PageSize int32 `di:"-"`
}

func NewSearchDao(dic *xdi.DiContainer) *SearchDao {
	repo := &SearchDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	repo.PageSize = repo.Config.MySqlConfig.PageSize
	return repo
}

func (s *SearchDao) SearchUser(against string, page int32) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := SearchHelper(s.Db, &po.User{}, s.PageSize, page, "username, profile", against, &users)
	return users, total
}

func (s *SearchDao) SearchVideo(against string, page int32) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := SearchHelper(s.Db, &po.Video{}, s.PageSize, page, "title, description", against, &videos)
	for _, video := range videos {
		s.VideoDao.WrapVideo(video)
	}
	return videos, total
}
