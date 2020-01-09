package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type PassDao struct {
	config *config.MySqlConfig
	db     *gorm.DB
}

func PassRepository(config *config.MySqlConfig) *PassDao {
	return &PassDao{
		config: config,
		db:     database.SetupDBConn(config),
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
	rdb := p.db.Model(&po.PassRecord{}).Create(pass) // cascade create
	if database.IsDuplicateError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (p *PassDao) Update(pass *po.PassRecord) database.DbStatus {
	rdb := p.db.Model(&po.PassRecord{}).Update(pass)
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
