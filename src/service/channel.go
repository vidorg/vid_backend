package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gorm.io/gorm"
)

type ChannelService struct {
	db             *gorm.DB
	common         *CommonService
	userService    *UserService
	orderbyService *OrderbyService
}

func NewChannelService() *ChannelService {
	return &ChannelService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userService:    xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (c *ChannelService) QueryAll(pp *param.PageOrderParam) ([]*po.Channel, int32, error) {
	total := int64(0)
	rdb := c.db.Model(&po.Channel{}).Count(&total)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	channels := make([]*po.Channel, 0)
	rdb = xgorm.WithDB(c.db).Pagination(pp.Limit, pp.Page).Model(&po.Channel{}).Order(c.orderbyService.Channel(pp.Order)).Find(&channels)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return channels, int32(total), nil
}

func (c *ChannelService) QueryByUid(uid uint64, pp *param.PageOrderParam) ([]*po.Channel, int32, error) {
	author, err := c.userService.QueryByUid(uid)
	if err != nil {
		return nil, 0, err
	} else if author == nil {
		return nil, 0, nil
	}

	total := int64(0)
	rdb := c.db.Model(&po.Channel{}).Where("author_uid = ?", uid).Count(&total)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	channels := make([]*po.Channel, 0)
	rdb = xgorm.WithDB(c.db).Pagination(pp.Limit, pp.Page).Model(&po.Channel{}).Where("author_uid = ?", uid).Order(c.orderbyService.Channel(pp.Order)).Find(&channels)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return channels, int32(total), nil
}

func (c *ChannelService) QueryByCids(cids []uint64) ([]*po.Channel, error) {
	if len(cids) == 0 {
		return []*po.Channel{}, nil
	}

	channels := make([]*po.Channel, 0)
	where := c.common.BuildOrExpr("cid", cids)
	rdb := c.db.Model(&po.Channel{}).Where(where).Find(&channels)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]*po.Channel, len(channels))
	for _, channel := range channels {
		bucket[channel.Cid] = channel
	}
	out := make([]*po.Channel, len(cids))
	for idx, cid := range cids {
		if channel, ok := bucket[cid]; ok {
			out[idx] = channel
		}
	}
	return out, nil
}

func (c *ChannelService) QueryByCid(cid uint64) (*po.Channel, error) {
	channel := &po.Channel{}
	rdb := c.db.Model(&po.Channel{}).Where("cid = ?", cid).First(&channel)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return channel, nil
}

func (c *ChannelService) QueryCountByUids(uids []uint64) ([]int32, error) {
	if len(uids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntScanResult, 0)
	where := c.common.BuildOrExpr("author_uid", uids)
	rdb := c.db.Model(&po.Channel{}).Select("author_uid as id, count(*) as cnt").Where(where).Group("author_uid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, r := range counts {
		bucket[r.Id] = r.Cnt
	}
	out := make([]int32, len(uids))
	for idx, uid := range uids {
		if cnt, ok := bucket[uid]; ok {
			out[idx] = cnt
		}
	}
	return out, nil
}

func (c *ChannelService) Existed(cid uint64) (bool, error) {
	cnt := int64(0)
	rdb := c.db.Model(&po.Channel{}).Where("cid = ?", cid).Count(&cnt)
	if rdb.Error != nil {
		return false, rdb.Error
	}

	return cnt > 0, nil
}

func (c *ChannelService) Insert(pa *param.InsertChannelParam, uid uint64) (xstatus.DbStatus, error) {
	channel := pa.ToChannelPo()
	channel.AuthorUid = uid

	rdb := c.db.Model(&po.Channel{}).Create(channel)
	return xgorm.CreateErr(rdb)
}

func (c *ChannelService) Update(cid uint64, pa *param.UpdateChannelParam) (xstatus.DbStatus, error) {
	rdb := c.db.Model(&po.Channel{}).Where("cid = ?", cid).Updates(pa.ToMap())
	return xgorm.UpdateErr(rdb)
}

func (c *ChannelService) Delete(cid uint64) (xstatus.DbStatus, error) {
	rdb := c.db.Model(&po.Channel{}).Where("cid = ?", cid).Delete(&po.Channel{})
	return xgorm.DeleteErr(rdb)
}
