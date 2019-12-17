package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/po"
)

type passDao struct{}

var PassDao = new(passDao)

func (p *passDao) QueryByUsername(username string) *po.PassRecord {
	// SELECT * FROM `tbl_user`  WHERE `tbl_user`.`deleted_at` IS NULL AND ((`tbl_user`.`username` = '81')) ORDER BY `tbl_user`.`uid` ASC LIMIT 1
	// SELECT * FROM `tbl_password`  WHERE `tbl_password`.`deleted_at` IS NULL AND `tbl_password`.`uid` = 8 AND ((`tbl_password`.`uid` = 8)) ORDER BY `tbl_password`.`uid` ASC LIMIT 1
	user := &po.User{Username: username}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	pass := &po.PassRecord{Uid: user.Uid}
	if DB.Where(pass).First(pass).RecordNotFound() {
		return nil
	}
	pass.User = user
	return pass
}

func (p *passDao) Insert(pass *po.PassRecord) DbStatus {
	// INSERT  INTO `tbl_password` (`encrypted_pass`,`created_at`,`updated_at`,`deleted_at`) VALUES ('e10adc3949ba59abbe56e057f20f883e','2019-12-17 11:21:40','2019-12-17 11:21:40',NULL)
	// UPDATE `tbl_user` SET `username` = '81', `sex` = '', `profile` = '', `avatar_url` = '', `birth_time` = '0001-01-01 00:00:00', `authority` = '', `register_ip` = '::1', `phone_number` = '', `updated_at` = '2019-12-17 11:21:40', `deleted_at` = NULL  WHERE `tbl_user`.`deleted_at` IS NULL AND `tbl_user`.`uid` = 8
	// SELECT * FROM `tbl_user`  WHERE `tbl_user`.`deleted_at` IS NULL AND `tbl_user`.`uid` = 8 ORDER BY `tbl_user`.`uid` ASC LIMIT 1
	// INSERT  INTO `tbl_user` (`uid`,`username`,`profile`,`avatar_url`,`register_ip`,`phone_number`,`created_at`,`updated_at`,`deleted_at`) VALUES (8,'81','','','::1','','2019-12-17 11:21:40','2019-12-17 11:21:40',NULL)
	// SELECT `sex`, `birth_time`, `authority` FROM `tbl_user`  WHERE (uid = 8)
	if err := DB.Create(pass).Error; err != nil {
		if IsDuplicateError(err) {
			return DbExisted
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (p *passDao) Update(pass *po.PassRecord) DbStatus {
	if err := DB.Model(pass).Update(pass).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}
