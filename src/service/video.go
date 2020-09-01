package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"strings"
)

type VideoService struct {
	db          *gorm.DB
	userService *UserService
	orderBy     func(string) string
}

func NewVideoService() *VideoService {
	return &VideoService{
		db:          xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService: xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderBy:     xgorm.OrderByFunc(dto.BuildVideoPropertyMapper()),
	}
}

func (v *VideoService) QueryAll(pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	total := int32(0)
	v.db.Model(&po.Video{}).Count(&total)

	videos := make([]*po.Video, 0)
	rdb := xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).Model(&po.User{}).Order(v.orderBy(pp.Order)).Find(&videos)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return videos, total, nil
}

func (v *VideoService) QueryByUid(uid uint64, pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	author, err := v.userService.QueryByUid(uid)
	if err != nil {
		return nil, 0, err
	} else if author == nil {
		return nil, 0, nil
	}

	total := int32(0)
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	videos := make([]*po.Video, 0)
	rdb := xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).Model(&po.Video{}).Order(v.orderBy(pp.Order)).Where(&po.Video{AuthorUid: uid}).Find(&videos)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return videos, total, nil
}

func (v *VideoService) QueryByVid(vid uint64) (*po.Video, error) {
	video := &po.Video{}
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).First(video)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return video, nil
}

func (v *VideoService) QueryCountByUids(uids []uint64) ([]int32, error) {
	if len(uids) == 0 {
		return []int32{}, nil
	}

	type result struct {
		Id  uint64
		Cnt int32
	}

	sp := strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("author_uid = %d OR ", uid))
	}
	where := sp.String()[:sp.Len()-4]
	counts := make([]*result, 0)
	rdb := v.db.Model(&po.Video{}).Select("author_uid as id, count(*) as cnt").Where(where).Group("author_uid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, cnt := range counts {
		bucket[cnt.Id] = cnt.Cnt
	}

	out := make([]int32, len(uids))
	for idx, uid := range uids {
		cnt, ok := bucket[uid]
		if ok {
			out[idx] = cnt
		}
	}
	return out, nil
}

func (v *VideoService) Existed(vid uint64) (bool, error) {
	cnt := 0
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Count(&cnt)
	if rdb.Error != nil {
		return false, rdb.Error
	}

	return cnt > 0, nil
}

func (v *VideoService) Insert(pa *param.InsertVideoParam, uid uint64) (xstatus.DbStatus, error) {
	video := pa.ToPo()
	video.AuthorUid = uid

	rdb := v.db.Model(&po.Video{}).Create(video)
	return xgorm.CreateErr(rdb)
}

func (v *VideoService) Update(vid uint64, video *param.UpdateVideoParam) (xstatus.DbStatus, error) {
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Update(video.ToMap())
	return xgorm.UpdateErr(rdb)
}

func (v *VideoService) Delete(vid uint64) (xstatus.DbStatus, error) {
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Delete(&po.Video{Vid: vid})
	return xgorm.DeleteErr(rdb)
}
