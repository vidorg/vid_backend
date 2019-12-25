# vid backend
Backend of repo https://github.com/vidorg/vid_vue

## Version: 1.1

### Terms of service
https://github.com/vidorg

**License:** [MIT](https://github.com/vidorg/vid_backend/blob/master/LICENSE)

### /auth/

#### GET
##### Summary:

查看当前登录用户

##### Description:

根据认证令牌，查看当前登录用户

| Code | Message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /auth/login

#### POST
##### Summary:

登录

##### Description:

用户登录

| Code | Message |
| --- | --- |
| 400 | request form data error |
| 401 | password error |
| 404 | user not found |
| 500 | login failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名 | Yes | string |
| password | formData | 用户密码 | Yes | string |
| expire | formData | 登录有效期，默认为七天 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "user": {       "uid": 10,       "username": "aoihosizora",       "sex": "unknown",       "profile": "",       "avatar_url": "",       "birth_time": "2000-01-01",       "authority": "normal"     },     "token": "Bearer xxx"   } } ``` |

### /auth/pass

#### POST
##### Summary:

修改密码

##### Description:

用户修改密码

| Code | Message |
| --- | --- |
| 400 | request form data error |
| 400 | request format error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 500 | update password failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| password | formData | 用户密码  *Minimum length* : 8 *Maximin length* : 30 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /auth/register

#### POST
##### Summary:

注册

##### Description:

用户注册

| Code | Message |
| --- | --- |
| 400 | request form data error |
| 400 | request format error |
| 500 | username duplicated |
| 500 | register failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名  *Minimum length* : 5 *Maximin length* : 30 | Yes | string |
| password | formData | 用户密码  *Minimum length* : 8 *Maximin length* : 30 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

### /user/

#### DELETE
##### Summary:

删除用户

##### Description:

删除用户所有信息

| Code | Message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 404 | user delete failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success" } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

#### PUT
##### Summary:

更新用户

##### Description:

更新用户信息

| Code | Message |
| --- | --- |
| 400 | request format error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 500 | username duplicated |
| 500 | user update failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| username | formData | 用户名  *Minimum length* : 8 *Maximin length* : 30 | No | string |
| sex | formData | 用户性别 | No | string |
| profile | formData | 用户简介  *Minimum length* : 0 *Maximin length* : 255 | No | string |
| birth_time | formData | 用户生日，固定格式为2000-01-01 | No | string |
| phone_number | formData | 用户手机号码 | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "male",     "profile": "Demo Profile",     "avatar_url": "",     "birth_time": "2019-12-18",     "authority": "admin"   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /user/sub?uid

#### POST
##### Summary:

关注用户

##### Description:

关注某一用户

| Code | Message |
| --- | --- |
| 400 | request query param error |
| 400 | subscribe oneself invalid |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 500 | subscribe failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| uid | query | 关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "me": 10,     "up": 3,     "action": "subscribe"   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /user/unsub?uid

#### POST
##### Summary:

取消关注用户

##### Description:

取消关注某一用户

| Code | Message |
| --- | --- |
| 400 | request query param error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 500 | unsubscribe failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| uid | query | 取消关注用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "me": 10,     "up": 3,     "action": "unsubscribe"   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /user/{uid}

#### GET
##### Summary:

查询用户

##### Description:

查询用户信息

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "user": {       "uid": 10,       "username": "aoihosizora",       "sex": "unknown",       "profile": "",       "avatar_url": "",       "birth_time": "2000-01-01",       "authority": "admin"     },     "extra": {       "subscribing_cnt": 1,       "subscriber_cnt": 2,       "video_cnt": 0,       "playlist_cnt": 0     }   } } ``` |

### /user/{uid}/subscriber

#### GET
##### Summary:

用户粉丝

##### Description:

查询用户所有粉丝，返回分页数据

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 查询的用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "count": 1,     "page": 1,     "data": [       {         "uid": 2,         "username": "User2",         "sex": "unknown",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "normal"       }     ]   } } ``` |

### /user/{uid}/subscribing

#### GET
##### Summary:

用户关注

##### Description:

查询用户所有关注，返回分页数据

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 查询的用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "count": 1,     "page": 1,     "data": [       {         "uid": 1,         "username": "User1",         "sex": "male",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "admin"       }     ]   } } ``` |

### /user?page

#### GET
##### Summary:

查询所有用户

##### Description:

管理员查询所有用户，返回分页数据，Admin

| Code | Message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 401 | need admin authority |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "count": 1,     "page": 1,     "data": [       {         "uid": 1,         "username": "User1",         "sex": "male",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "admin"       }     ]   } } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /video/

#### POST
##### Summary:

新建视频

##### Description:

新建用户视频

| Code | Message |
| --- | --- |
| 400 | request form data error |
| 400 | request format error |
| 401 | authorization failed |
| 401 | token has expired |
| 500 | video existed failed |
| 500 | video insert failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| title | formData | 视频标题  *Minimum length* : 5 *Maximin length* : 100 | Yes | string |
| description | formData | 视频简介  *Minimum length* : 0 *Maximin length* : 255 | Yes | string |
| video_url | formData | 视频资源链接 | Yes | string |
| cover_url | formData | 视频封面链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /video/uid/{uid}?page

#### GET
##### Summary:

查询用户视频

##### Description:

查询作者为用户的所有视频，返回分页数据

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

### /video/vid/{vid}

#### GET
##### Summary:

查询视频

##### Description:

查询视频信息

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 404 | video not found |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 视频id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

### /video/{vid}

#### DELETE
##### Summary:

删除视频

##### Description:

删除用户视频

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | video not found |
| 500 | video delete failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| vid | path | 删除视频id | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

#### POST
##### Summary:

更新视频

##### Description:

更新用户视频信息

| Code | Message |
| --- | --- |
| 400 | request route param error |
| 400 | request format error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | video not found |
| 500 | video update failed |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| vid | path | 更新视频id | Yes | string |
| title | formData | 视频标题  *Minimum length* : 5 *Maximin length* : 100 | No | string |
| description | formData | 视频简介  *Minimum length* : 0 *Maximin length* : 255 | No | string |
| video_url | formData | 视频资源链接 | No | string |
| cover_url | formData | 视频封面链接 | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |

### /video?page

#### GET
##### Summary:

查询所有视频

##### Description:

管理员查询所有视频，返回分页数据，Admin

| Code | Message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 401 | need admin authority |




##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户登录令牌 | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

##### Security

| Security Schema | Scopes | |
| --- | --- | --- |
| basicAuth | [] | |
