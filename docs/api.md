# vid backend
Backend of repo https://github.com/vidorg/vid_vue

## Version: 1.1

### Terms of service
https://github.com/vidorg

**License:** [MIT](https://github.com/vidorg/vid_backend/blob/master/LICENSE)

### Security
**Jwt**  

|apiKey|*API Key*|
|---|---|
|In|header|
|Name|Authorization|

### /ping

#### GET
##### Summary:

Ping

##### Description:

Ping

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json Content-Type: application/json; charset=utf-8 ``` |

### /v1/auth/

#### GET
##### Summary:

当前登录用户

##### Description:

根据认证信息，查看当前登录用户

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 1,         "username": "admin",         "sex": "male",         "profile": "Demo admin profile",         "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "birth_time": "2020-01-10",         "authority": "admin",         "phone_number": "13512345678",         "register_time": "2020-01-10 00:30:49"     } } ``` |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

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
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "user": {             "uid": 1,             "username": "admin",             "sex": "male",             "profile": "Demo admin profile",             "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",             "birth_time": "2020-01-10",             "authority": "admin",             "phone_number": "13512345678",             "register_time": "2020-01-10 00:30:49"         },         "token": "Bearer xxx",         "expire": 604800     } } ``` |
| 400 | "request param error" |
| 401 | "password error" |
| 404 | "user not found" |
| 500 | "login failed" |

### /v1/auth/logout

#### POST
##### Summary:

注销

##### Description:

用户注销

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 500 | "logout failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/auth/password

#### PUT
##### Summary:

修改密码

##### Description:

用户修改密码

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| password | formData | 用户密码，长度在 [8, 30] 之间 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error", "request format error" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 404 | "user not found" |
| 500 | "update password failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/auth/register

#### POST
##### Summary:

注册

##### Description:

注册新用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名，长度在 [5, 30] 之间 | Yes | string |
| password | formData | 用户密码，长度在 [8, 30] 之间 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 201,     "message": "created",     "data": {         "uid": 1,         "username": "admin",         "sex": "male",         "profile": "Demo admin profile",         "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "birth_time": "2020-01-10",         "authority": "admin",         "phone_number": "13512345678",         "register_time": "2020-01-10 00:30:49"     } } ``` |
| 400 | "request param error", "request format error" |
| 500 | "username has been used", "register failed" |

### /v1/raw/image

#### POST
##### Summary:

上传图片

##### Description:

上传公共图片，包括用户头像和视频封面

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| image | formData | 上传的图片，大小限制在2M，允许后缀名为 {.jpg, .jpeg, .png, .bmp, .gif} | Yes | file |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "url": "http://localhost:3344/v1/raw/image/20200110130323908439.jpg",         "size": 381952     } } ``` |
| 400 | "request param error", "image type not supported" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 413 | "request body too large" |
| 500 | "image save failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/raw/image/{filename}

#### GET
##### Summary:

获取图片

##### Description:

获取用户头像图片以及视频封面

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| filename | path | 图片文件名，jpg后缀名 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json Content-Type: image/jpeg ``` |
| 404 | "image not found" |

### /v1/user/

#### DELETE
##### Summary:

删除用户

##### Description:

删除用户账户以及所有信息

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 404 | "user not found" |
| 500 | "user delete failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

#### PUT
##### Summary:

更新用户

##### Description:

更新用户个人信息

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名，长度在 [8, 30] 之间 | Yes | string |
| sex | formData | 用户性别，允许值为 {male, female, unknown} | Yes | string |
| profile | formData | 用户简介，长度在 [0, 255] 之间 | Yes | string |
| birth_time | formData | 用户生日，固定格式为 2000-01-01 | Yes | string |
| phone_number | formData | 用户手机号码，长度为 11，仅限中国大陆手机号码 | Yes | string |
| avatar_url | formData | 用户头像链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 1,         "username": "admin",         "sex": "male",         "profile": "Demo admin profile",         "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "birth_time": "2020-01-10",         "authority": "admin",         "phone_number": "13512345678",         "register_time": "2020-01-10 00:30:49"     } } ``` |
| 400 | "request param error", "request format error", "username has been used" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 404 | "user not found" |
| 500 | "user update failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/user/admin/{uid}

