# vid backend
Backend of repo https://github.com/vidorg/vid_vue

## Version: 1.1

### Terms of service
https://github.com/vidorg

**License:** [MIT](https://github.com/vidorg/vid_backend/blob/master/LICENSE)

### /ping

#### GET
##### Summary:

Ping

##### Description:

Ping

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "ping": "pong" } ``` |

### /v1/auth/

#### GET
##### Summary:

查看当前登录用户

##### Description:

根据认证令牌，查看当前登录用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 10,         "username": "aoihosizora",         "sex": "male",         "profile": "Demo profile",         "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",         "birth_time": "2001-02-03",         "authority": "normal",         "phone_number": "13512345678"     } } ``` |
| 401 | "authorization failed" / "token has expired" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/auth/login

#### POST
##### Summary:

登录

##### Description:

用户登录

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名 | Yes | string |
| password | formData | 用户密码 | Yes | string |
| expire | formData | 登录有效期，默认为七天 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "user": {             "uid": 10,             "username": "aoihosizora",             "sex": "male",             "profile": "Demo profile",             "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",             "birth_time": "2001-02-03",             "authority": "normal",             "phone_number": "13512345678"         },         "token": "Bearer xxx",         "expire": 604800     } } ``` |
| 400 | "request param error" |
| 401 | "password error" |
| 404 | "user not found" |
| 500 | "login failed" |

### /v1/auth/password

#### PUT
##### Summary:

修改密码

##### Description:

用户修改密码

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| password | formData | 用户密码 *Minimum length* : 8 *Maximum length* : 30 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error" / "request format error" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "user not found" |
| 500 | "update password failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/auth/register

#### POST
##### Summary:

注册

##### Description:

用户注册

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名 *Minimum length* : 5 *Maximum length* : 30 | Yes | string |
| password | formData | 用户密码 *Minimum length* : 8 *Maximum length* : 30 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 10,         "username": "aoihosizora",         "sex": "male",         "profile": "Demo profile",         "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",         "birth_time": "2001-02-03",         "authority": "normal",         "phone_number": "13512345678"     } } ``` |
| 400 | "request param error" / "request format error" |
| 500 | "username has been used" / "register failed" |

### /v1/raw/image/{uid}/{filename}

#### GET
##### Summary:

获取图片

##### Description:

获取用户头像图片以及视频封面

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id，或者default | Yes | string |
| filename | path | 图片文件名，jpg后缀名 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "Content-Type": "image/jpeg" } ``` |
| 400 | "request param error" |
| 404 | "image not found" |

### /v1/user/

#### DELETE
##### Summary:

删除用户

##### Description:

删除用户所有信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 401 | "authorization failed" / "token has expired" |
| 404 | "user not found" |
| 500 | "user delete failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

#### PUT
##### Summary:

更新用户

##### Description:

更新用户信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| username | formData | 用户名 *Minimum length* : 8 *Maximum length* : 30 | Yes | string |
| sex | formData | 用户性别 | Yes | string |
| profile | formData | 用户简介 *Minimum length* : 0 *Maximum length* : 255 | Yes | string |
| birth_time | formData | 用户生日，固定格式为2000-01-01 | Yes | string |
| phone_number | formData | 用户手机号码 | Yes | string |
| avatar | formData | 用户头像，默认不修改 | No | file |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 10,         "username": "aoihosizora",         "sex": "male",         "profile": "Demo profile",         "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",         "birth_time": "2001-02-03",         "authority": "normal",         "phone_number": "13512345678"     } } ``` |
| 400 | "request param error" / "request format error" / "request body too large" / "username has been used" / "image type not supported" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "user not found" |
| 500 | "image save failed" / "user update failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/user/subscribing

#### DELETE
##### Summary:

取消关注用户

##### Description:

取消关注某一用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| to | formData | 取消关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "me_uid": 10,         "to_uid": 3,         "action": "unsubscribe"     } } ``` |
| 400 | "request param error" / "request format error" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "user not found" |
| 500 | "unsubscribe failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

#### PUT
##### Summary:

关注用户

##### Description:

关注某一用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| to | formData | 关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "me_uid": 10,         "to_uid": 3,         "action": "subscribe"     } } ``` |
| 400 | "request param error" / "request format error" / "subscribe oneself invalid" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "user not found" |
| 500 | "subscribe failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/user/{uid}

#### GET
##### Summary:

查询用户

##### Description:

查询用户信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "user": {             "uid": 10,             "username": "aoihosizora",             "sex": "male",             "profile": "Demo profile",             "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",             "birth_time": "2001-02-03",             "authority": "normal",             "phone_number": "13512345678"         },         "extra": {             "subscribing_cnt": 3,             "subscriber_cnt": 2,             "video_cnt": 3         }     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/subscriber

#### GET
##### Summary:

用户粉丝

##### Description:

查询用户所有粉丝，返回分页数据

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 查询的用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "uid": 10,                 "username": "aoihosizora",                 "sex": "male",                 "profile": "Demo profile",                 "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",                 "birth_time": "2001-02-03",                 "authority": "normal",                 "phone_number": "13512345678"             }         ]     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/subscribing

#### GET
##### Summary:

用户关注

##### Description:

