-- db: db_vid

-- host: localhost:3306
-- charset: utf8

-- tbl: tbl_password

CREATE TABLE `tbl_password`
(
    `uid`            int(11)   NOT NULL AUTO_INCREMENT,
    `encrypted_pass` char(48)  NOT NULL,
    `created_at`     timestamp NULL DEFAULT NULL,
    `updated_at`     timestamp NULL DEFAULT NULL,
    `deleted_at`     timestamp NULL DEFAULT NULL,
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
    `birth_time`   timestamp NULL                   DEFAULT NULL,
    `authority`    enum ('admin','normal')          DEFAULT 'normal',
    `register_ip`  varchar(255)                     DEFAULT NULL,
    `phone_number` varchar(11)                      DEFAULT NULL,
    `created_at`   timestamp NULL                   DEFAULT NULL,
    `updated_at`   timestamp NULL                   DEFAULT NULL,
    `deleted_at`   timestamp NULL                   DEFAULT NULL,
    PRIMARY KEY (`uid`),
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- tbl: tbl_subscribe

CREATE TABLE `tbl_subscribe`
(
    `up_uid`         int(11) NOT NULL,
    `subscriber_uid` int(11) NOT NULL,
    PRIMARY KEY (`up_uid`, `subscriber_uid`)
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
    `upload_time` datetime       DEFAULT CURRENT_TIMESTAMP,
    `author_uid`  int(11)        DEFAULT NULL,
    `created_at`  timestamp NULL DEFAULT NULL,
    `updated_at`  timestamp NULL DEFAULT NULL,
    `deleted_at`  timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`vid`),
    UNIQUE KEY `video_url` (`video_url`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
