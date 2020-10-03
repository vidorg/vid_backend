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
	userService    *UserService
	orderbyService *OrderbyService
}

func NewChannelService() *ChannelService {
	return &ChannelService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
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
