package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/model/constant"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gorm.io/gorm"
)

type UserService struct {
	db             *gorm.DB
	common         *CommonService
	orderbyService *OrderbyService
}

func NewUserService() *UserService {
	return &UserService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (u *UserService) QueryAll(pp *param.PageOrderParam) ([]*po.User, int32, error) {
	total := int64(0)
	rdb := u.db.Model(&po.User{}).Count(&total)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	users := make([]*po.User, 0)
	rdb = xgorm.WithDB(u.db).Pagination(pp.Limit, pp.Page).Model(&po.User{}).Order(u.orderbyService.User(pp.Order)).Find(&users)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return users, int32(total), nil
}

func (u *UserService) QueryByUid(uid uint64) (*po.User, error) {
	user := &po.User{}
	rdb := u.db.Model(&po.User{}).Where("uid = ?", uid).First(user)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return user, nil
}

func (u *UserService) QueryByUids(uids []uint64) ([]*po.User, error) {
	if len(uids) == 0 {
		return []*po.User{}, nil
	}

	users := make([]*po.User, 0)
	where := u.common.BuildOrExpr("uid", uids)
	rdb := u.db.Model(&po.User{}).Where(where).Find(&users)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]*po.User, len(uids))
	for _, user := range users {
		bucket[user.Uid] = user
	}
	out := make([]*po.User, len(uids))
	for idx, uid := range uids {
		if user, ok := bucket[uid]; ok {
			out[idx] = user
		}
	}
	return out, nil
}

func (u *UserService) Existed(uid uint64) (bool, error) {
	cnt := int64(0)
	rdb := u.db.Model(&po.User{}).Where("uid = ?", uid).Count(&cnt)
	if rdb.Error != nil {
		return false, rdb.Error
	}

	return cnt > 0, nil
}

func (u *UserService) Update(uid uint64, user *param.UpdateUserParam) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where("uid = ?", uid).Updates(user.ToMap())
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) UpdateRole(uid uint64, role string) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where("uid = ?", uid).Update("role", role)
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) UpdateState(uid uint64, state constant.UserState) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where("uid = ?", uid).Update("state", state)
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) Delete(uid uint64) (xstatus.DbStatus, error) {
	tx := u.db.Begin()

	rdb := tx.Model(&po.User{}).Where("uid = ?", uid).Delete(&po.User{})
	status, err := xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status, err
	}

	rdb = tx.Model(&po.Account{}).Where("uid = ?", uid).Delete(&po.Account{})
	status, err = xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status, err
	}

	tx.Commit()
	return xstatus.DbSuccess, nil
}
