package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SearchDao struct {
	Config *config.ServerConfig `di:"~"`
	Db     *gorm.DB             `di:"~"`

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
	user := make([]*po.User, 0)
	var total int32
	rdb := s.Db.Model(&po.User{}).Where("match (username, profile) against (? in boolean mode)", against)
	rdb.Limit(s.PageSize).Offset((page - 1) * s.PageSize).Find(&user)
	rdb.Count(&total)
	return user, total
}
