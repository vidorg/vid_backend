package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SubDao struct {
	db       *gorm.DB
	pageSize int
	userDao  *UserDao

	colSubscribers  string
	colSubscribings string
}

func SubRepository(config *config.DatabaseConfig) *SubDao {
	return &SubDao{
		db:       database.SetupDBConn(config),
		pageSize: config.PageSize,
		userDao:  UserRepository(config),

		colSubscribers:  "Subscribers",
		colSubscribings: "Subscribings",
	}
}

func (s *SubDao) QuerySubscriberUsers(uid int, page int) (users []*po.User, count int) {
	user := &po.User{Uid: uid}
	if !s.userDao.Exist(uid) {
		return nil, 0
	}
	count = s.db.Model(user).Association(s.colSubscribers).Count() // 开始关联模式
	s.db.Limit(s.pageSize).Offset((page-1)*s.pageSize).Model(user).Related(&users, s.colSubscribers)
	return users, count
}

func (s *SubDao) QuerySubscribingUsers(uid int, page int) (users []*po.User, count int) {
	user := &po.User{Uid: uid}
	if !s.userDao.Exist(uid) {
		return nil, 0
	}
	count = s.db.Model(user).Association(s.colSubscribings).Count()
	s.db.Limit(s.pageSize).Offset((page-1)*s.pageSize).Model(user).Related(&users, s.colSubscribings)
	return users, count
}

func (s *SubDao) QuerySubCnt(uid int) (subscribingCnt int, subscriberCnt int, status database.DbStatus) {
	if !s.userDao.Exist(uid) {
		return 0, 0, database.DbNotFound
	}
	user := &po.User{Uid: uid}
	subscribingCnt = s.db.Model(user).Association(s.colSubscribings).Count()
	subscriberCnt = s.db.Model(user).Association(s.colSubscribers).Count()
	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubDao) SubscribeUser(meUid int, toUid int) database.DbStatus {
	if !s.userDao.Exist(toUid) {
		return database.DbNotFound
	}

	ass := s.db.Model(&po.User{Uid: toUid}).Association(s.colSubscribers).Append(&po.User{Uid: meUid})
	if ass.Error != nil {
		if database.IsNotFoundError(ass.Error) {
			return database.DbNotFound
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}

func (s *SubDao) UnSubscribeUser(meUid int, toUid int) database.DbStatus {
	if !s.userDao.Exist(toUid) {
		return database.DbNotFound
	}

	ass := s.db.Model(&po.User{Uid: toUid}).Association(s.colSubscribers).Delete(&po.User{Uid: meUid})
	if ass.Error != nil {
		if database.IsNotFoundError(ass.Error) {
			return database.DbNotFound
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}
