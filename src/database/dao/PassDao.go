package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type PassDao struct {
	db       *gorm.DB
	pageSize int
}

func PassRepository(config *config.DatabaseConfig) *PassDao {
	return &PassDao{
		db:       database.SetupDBConn(config),
		pageSize: config.PageSize,
	}
}

func (p *PassDao) QueryByUsername(username string) *po.PassRecord {
	user := &po.User{Username: username}
	rdb := p.db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	pass := &po.PassRecord{Uid: user.Uid}
	rdb = p.db.Model(&po.PassRecord{}).Where(pass).First(pass)
	if rdb.RecordNotFound() {
		return nil
	}
	pass.User = user
	return pass
}

func (p *PassDao) Insert(pass *po.PassRecord) database.DbStatus {
	rdb := p.db.Model(&po.PassRecord{}).Create(pass)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}

func (p *PassDao) Update(pass *po.PassRecord) database.DbStatus {
	rdb := p.db.Model(&po.User{}).Update(pass)
	if rdb.Error != nil {
		if database.IsNotFoundError(rdb.Error) {
			return database.DbNotFound
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}
