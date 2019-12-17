package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/dto"
	"vid/app/model/po"
)

type userDao struct{}

var UserDao = new(userDao)

func (u *userDao) QueryAll(page int) (users []po.User, count int) {
	DB.Model(&po.User{}).Count(&count)
	DB.Limit(PageSize).Offset((page - 1) * PageSize).Find(&users)
	return users, count
}

func (u *userDao) QueryByUid(uid int) *po.User {
	user := &po.User{Uid: uid}
	// NewRecord: 防止以主键查询时主键为 0
	if DB.NewRecord(user) || DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	return user
}

func (u *userDao) QueryByUsername(username string) *po.User {
	user := &po.User{Username: username}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	return user
}

func (u *userDao) Update(user *po.User) DbStatus {
	if DB.NewRecord(user) {
		return DbNotFound
	}
	if err := DB.Model(user).Update(user).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else if IsDuplicateError(err) {
			return DbExisted
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	// user.BirthTime == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	// 0001-01-01 00:00:00 +0000 UTC
	return DbSuccess
}

func (u *userDao) Delete(uid int) (*po.User, DbStatus) {
	user := u.QueryByUid(uid)
	if user == nil {
		return nil, DbNotFound
	}
	if err := DB.Model(user).Update(user).Error; err != nil {
		if IsNotFoundError(err) {
			return nil, DbNotFound
		} else {
			log.Println(err)
			return nil, DbFailed
		}
	}
	return user, DbSuccess
}

func (u *userDao) QueryUserExtraInfo(isSelfOrAdmin bool, user *po.User) (*dto.UserExtraInfo, DbStatus) {
	phoneNumber := user.PhoneNumber
	if !isSelfOrAdmin {
		phoneNumber = ""
	}
	subscribingCnt, subscriberCnt, status := SubDao.QuerySubCnt(user.Uid)
	if status == DbNotFound {
		return nil, status
	}

	// video_cnt, err := VideoDao.QueryUserVideoCnt(user.Uid)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// playlist_cnt, err := PlaylistDao.QueryUserPlaylistCnt(user.Uid)
	// if err != nil {
	// 	return nil, err
	// }

	return &dto.UserExtraInfo{
		PhoneNumber:      phoneNumber,
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       0,
		PlaylistCount:    0,
	}, DbSuccess
}
