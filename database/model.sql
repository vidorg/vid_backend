
-- tbl_user

CREATE TABLE `tbl_user` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL,
  `profile` varchar(120) DEFAULT NULL,
  `register_time` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8

-- tbl_passrecord

CREATE TABLE `tbl_passrecord` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `encrypted_pass` char(128) NOT NULL,
  PRIMARY KEY (`uid`),
  CONSTRAINT `tbl_passrecord_uid_tbl_user_uid_foreign` FOREIGN KEY (`uid`) REFERENCES `tbl_user` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- tbl_subscribe

CREATE TABLE `tbl_subscribe` (
  `up_uid` int(11) NOT NULL,
  `subscriber_uid` int(11) NOT NULL,
  PRIMARY KEY (`up_uid`,`subscriber_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8