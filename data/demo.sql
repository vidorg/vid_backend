-- Host: localhost  (Version 8.0.15)
-- Date: 2020-01-10 01:18:05
-- Generator: MySQL-Front 6.0  (Build 2.20)

-- db: db_vid
-- host: localhost:3306
-- charset: utf8


-- tbl: tbl_account
-- passwords are all 12345678

INSERT INTO `tbl_account` (`uid`, `encrypted_pass`, `created_at`, `updated_at`, `deleted_at`)
VALUES (1, '$2a$10$iq6j.QEU0Nvj5EmZtj2fsu4WTJGNdM5rJBSn0.qbafxc8ZYXviq/y', '2020-01-10 00:30:49', '2020-01-10 00:47:02',
        '2000-01-01 00:00:00'),
       (2, '$2a$10$NSV3DOeIugCPLxJYkpjdLeNJJ/EeOUlEqtlbl9/nz4q1tPCXIU63y', '2020-01-10 00:31:04', '2020-01-10 00:31:04',
        '2000-01-01 00:00:00'),
       (3, '$2a$10$HA90TCha0MfxFsIH7i8h3.Lu2fXG/wC3LO6.wKwNzHKq5F9oK3XWW', '2020-01-10 00:31:08', '2020-01-10 00:31:08',
        '2000-01-01 00:00:00'),
       (4, '$2a$10$1qpTH7aI19n/FdiU1bfJC.Xd6sWU.ESIBoYgWlTbwVMCgysR9g8mm', '2020-01-10 00:31:11', '2020-01-10 00:31:11',
        '2000-01-01 00:00:00'),
       (5, '$2a$10$HGDcYm6edBsIBV/mP5bDXuYyVYSAEdI0jGAn5sXVMrvMuSEUaEhQm', '2020-01-10 00:31:15', '2020-01-10 00:31:15',
        '2000-01-01 00:00:00'),
       (6, '$2a$10$TQ/.siAIhrrn4O/W13YjNe9LYnZ55qMabESqiyrFTi.K8tzg6tz32', '2020-01-10 00:31:18', '2020-01-10 00:31:18',
        '2020-01-10 00:41:48');


-- tbl: tbl_user

INSERT INTO `tbl_user` (`uid`, `username`, `sex`, `profile`, `avatar_url`, `birthday`, `authority`, `register_ip`,
                        `phone_number`, `created_at`, `updated_at`, `deleted_at`)
VALUES (1, 'admin', 'male', 'Demo admin profile', '', '2020-01-10', 'admin', '::1', '13512345678',
        '2020-01-10 00:30:49', '2020-01-10 01:16:52', '2000-01-01 00:00:00'),
       (2, 'testuser1', 'unknown', 'Demo profile', '', '2000-01-01', 'normal', '::1', '13312345678',
        '2020-01-10 00:31:04', '2020-01-10 00:31:04', '2000-01-01 00:00:00'),
       (3, 'testuser2', 'male', '', '', '2000-01-01', 'normal', '::1', '', '2020-01-10 00:31:08',
        '2020-01-10 00:31:08', '2000-01-01 00:00:00'),
       (4, 'testuser3', 'female', '', '', '2000-01-01', 'normal', '::1', '', '2020-01-10 00:31:11',
        '2020-01-10 00:31:11', '2000-01-01 00:00:00'),
       (5, 'testuser4', 'unknown', '', '', '2000-01-01', 'normal', '::1', '', '2020-01-10 00:31:15',
        '2020-01-10 00:31:15', '2000-01-01 00:00:00'),
       (6, 'testuser_del', 'male', 'Demo deleted user', '', '1999-01-31', 'normal', '::1', '15912345678',
        '2020-01-10 00:31:18', '2020-01-10 00:39:09', '2020-01-10 00:41:48');


-- tbl: tbl_subscribe

INSERT INTO `tbl_subscribe` (`subscriber_uid`, `up_uid`)
VALUES (1, 2),
       (1, 3),
       (2, 1),
       (2, 3),
       (2, 4),
       (2, 5),
       (3, 2),
       (3, 6),
       (5, 1),
       (5, 4);


-- tbl: tbl_video

INSERT INTO `tbl_video` (`vid`, `title`, `description`, `video_url`, `cover_url`, `upload_time`, `author_uid`,
                         `created_at`, `updated_at`, `deleted_at`)
VALUES (1, 'The First Video', 'This is the first video uploaded', '123', '',  1, '2020-01-10 00:55:36',  '2020-01-10 00:55:36', '2000-01-01 00:00:00'),
       (2, 'The Second Video', 'This is the second video uploaded', '123', '',  2, '2020-01-10 00:55:51', '2020-01-10 00:55:51', '2000-01-01 00:00:00'),
       (3, 'Test Video', 'Demo Desc', '12345', '', '2020-01-10 00:56:30', 1, '2020-01-10 01:16:52', '2000-01-01 00:00:00'),
       (4, 'Test Video2', '', '123', '', 2, '2020-01-10 01:11:09', '2020-01-10 01:11:09', '2020-01-10 01:12:23'),
       (5, 'Test Video3', '', '123', '', 4, '2020-01-10 01:11:52', '2020-01-10 01:11:52', '2000-01-01 00:00:00'),
       (6, 'Test Video4', '', '123', '', 3, '2020-01-10 01:11:53', '2020-01-10 01:11:53', '2000-01-01 00:00:00'),
       (7, 'Test Video5', '', '123', '', 4, '2020-01-10 01:11:54', '2020-01-10 01:11:54', '2000-01-01 00:00:00'),
       (8, 'Test Video6', '', '123', '', 2, '2020-01-10 01:11:55', '2020-01-10 01:11:55', '2000-01-01 00:00:00'),
       (9, 'Test Video7', '', '123', '', 4, '2020-01-10 01:11:56', '2020-01-10 01:11:56', '2000-01-01 00:00:00'),
       (10, 'Test Video8', '', '123', '', 1, '2020-01-10 01:11:56', '2020-01-10 01:11:56', '2020-01-10 01:17:11');
