-- Host: localhost  (Version 8.0.15)
-- Date: 2020-09-02 22:25:42
-- Generator: MySQL-Front 6.0  (Build 2.20)


--
-- Structure for table "tbl_account"
--

CREATE TABLE `tbl_account`
(
    `uid`        bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `password`   varchar(255)        NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT '1970-01-01 00:00:00',
    PRIMARY KEY (`uid`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 12
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_account"
--

INSERT INTO `tbl_account`
VALUES (1, '$2a$10$.PQPJREfWgtXM/PGTselSeHb5pzcHbtsNFF1sKCEZAmu3.YK/8duS', '2020-09-01 10:36:21', '2020-09-01 10:55:08',
        '1970-01-01 00:00:00'),
       (2, '$2a$10$gWud8tX4BBji5vA2PCRK1eR/MQhjTPUPzTLMEgKHCZ8eF.3kXRFYe', '2020-09-01 10:36:42', '2020-09-01 10:36:42',
        '1970-01-01 00:00:00'),
       (3, '$2a$10$YjsImy0zvxpW7U5.psHdt.vTPQ5vu/FfU7WhFo5Y0a6fNjGWs2eiW', '2020-09-01 10:36:51', '2020-09-01 10:36:51',
        '1970-01-01 00:00:00'),
       (4, '$2a$10$VVgcDut9KwOqpY9Zju1freZHI6LclDdp43ps/u/R2PL/6lcjxWcvi', '2020-09-01 10:37:18', '2020-09-01 10:37:18',
        '1970-01-01 00:00:00'),
       (5, '$2a$10$IEzFC/Wy1z1LpBd/yH68keNYbAuXpRfOu9RWrYRCydEtdK39af4d.', '2020-09-01 10:37:21', '2020-09-01 10:37:21',
        '1970-01-01 00:00:00'),
       (6, '$2a$10$//PLRBNUw1TDTTaE.9463.tuewvME7.4E1zBRQ/1blgJCik6.lbA2', '2020-09-01 10:37:26', '2020-09-01 10:37:26',
        '1970-01-01 00:00:00'),
       (7, '$2a$10$NuA9YrdR0q3LkP/NPme9k.KBxpHJ6RnGxDVDflAwUJf01BJ5nkPNi', '2020-09-01 10:37:30', '2020-09-01 10:37:30',
        '1970-01-01 00:00:00'),
       (8, '$2a$10$ZAvv1VF.McUuSaHsixAPn.hWy5dDuWiy0A9/a2KbgSItI95hHxvuG', '2020-09-01 10:37:32', '2020-09-01 10:37:32',
        '1970-01-01 00:00:00');

--
-- Structure for table "tbl_block"
--

CREATE TABLE `tbl_block`
(
    `from_uid` bigint(20) unsigned NOT NULL,
    `to_uid`   bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`from_uid`, `to_uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_block"
--

INSERT INTO `tbl_block`
VALUES (1, 2),
       (1, 4),
       (1, 6),
       (1, 8),
       (3, 1);

--
-- Structure for table "tbl_casbin_rule"
--

CREATE TABLE `tbl_casbin_rule`
(
    `p_type` varchar(100) DEFAULT NULL,
    `v0`     varchar(100) DEFAULT NULL,
    `v1`     varchar(100) DEFAULT NULL,
    `v2`     varchar(100) DEFAULT NULL,
    `v3`     varchar(100) DEFAULT NULL,
    `v4`     varchar(100) DEFAULT NULL,
    `v5`     varchar(100) DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_casbin_rule"
--

INSERT INTO `tbl_casbin_rule`
VALUES ('g', 'sub_root', 'normal', '', '', '', ''),
       ('p', 'normal', '/v1/auth/user', 'GET', '', '', ''),
       ('p', 'normal', '/v1/auth/logout', 'DELETE', '', '', ''),
       ('p', 'normal', '/v1/auth/password', 'PUT', '', '', ''),
       ('p', 'normal', '/v1/auth/activate', 'POST', '', '', ''),
       ('p', 'normal', '/v1/user', 'PUT', '', '', ''),
       ('p', 'normal', '/v1/user', 'DELETE', '', '', ''),
       ('p', 'normal', '/v1/user/subscribing/:uid', 'DELETE', '', '', ''),
       ('p', 'normal', '/v1/user/subscribing/:uid', 'POST', '', '', ''),
       ('p', 'normal', '/v1/video', 'POST', '', '', ''),
       ('p', 'normal', '/v1/video/:vid', 'PUT', '', '', ''),
       ('p', 'normal', '/v1/video/:vid', 'DELETE', '', '', ''),
       ('p', 'sub_root', '/v1/video', 'GET', '', '', ''),
       ('p', 'sub_root', '/v1/user', 'GET', '', '', '');

--
-- Structure for table "tbl_favorite"
--

CREATE TABLE `tbl_favorite`
(
    `uid` bigint(20) unsigned NOT NULL,
    `vid` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`uid`, `vid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_favorite"
--

INSERT INTO `tbl_favorite`
VALUES (1, 1),
       (1, 4),
       (2, 1),
       (2, 2),
       (2, 4),
       (3, 1),
       (3, 2),
       (3, 4),
       (4, 2),
       (4, 5),
       (5, 2),
       (5, 3),
       (5, 4),
       (5, 5);

--
-- Structure for table "tbl_subscribe"
--

CREATE TABLE `tbl_subscribe`
(
    `from_uid` bigint(20) unsigned NOT NULL,
    `to_uid`   bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`from_uid`, `to_uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_subscribe"
--

INSERT INTO `tbl_subscribe`
VALUES (1, 2),
       (1, 3),
       (1, 4),
       (1, 6),
       (1, 8),
       (2, 1),
       (2, 3),
       (2, 7),
       (3, 2),
       (3, 5),
       (3, 6),
       (3, 8),
       (5, 1),
       (5, 2),
       (5, 4),
       (5, 8),
       (6, 1),
       (6, 5),
       (6, 8);

--
-- Structure for table "tbl_user"
--

CREATE TABLE `tbl_user`
(
    `uid`        bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `username`   varchar(127)        NOT NULL,
    `email`      varchar(255)        NOT NULL,
    `nickname`   varchar(127)        NOT NULL,
    `gender`     tinyint(4)          NOT NULL DEFAULT '0',
    `profile`    varchar(255)        NOT NULL,
    `avatar`     varchar(255)        NOT NULL,
    `birthday`   date                NOT NULL DEFAULT '2000-01-01',
    `role`       varchar(255)        NOT NULL DEFAULT 'normal',
    `state`      tinyint(4)          NOT NULL DEFAULT '0',
    `phone`      varchar(127)        NOT NULL,
    `created_at` datetime                     DEFAULT NULL,
    `updated_at` datetime                     DEFAULT NULL,
    `deleted_at` datetime                     DEFAULT '1970-01-01 00:00:00',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `uk_username` (`username`, `deleted_at`),
    UNIQUE KEY `uk_email` (`email`, `deleted_at`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 9
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_user"
--

INSERT INTO `tbl_user`
VALUES (1, 'string', 'aoihosizora@hotmail.com', 'string', 1, 'lll', 'https://aaa.bbb.ccc', '1999-12-20', 'root', 1,
        '13512345678', '2020-09-01 10:36:21', '2020-09-01 11:02:06', '1970-01-01 00:00:00'),
       (2, 'string2', 'a970335605@gmail.com', 'string2', 2, 'string', 'https://aaa.bbb.ccc', '2000-01-01', 'sub_root',
        1, '13512345678', '2020-09-01 10:36:42', '2020-09-01 11:24:45', '1970-01-01 00:00:00'),
       (3, 'string3', 'aaa3@bbb.ccc', 'string3', 1, 'string', 'https://aaa.bbb.ccc', '2020-09-01', 'normal', 0,
        '13512345678', '2020-09-01 10:36:51', '2020-09-01 11:25:49', '1970-01-01 00:00:00'),
       (4, 'string4', 'aaa4@bbb.ccc', 'user20200901103718309483', 0, '', '', '2000-01-01', 'normal', 0, '',
        '2020-09-01 10:37:18', '2020-09-01 10:37:18', '1970-01-01 00:00:00'),
       (5, 'string5', 'aaa5@bbb.ccc', 'user20200901103720772571', 0, '', '', '2000-01-01', 'normal', 0, '',
        '2020-09-01 10:37:21', '2020-09-01 10:37:21', '1970-01-01 00:00:00'),
       (6, 'string6', 'aaa6@bbb.ccc', 'user20200901103726087155', 0, '', '', '2000-01-01', 'normal', 0, '',
        '2020-09-01 10:37:26', '2020-09-01 10:37:26', '1970-01-01 00:00:00'),
       (7, 'string7', 'aaa7@bbb.ccc', 'user20200901103729785780', 0, '', '', '2000-01-01', 'normal', 0, '',
        '2020-09-01 10:37:30', '2020-09-01 10:37:30', '1970-01-01 00:00:00'),
       (8, 'string8', 'aaa8@bbb.ccc', 'user20200901103732341677', 0, '', '', '2000-01-01', 'normal', 0, '',
        '2020-09-01 10:37:32', '2020-09-01 10:37:32', '1970-01-01 00:00:00');

--
-- Structure for table "tbl_video"
--

CREATE TABLE `tbl_video`
(
    `vid`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `title`       varchar(255)        NOT NULL,
    `description` varchar(1023)       NOT NULL,
    `video_url`   varchar(255)        NOT NULL,
    `cover_url`   varchar(255)        NOT NULL,
    `author_uid`  bigint(20) unsigned NOT NULL,
    `created_at`  datetime DEFAULT NULL,
    `updated_at`  datetime DEFAULT NULL,
    `deleted_at`  datetime DEFAULT '1970-01-01 00:00:00',
    PRIMARY KEY (`vid`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  DEFAULT CHARSET = utf8;

--
-- Data for table "tbl_video"
--

INSERT INTO `tbl_video`
VALUES (1, 'string', 'string', 'https://aaa.bbb.ccc', 'https://aaa.bbb.ccc', 1, '2020-09-01 11:10:18',
        '2020-09-01 11:10:18', '1970-01-01 00:00:00'),
       (2, 'string2', 'string', 'https://aaa.bbb.ccc', 'https://aaa.bbb.ccc', 1, '2020-09-01 11:10:21',
        '2020-09-01 11:11:39', '1970-01-01 00:00:00'),
       (3, 'string3', 'string', 'https://aaa.bbb.ccc', 'https://aaa.bbb.ccc', 1, '2020-09-01 11:10:22',
        '2020-09-01 11:11:45', '2020-09-01 11:12:02'),
       (4, 'string4', 'string', 'https://aaa.bbb.ccc', 'https://aaa.bbb.ccc', 2, '2020-09-01 11:10:22',
        '2020-09-01 11:11:48', '1970-01-01 00:00:00'),
       (5, 'string5', 'string', 'https://aaa.bbb.ccc', 'https://aaa.bbb.ccc', 3, '2020-09-01 11:10:23',
        '2020-09-01 11:11:52', '1970-01-01 00:00:00');
