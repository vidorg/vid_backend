-- db: db_vid

-- tbl: tbl_user

CREATE TABLE `tbl_user` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL,
  `profile` varchar(255) DEFAULT NULL,
  `register_time` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT DEFAULT CHARSET=utf8

-- tbl: tbl_passrecord

CREATE TABLE `tbl_passrecord` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `encrypted_pass` char(48) NOT NULL,
  PRIMARY KEY (`uid`),
  CONSTRAINT `tbl_passrecord_uid_tbl_user_uid_foreign` FOREIGN KEY (`uid`) REFERENCES `tbl_user` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT DEFAULT CHARSET=utf8

-- tbl: tbl_subscribe

CREATE TABLE `tbl_subscribe` (
  `up_uid` int(11) NOT NULL,
  `subscriber_uid` int(11) NOT NULL,
  PRIMARY KEY (`up_uid`,`subscriber_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- tbl: tbl_video

CREATE TABLE `tbl_video` (
  `vid` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `video_url` varchar(255) DEFAULT NULL,
  `upload_time` timestamp NULL DEFAULT NULL,
  `author_uid` int(11) DEFAULT NULL,
  PRIMARY KEY (`vid`),
  KEY `tbl_video_author_uid_tbl_user_uid_foreign` (`author_uid`),
  CONSTRAINT `tbl_video_author_uid_tbl_user_uid_foreign` FOREIGN KEY (`author_uid`) REFERENCES `tbl_user` (`uid`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT DEFAULT CHARSET=utf8