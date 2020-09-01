package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/constant"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"strings"
)

type UserService struct {
	db      *gorm.DB
	orderBy func(string) string
}

func NewUserService() *UserService {
	return &UserService{
		db:      xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		orderBy: xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
	}
}

func (u *UserService) QueryAll(pp *param.PageOrderParam) ([]*po.User, int32, error) {
	total := int32(0)
	u.db.Model(&po.User{}).Count(&total)

	users := make([]*po.User, 0)
	rdb := xgorm.WithDB(u.db).Pagination(pp.Limit, pp.Page).Model(&po.User{}).Order(u.orderBy(pp.Order)).Find(&users)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return users, total, nil
}

func (u *UserService) QueryByUids(uids []uint64) ([]*po.User, error) {
	if len(uids) == 0 {
		return []*po.User{}, nil
	}

	sp := strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("uid = %d OR ", uid))
	}
	where := sp.String()[:sp.Len()-4]
	users := make([]*po.User, 0)
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
		user, ok := bucket[uid]
		if ok {
			out[idx] = user
		}
	}

	return out, nil
}

func (u *UserService) QueryByUid(uid uint64) (*po.User, error) {
	user := &po.User{}
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).First(user)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return user, nil
}

func (u *UserService) QueryByUsername(username string) (*po.User, error) {
	user := &po.User{}
	rdb := u.db.Model(&po.User{}).Where(&po.User{Username: username}).First(user)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return user, nil
}

func (u *UserService) Existed(uid uint64) (bool, error) {
	cnt := 0
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).Count(&cnt)
	if rdb.Error != nil {
		return false, rdb.Error
	}

	return cnt > 0, nil
}

func (u *UserService) Update(uid uint64, user *param.UpdateUserParam) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).Updates(user.ToMap())
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) UpdateRole(uid uint64, role string) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).Update("role", role)
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) UpdateState(uid uint64, state constant.UserState) (xstatus.DbStatus, error) {
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).Update("state", state)
	return xgorm.UpdateErr(rdb)
}

func (u *UserService) Delete(uid uint64) (xstatus.DbStatus, error) {
	tx := u.db.Begin()

	rdb := tx.Model(&po.User{}).Where(&po.User{Uid: uid}).Delete(&po.User{Uid: uid})
	status, err := xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status, err
	}

	rdb = tx.Model(&po.Account{}).Where(&po.Account{Uid: uid}).Delete(&po.Account{Uid: uid})
	status, err = xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status, err
	}

	tx.Commit()
	return xstatus.DbSuccess, nil
}