查询用户所有关注，返回分页数据

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 查询的用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "uid": 10,                 "username": "aoihosizora",                 "sex": "male",                 "profile": "Demo profile",                 "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",                 "birth_time": "2001-02-03",                 "authority": "normal",                 "phone_number": "13512345678"             }         ]     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/video?page

#### GET
##### Summary:

查询用户视频

##### Description:

查询作者为用户的所有视频，返回分页数据

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "vid": 10,                 "title": "Video Title",                 "description": "Video Demo Description",                 "video_url": "",                 "cover_url": "http://localhost:3344/raw/image/default/cover.jpg",                 "upload_time": "2019-12-26 14:14:04",                 "update_time": "2019-12-30 21:04:51",                 "author": {                     "uid": 10,                     "username": "aoihosizora",                     "sex": "male",                     "profile": "Demo profile",                     "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",                     "birth_time": "2001-02-03",                     "authority": "normal"                 }             }         ]     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user?page

#### GET
##### Summary:

查询所有用户

##### Description:

管理员查询所有用户，返回分页数据，Admin

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "uid": 10,                 "username": "aoihosizora",                 "sex": "male",                 "profile": "Demo profile",                 "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",                 "birth_time": "2001-02-03",                 "authority": "normal",                 "phone_number": "13512345678"             }         ]     } } ``` |
| 400 | "request param error" |
| 401 | "authorization failed" / "token has expired" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/video/

#### POST
##### Summary:

新建视频

##### Description:

新建用户视频

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| title | formData | 视频标题 *Minimum length* : 5 *Maximum length* : 100 | Yes | string |
| description | formData | 视频简介 *Minimum length* : 0 *Maximum length* : 255 | Yes | string |
| video_url | formData | 视频资源链接 | Yes | string |
| cover | formData | 视频封面 | No | file |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {        "code": 200,        "message": "success",        "data": {"vid": 10, "title": "Video Title", "description": "Video Demo Description", "video_url": "", "cover_url": "http://localhost:3344/raw/image/default/cover.jpg", "upload_time": "2019-12-26 14:14:04", "update_time": "2019-12-30 21:04:51", "author": {"uid": 10, "username": "aoihosizora", "sex": "male", "profile": "Demo profile", "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg", "birth_time": "2001-02-03", "authority": "normal"}}        }        } ``` |
| 400 | "request param error" / "request format error" / "request body too large" / "image type not supported" / "video resource has been used" |
| 401 | "authorization failed" / "token has expired" |
| 500 | "image save failed" / "video insert failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/video/{vid}

#### DELETE
##### Summary:

删除视频

##### Description:

删除用户视频

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| vid | path | 删除视频id | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "video not found" |
| 500 | "video delete failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

#### GET
##### Summary:

查询视频

##### Description:

查询视频信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 视频id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {        "code": 200,        "message": "success",        "data": {"vid": 10, "title": "Video Title", "description": "Video Demo Description", "video_url": "", "cover_url": "http://localhost:3344/raw/image/default/cover.jpg", "upload_time": "2019-12-26 14:14:04", "update_time": "2019-12-30 21:04:51", "author": {"uid": 10, "username": "aoihosizora", "sex": "male", "profile": "Demo profile", "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg", "birth_time": "2001-02-03", "authority": "normal"}}        }        } ``` |
| 400 | "request param error" |
| 404 | "video not found" |

#### POST
##### Summary:

更新视频

##### Description:

更新用户视频信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| vid | path | 更新视频id | Yes | string |
| title | formData | 视频标题 *Minimum length* : 5 *Maximum length* : 100 | No | string |
| description | formData | 视频简介 *Minimum length* : 0 *Maximum length* : 255 | No | string |
| cover | formData | 视频封面 | No | file |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "vid": 10,         "title": "Video Title",         "description": "Video Demo Description",         "video_url": "",         "cover_url": "http://localhost:3344/raw/image/default/cover.jpg",         "upload_time": "2019-12-26 14:14:04",         "update_time": "2019-12-30 21:04:51",         "author": {             "uid": 10,             "username": "aoihosizora",             "sex": "male",             "profile": "Demo profile",             "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",             "birth_time": "2001-02-03",             "authority": "normal"         }     } } ``` |
| 400 | "request param error" / "request format error" / "request body too large" / "image type not supported" |
| 401 | "authorization failed" / "token has expired" |
| 404 | "video not found" |
| 500 | "image save failed" / "video update failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |

### /v1/video?page

#### GET
##### Summary:

查询所有视频

##### Description:

管理员查询所有视频，返回分页数据，Admin

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "vid": 10,                 "title": "Video Title",                 "description": "Video Demo Description",                 "video_url": "",                 "cover_url": "http://localhost:3344/raw/image/default/cover.jpg",                 "upload_time": "2019-12-26 14:14:04",                 "update_time": "2019-12-30 21:04:51",                 "author": {                     "uid": 10,                     "username": "aoihosizora",                     "sex": "male",                     "profile": "Demo profile",                     "avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",                     "birth_time": "2001-02-03",                     "authority": "normal"                 }             }         ]     } } ``` |
| 400 | "request param error" |
| 401 | "authorization failed" / "token has expired" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| basicAuth | |
