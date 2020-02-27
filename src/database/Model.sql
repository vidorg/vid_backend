-- db: db_vid

-- host: localhost:3306
-- charset: utf8

-- tbl: tbl_account

CREATE TABLE `tbl_account`
(
    `uid`            int(11)      NOT NULL AUTO_INCREMENT,
    `encrypted_pass` varchar(255) NOT NULL,
    `created_at`     timestamp    NULL DEFAULT NULL,
    `updated_at`     timestamp    NULL DEFAULT NULL,
    `deleted_at`     timestamp    NULL DEFAULT '2000-01-01 00:00:00',
    PRIMARY KEY (`uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- tbl: tbl_user

CREATE TABLE `tbl_user`
(
    `uid`          int(11)   NOT NULL AUTO_INCREMENT,
    `username`     varchar(30)                      DEFAULT NULL,
    `sex`          enum ('unknown','male','female') DEFAULT 'unknown',
    `profile`      varchar(255)                     DEFAULT NULL,
    `avatar_url`   varchar(255)                     DEFAULT NULL,
    `birthday`     date                             DEFAULT '2000-01-01',
    `authority`    enum ('admin','normal')          DEFAULT 'normal',
    `register_ip`  varchar(15)                      DEFAULT NULL,
    `phone_number` varchar(11)                      DEFAULT NULL,
    `created_at`   timestamp NULL                   DEFAULT NULL,
    `updated_at`   timestamp NULL                   DEFAULT NULL,
    `deleted_at`   timestamp NULL                   DEFAULT '2000-01-01 00:00:00',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `idx_user_username_deleted_at_unique` (`username`, `deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- tbl: tbl_subscribe

CREATE TABLE `tbl_subscribe`
(
    `subscriber_uid` int(11) NOT NULL,
    `up_uid`         int(11) NOT NULL,
    PRIMARY KEY (`subscriber_uid`, `up_uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- tbl: tbl_video

CREATE TABLE `tbl_video`
(
    `vid`         int(11)   NOT NULL AUTO_INCREMENT,
    `title`       varchar(100)   DEFAULT NULL,
    `description` varchar(255)   DEFAULT NULL,
    `video_url`   varchar(255)   DEFAULT NULL,
    `cover_url`   varchar(255)   DEFAULT NULL,
    `author_uid`  int(11)        DEFAULT NULL,
    `created_at`  timestamp NULL DEFAULT NULL,
    `updated_at`  timestamp NULL DEFAULT NULL,
    `deleted_at`  timestamp NULL DEFAULT '2000-01-01 00:00:00',
    PRIMARY KEY (`vid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
