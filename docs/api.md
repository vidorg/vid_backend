# vid backend
Backend of repo https://github.com/vidorg/vid_vue

## Version: 1.1

### Terms of service
https://github.com/vidorg

**License:** [MIT](https://github.com/vidorg/vid_backend/blob/master/LICENSE)

### /auth/

#### GET
##### Summary:

查看当前用户

##### Description:

根据认证 token 查看当前用户，Auth

| code | message |
| --- | --- |
| 400 | request form data error |
| 401 | authorization failed |
| 401 | token has expired | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

### /auth/login

#### POST
##### Summary:

登录

##### Description:

用户登录，Non-Auth

| code | message |
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
| expire | formData | 登录有效期，默认一个小时 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "user": {       "uid": 10,       "username": "aoihosizora",       "sex": "unknown",       "profile": "",       "avatar_url": "",       "birth_time": "2000-01-01",       "authority": "normal"     },     "token": "Bearer xxx"   } } ``` |

### /auth/pass

#### POST
##### Summary:

修改密码

##### Description:

用户修改密码，Auth

| code | message |
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
| Authorization | header | 用户 Token | Yes | string |
| password | formData | 用户新密码 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

### /auth/register

#### POST
##### Summary:

注册

##### Description:

用户注册，Non-Auth

| code | message |
| --- | --- |
| 400 | request form data error |
| 400 | request format error |
| 500 | username duplicated |
| 500 | register failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | formData | 用户名 | Yes | string |
| password | formData | 用户密码 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "unknown",     "profile": "",     "avatar_url": "",     "birth_time": "2000-01-01",     "authority": "normal"   } } ``` |

### /user/

#### DELETE
##### Summary:

删除用户

##### Description:

删除用户所有信息，Auth

| code | message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 404 | user delete failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success" } ``` |

#### PUT
##### Summary:

更新用户

##### Description:

更新用户信息，Auth

| code | message |
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
| Authorization | header | 用户 Token | Yes | string |
| username | formData | 新用户名 | No | string |
| sex | formData | 新用户性别，只允许为 (male, female, unknown) | No | string |
| profile | formData | 新用户简介 | No | string |
| birth_time | formData | 新用户生日，固定格式为 2000-01-01 | No | string |
| phone_number | formData | 新用户电话号码 | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "uid": 10,     "username": "aoihosizora",     "sex": "male",     "profile": "Demo Profile",     "avatar_url": "",     "birth_time": "2019-12-18",     "authority": "admin"   } } ``` |

### /user/sub?uid

#### POST
##### Summary:

关注用户

##### Description:

关注某一用户，Auth

| code | message |
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
| Authorization | header | 用户 Token | Yes | string |
| uid | query | 对方用户 id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "me": 10,     "up": 3,     "action": "subscribe"   } } ``` |

### /user/unsub?uid

#### POST
##### Summary:

取消关注用户

##### Description:

取消关注某一用户，Auth

| code | message |
| --- | --- |
| 400 | request query param error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
| 500 | unsubscribe failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| uid | query | 对方用户 id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "me": 10,     "up": 3,     "action": "unsubscribe"   } } ``` |

### /user/{uid}

#### GET
##### Summary:

查询用户

##### Description:

查询用户信息，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户 id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "user": {       "uid": 10,       "username": "aoihosizora",       "sex": "unknown",       "profile": "",       "avatar_url": "",       "birth_time": "2000-01-01",       "authority": "admin"     },     "extra": {       "subscribing_cnt": 1,       "subscriber_cnt": 2,       "video_cnt": 0,       "playlist_cnt": 0     }   } } ``` |

### /user/{uid}/subscriber

#### GET
##### Summary:

用户粉丝

##### Description:

查询用户所有粉丝，返回分页数据，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 所查询的用户 id | Yes | integer |
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

查询用户所有关注，返回分页数据，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 所查询的用户 id | Yes | integer |
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

| code | message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 401 | need admin authority | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "count": 1,     "page": 1,     "data": [       {         "uid": 1,         "username": "User1",         "sex": "male",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "admin"       }     ]   } } ``` |

### /video/

#### POST
##### Summary:

新建视频

##### Description:

新建用户视频，Auth

| code | message |
| --- | --- |
| 400 | request form data error |
| 401 | authorization failed |
| 401 | token has expired |
| 500 | video existed failed |
| 401 | video insert failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| title | formData | 视频标题 | Yes | string |
| description | formData | 视频简介 | Yes | string |
| video_url | formData | 视频资源链接 | Yes | string |
| cover_url | formData | 视频封面链接 | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

### /video/uid/{uid}?page

#### GET
##### Summary:

查询用户视频

##### Description:

查询作者为用户的所有视频，返回分页数据，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 404 | user not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 用户 id | Yes | integer |
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

查询视频信息，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 404 | video not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| vid | path | 视频 id | Yes | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

### /video/{vid}

#### DELETE
##### Summary:

删除视频

##### Description:

删除用户视频，Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | video not found |
| 500 | video delete failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| vid | path | 删除视频 id | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

#### POST
##### Summary:

更新视频

##### Description:

更新用户视频信息，Auth

| code | message |
| --- | --- |
| 400 | request route param error |
| 400 | request format error |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | video not found |
| 404 | video update failed | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| vid | path | 更新视频 id | Yes | string |
| title | formData | 视频新标题 | No | string |
| description | formData | 视频新简介 | No | string |
| video_url | formData | 视频新资源链接 | No | string |
| cover_url | formData | 视频新封面链接 | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |

### /video?page

#### GET
##### Summary:

查询所有视频

##### Description:

管理员查询所有视频，返回分页数据，Admin

| code | message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 401 | need admin authority | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | 用户 Token | Yes | string |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {} } ``` |
