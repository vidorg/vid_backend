package database

import (
	"time"

	"vid/config"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"
)

type userDao struct{}

var UserDao = new(userDao)

const (
	col_user_uid           = "uid"
	col_user_username      = "username"
	col_user_profile       = "profile"
	col_user_sex           = "sex"
	col_user_avatar_url    = "avatar_url"
	col_user_birth_time    = "birth_time"
	col_user_register_time = "register_time"
)

// db 查询所有用户
//
// @return `[]User`
func (u *userDao) QueryAllUsers() (users []User) {
	DB.Find(&users)
	for k, _ := range users {
		users[k].ToServer()
	}
	return users
}

// db 查询 uid 用户
//
// @return `*user` `isUserExist`
func (u *userDao) QueryUserByUid(uid int) (*User, bool) {
	var user User
	nf := DB.Where(col_user_uid+" = ?", uid).Find(&user).RecordNotFound()
	if nf {
		return nil, false
	} else {
		user.ToServer()
		return &user, true
	}
}

// db 查询 username 用户
//
// @return `*user` `isUserExist`
func (u *userDao) QueryUserByUserName(username string) (*User, bool) {
	var user User
	nf := DB.Where(col_user_username+" = ?", username).Find(&user).RecordNotFound()
	if nf {
		return nil, false
	} else {
		user.ToServer()
		return &user, true
	}
}

// db 根据用户名模糊查询用户
//
// @return `[]user`
func (u *userDao) SearchByUserName(username string) (users []User) {
	DB.Where(col_user_username+" like ?", "%"+username+"%").Find(&users).RecordNotFound()
	for _, v := range users {
		v.ToServer()
	}
	return users
}

///////////////////////////////////////////////////////////////////////////////////////////

// db 更新用户名和简介
//
// @return `*user` `err`
//
// @error `UserNotExistException` `UserInfoException` `UserNameUsedException` `NotUpdateUserException`
func (u *userDao) UpdateUser(user User) (*User, error) {
	// 检查用户信息
	queryBefore, ok := u.QueryUserByUid(user.Uid)
	if !ok {
		return nil, UserNotExistException
	}

	user.ToDB()

	// 更新空字段
	if user.Username == "" {
		user.Username = queryBefore.Username
	}
	if user.Profile == config.AppCfg.MagicToken {
		user.Profile = queryBefore.Profile
	}
	if user.Sex == "" {
		user.Sex = queryBefore.Sex
	}
	if user.AvatarUrl == "" {
		user.AvatarUrl = queryBefore.AvatarUrl
	}
	if user.BirthTime == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
		// 0001-01-01 00:00:00 +0000 UTC
		user.BirthTime = queryBefore.BirthTime
	}

	// 检查格式
	if !user.CheckFormat() {
		return nil, UserInfoException
	}
	// 检查同名
	if _, ok = u.QueryUserByUserName(user.Username); ok && user.Username != queryBefore.Username {
		return nil, UserNameUsedException
	}

	DB.Model(&user).Updates(map[string]interface{}{
		col_user_username:   user.Username,
		col_user_profile:    user.Profile,
		col_user_sex:        user.Sex,
		col_user_avatar_url: user.AvatarUrl,
		col_user_birth_time: user.BirthTime,
	})
	// 检查更新后
	query, _ := u.QueryUserByUid(user.Uid)
	if queryBefore.Equals(query) {
		// 数据不变
		query.ToServer()
		return query, NotUpdateUserException
	} else {
		// 正常
		query.ToServer()
		return query, nil
	}
}

// db 删除用户和用户密码 (cascade)
//
// @return `*user` `err`
//
// @error `UserNotExistException` `DeleteUserException`
func (u *userDao) DeleteUser(uid int) (*User, error) {

	query, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}

	if DB.Delete(query).RowsAffected != 1 {
		return nil, DeleteUserException
	} else {
		query.ToServer()
		return query, nil
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

// db `suberUip` 关注 `upUid`
//
// @return `err`
//
// @error `UserNotExistException` `SubscribeOneSelfException`
func (u *userDao) SubscribeUser(suberUid int, upUid int) error {
	upUser, ok := u.QueryUserByUid(upUid)
	if !ok {
		return UserNotExistException
	}
	suberUser, ok := u.QueryUserByUid(suberUid)
	if !ok {
		return UserNotExistException
	}
	if upUid == suberUid {
		return SubscribeOneSelfException
	}
	DB.Model(upUser).Association("Subscribers").Append(suberUser)
	return nil
}

// db `suberUip` 取消关注 `upUid`
//
// @return `err`
//
// @error `UserNotExistException` `SubscribeOneSelfException`
func (u *userDao) UnSubscribeUser(suberUid int, upUid int) error {
	upUser, ok := u.QueryUserByUid(upUid)
	if !ok {
		return UserNotExistException
	}
	suberUser, ok := u.QueryUserByUid(suberUid)
	if !ok {
		return UserNotExistException
	}
	if upUid == suberUid {
		return SubscribeOneSelfException
	}
	DB.Model(upUser).Association("Subscribers").Delete(suberUser)
	return nil
}

// db 查询 uid 的粉丝
//
// @return `user[]` `err`
//
// error `UserNotExistException`
func (u *userDao) QuerySubscriberUsers(uid int) ([]User, error) {
	user, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}
	var users []User
	DB.Model(user).Related(&users, "Subscribers")
	for _, v := range users {
		v.ToServer()
	}
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`subscriber_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`user_uid` IN (5))
	return users, nil
}

// db 查询 uid 的关注
//
// @return `user[]` `err`
//
// @error `UserNotExistException`
func (u *userDao) QuerySubscribingUsers(uid int) ([]User, error) {
	user, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}
	var users []User
	DB.Model(user).Related(&users, "Subscribings")
	for _, v := range users {
		v.ToServer()
	}
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`user_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`subscriber_uid` IN (5))
	return users, nil
}

// db 查询 uid 的关注和粉丝数
//
// @return `subing_cnt` `suber_cnt` `err`
//
// @error `UserNotExistException`
func (u *userDao) QuerySubCnt(uid int) (int, int, error) {
	user, ok := u.QueryUserByUid(uid)
	if !ok {
		return 0, 0, UserNotExistException
	}
	var subing []User
	DB.Model(user).Related(&subing, "Subscribings")
	var suber []User
	DB.Model(user).Related(&suber, "Subscribers")
	return len(subing), len(suber), nil
}

///////////////////////////////////////////////////////////////////////////////////////////

// db 查询 uid 其他数据
//
// @return `*UserExtraInfo` `err`
//
// @error `UserNotExistException`
func (u *userDao) QueryUserExtraInfo(isAuth bool, user *User) (*UserExtraInfo, error) {
	subing_cnt, suber_cnt, err := UserDao.QuerySubCnt(user.Uid)
	if err != nil {
		return nil, err
	}

	video_cnt, err := VideoDao.QueryUserVideoCnt(user.Uid)
	if err != nil {
		return nil, err
	}

	playlist_cnt, err := PlaylistDao.QueryUserPlaylistCnt(user.Uid)
	if err != nil {
		return nil, err
	}

	phone_number := user.PhoneNumber
	if !isAuth {
		phone_number = 0
	}

	return &UserExtraInfo{
		PhoneNumber:     phone_number,
		Subscriber_cnt:  suber_cnt,
		Subscribing_cnt: subing_cnt,
		Video_cnt:       video_cnt,
		Playlist_cnt:    playlist_cnt,
	}, nil
}
