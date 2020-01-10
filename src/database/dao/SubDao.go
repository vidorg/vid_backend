package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SubDao struct {
	config  *config.MySqlConfig
	db      *gorm.DB
	userDao *UserDao

	colSubscribers  string
	colSubscribings string
}

func SubRepository(config *config.MySqlConfig) *SubDao {
	return &SubDao{
		config:  config,
		db:      database.SetupDBConn(config),
		userDao: UserRepository(config),

		colSubscribers:  "Subscribers",
		colSubscribings: "Subscribings",
	}
}

func (s *SubDao) QuerySubscriberUsers(uid int32, page int32) (users []*po.User, count int32, status database.DbStatus) {
	// https://gorm.io/docs/many_to_many.html
	user := &po.User{Uid: uid}
	if !s.userDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count = int32(s.db.Model(user).Association(s.colSubscribers).Count()) // 开始关联模式
	s.db.Limit(s.config.PageSize).Offset((page-1)*s.config.PageSize).Model(user).Related(&users, s.colSubscribers)
	return users, count, database.DbSuccess
}

func (s *SubDao) QuerySubscribingUsers(uid int32, page int32) (users []*po.User, count int32, status database.DbStatus) {
	user := &po.User{Uid: uid}
	if !s.userDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count = int32(s.db.Model(user).Association(s.colSubscribings).Count())
	s.db.Limit(s.config.PageSize).Offset((page-1)*s.config.PageSize).Model(user).Related(&users, s.colSubscribings)
	return users, count, database.DbSuccess
}

func (s *SubDao) QuerySubCnt(uid int32) (subscribingCnt int32, subscriberCnt int32, status database.DbStatus) {
	if !s.userDao.Exist(uid) {
		return 0, 0, database.DbNotFound
	}
	user := &po.User{Uid: uid}
	subscribingCnt = int32(s.db.Model(user).Association(s.colSubscribings).Count())
	subscriberCnt = int32(s.db.Model(user).Association(s.colSubscribers).Count())
	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubDao) SubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userDao.Exist(toUid) || !s.userDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.db.Model(&po.User{Uid: toUid}).Association(s.colSubscribers).Append(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (s *SubDao) UnSubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userDao.Exist(toUid) || !s.userDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.db.Model(&po.User{Uid: toUid}).Association(s.colSubscribers).Delete(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