#### DELETE
##### Summary:

管理员删除用户

##### Description:

删除用户账户，管理员权限

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |
| 404 | "user not found" |
| 500 | "user delete failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

#### PUT
##### Summary:

管理员更新用户

##### Description:

更新用户信息，管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名，长度在 [8, 30] 之间 | Yes | string |
| sex | formData | 用户性别，允许值为 {male, female, unknown} | Yes | string |
| profile | formData | 用户简介，长度在 [0, 255] 之间 | Yes | string |
| birth_time | formData | 用户生日，固定格式为 2000-01-01 | Yes | string |
| phone_number | formData | 用户手机号码，长度为 11，仅限中国大陆手机号码 | Yes | string |
| avatar_url | formData | 用户头像链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "uid": 1,         "username": "admin",         "sex": "male",         "profile": "Demo admin profile",         "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "birth_time": "2020-01-10",         "authority": "admin",         "phone_number": "13512345678",         "register_time": "2020-01-10 00:30:49"     } } ``` |
| 400 | "request param error", "request format error", "username has been used" |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |
| 404 | "user not found" |
| 500 | "user update failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/user/subscribing

#### DELETE
##### Summary:

取消关注

##### Description:

取消关注某一用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| to | formData | 取消关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error", "request format error" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 404 | "user not found" |
| 500 | "unsubscribe failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

#### PUT
##### Summary:

关注

##### Description:

关注某一用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| to | formData | 关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error", "request format error", "subscribe oneself invalid" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 404 | "user not found" |
| 500 | "subscribe failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/user/{uid}

#### GET
##### Summary:

查询用户

##### Description:

查询用户个人信息和数量信息，此处可见用户手机号码

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "user": {             "uid": 1,             "username": "admin",             "sex": "male",             "profile": "Demo admin profile",             "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",             "birth_time": "2020-01-10",             "authority": "admin",             "phone_number": "13512345678",             "register_time": "2020-01-10 00:30:49"         },         "extra": {             "subscribing_cnt": 3,             "subscriber_cnt": 2,             "video_cnt": 3         }     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/subscriber?page

#### GET
##### Summary:

查询粉丝

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
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "uid": 1,                 "username": "admin",                 "sex": "male",                 "profile": "Demo admin profile",                 "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",                 "birth_time": "2020-01-10",                 "authority": "admin",                 "phone_number": "13512345678",                 "register_time": "2020-01-10 00:30:49"             }         ]     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/subscribing?page

#### GET
##### Summary:

查询关注

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
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user/{uid}/video?page

#### GET
##### Summary:

查询用户视频

##### Description:

查询作者为指定用户的所有视频，返回分页数据

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "vid": 1,                 "title": "The First Video",                 "description": "This is the first video uploaded",                 "video_url": "123",                 "cover_url": "http://localhost:3344/v1/raw/image/avatar.jpg",                 "upload_time": "2020-01-10 00:55:36",                 "update_time": "2020-01-10 14:31:00",                 "author": {                     "uid": 1,                     "username": "admin",                     "sex": "male",                     "profile": "Demo admin profile",                     "avatar_url": "http://localhost:3344/v1/raw/image/cover.jpg",                     "birth_time": "2020-01-10",                     "authority": "admin",                     "register_time": "2020-01-10 00:30:49"                 }             }         ]     } } ``` |
| 400 | "request param error" |
| 404 | "user not found" |

### /v1/user?page

#### GET
##### Summary:

查询所有用户

##### Description:

管理员查询所有用户，返回分页数据，管理员权限，此处可见用户手机号码

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "uid": 1,                 "username": "admin",                 "sex": "male",                 "profile": "Demo admin profile",                 "avatar_url": "http://localhost:3344/v1/raw/image/avatar.jpg",                 "birth_time": "2020-01-10",                 "authority": "admin",                 "phone_number": "13512345678",                 "register_time": "2020-01-10 00:30:49"             }         ]     } } ``` |
| 400 | "request param error" |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/video/

