package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type BlockService struct {
	db          *gorm.DB
	userService *UserService
	common      *CommonService
	orderBy     func(string) string
	tblName     string
}

func NewBlockService() *BlockService {
	return &BlockService{
		db:          xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService: xdi.GetByNameForce(sn.SUserService).(*UserService),
		common:      xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		orderBy:     xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
		tblName:     "tbl_block",
	}
}

func (b *BlockService) blockingAssoDB(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Blockings") // association pattern
}

func (b *BlockService) blockingAsso(uid uint64) *gorm.Association {
	return b.blockingAssoDB(b.db, uid)
}

func (b *BlockService) QueryBlockings(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := b.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(b.blockingAsso(uid).Count())
	users := make([]*po.User, 0)
	rac := b.blockingAssoDB(xgorm.WithDB(b.db).Pagination(pp.Limit, pp.Page).Order(b.orderBy(pp.Order)), uid).Find(&users)
	if rac.Error != nil {
		return nil, 0, rac.Error
	}

	return users, total, nil
}

func (b *BlockService) CheckBlockings(me uint64, uids []uint64) ([]bool, error) {
	if len(uids) == 0 {
		return []bool{}, nil
	}

	blockings := make([]*_FromToUidPair, 0)
	where := b.common.BuildOrExp("to_uid", uids)
	rdb := b.db.Table(b.tblName).Select("from_uid, to_uid").Where("from_uid = ?", me).Where(where).Group("from_uid, to_uid").Scan(&blockings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]bool, len(uids))
	for _, pair := range blockings {
		bucket[pair.ToUid] = true
	}
	out := make([]bool, len(uids))
	for idx, uid := range uids {
		user, ok := bucket[uid]
		if ok {
			out[idx] = user
		}
	}
	return out, nil
}

func (b *BlockService) InsertBlocking(uid uint64, to uint64) (xstatus.DbStatus, error) {
	ok1, err1 := b.userService.Existed(uid)
	ok2, err2 := b.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	cnt := 0
	rdb := b.db.Table(b.tblName).Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil
	}

	ras := b.blockingAsso(uid).Append(&po.User{Uid: to})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}

func (b *BlockService) DeleteBlocking(uid uint64, to uint64) (xstatus.DbStatus, error) {
	ok1, err1 := b.userService.Existed(uid)
	ok2, err2 := b.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	cnt := 0
	rdb := b.db.Table(b.tblName).Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil
	}

	ras := b.blockingAsso(uid).Delete(&po.User{Uid: to})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}