#### POST
##### Summary:

新建视频

##### Description:

新建用户视频

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| title | formData | 视频标题，长度在 [1, 100] 之间 | Yes | string |
| description | formData | 视频简介，长度在 [0, 1024] 之间 | Yes | string |
| cover_url | formData | 视频封面链接 | Yes | string |
| video_url | formData | 视频资源链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 201,     "message": "created",     "data": {         "vid": 1,         "title": "The First Video",         "description": "This is the first video uploaded",         "video_url": "123",         "cover_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "upload_time": "2020-01-10 00:55:36",         "update_time": "2020-01-10 14:31:00",         "author": {             "uid": 1,             "username": "admin",             "sex": "male",             "profile": "Demo admin profile",             "avatar_url": "http://localhost:3344/v1/raw/image/cover.jpg",             "birth_time": "2020-01-10",             "authority": "admin",             "register_time": "2020-01-10 00:30:49"         }     } } ``` |
| 400 | "request param error", "request format error", "video has been updated" |
| 401 | "unauthorized user", "token has expired", "authorized user not found" |
| 500 | "video insert failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/video/{vid}

#### DELETE
##### Summary:

删除视频

##### Description:

删除用户视频，管理员或者作者本人可以操作

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 删除视频id | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success" } ``` |
| 400 | "request param error" |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |
| 404 | "video not found" |
| 500 | "video delete failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

#### GET
##### Summary:

查询视频

##### Description:

查询视频信息，作者id为-1表示已删除的用户

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 视频id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "vid": 1,         "title": "The First Video",         "description": "This is the first video uploaded",         "video_url": "123",         "cover_url": "http://localhost:3344/v1/raw/image/avatar.jpg",         "upload_time": "2020-01-10 00:55:36",         "update_time": "2020-01-10 14:31:00",         "author": {             "uid": 1,             "username": "admin",             "sex": "male",             "profile": "Demo admin profile",             "avatar_url": "http://localhost:3344/v1/raw/image/cover.jpg",             "birth_time": "2020-01-10",             "authority": "admin",             "register_time": "2020-01-10 00:30:49"         }     } } ``` |
| 400 | "request param error" |
| 404 | "video not found" |

#### POST
##### Summary:

更新视频

##### Description:

更新用户视频信息，管理员或者作者本人可以操作

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 更新视频id | Yes | string |
| title | formData | 视频标题，长度在 [1, 100] 之间 | Yes | string |
| description | formData | 视频简介，长度在 [0, 1024] 之间 | Yes | string |
| cover_url | formData | 视频封面链接 | Yes | string |
| video_url | formData | 视频资源链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 400 | "request param error", "request format error", "video has been updated" |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |
| 404 | "video not found" |
| 500 | "video update failed" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |

### /v1/video?page

#### GET
##### Summary:

查询所有视频

##### Description:

管理员查询所有视频，返回分页数据，管理员权限

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {     "code": 200,     "message": "success",     "data": {         "count": 1,         "page": 1,         "data": [             {                 "vid": 1,                 "title": "The First Video",                 "description": "This is the first video uploaded",                 "video_url": "123",                 "cover_url": "http://localhost:3344/v1/raw/image/avatar.jpg",                 "upload_time": "2020-01-10 00:55:36",                 "update_time": "2020-01-10 14:31:00",                 "author": {                     "uid": 1,                     "username": "admin",                     "sex": "male",                     "profile": "Demo admin profile",                     "avatar_url": "http://localhost:3344/v1/raw/image/cover.jpg",                     "birth_time": "2020-01-10",                     "authority": "admin",                     "register_time": "2020-01-10 00:30:49"                 }             }         ]     } } ``` |
| 400 | "request param error" |
| 401 | "unauthorized user", "token has expired", "authorized user not found", "need admin authority" |

##### Security

| Security Schema | Scopes |
| --- | --- |
| Jwt | |
