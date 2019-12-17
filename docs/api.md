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
| 400 | request form data exception |
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
| 400 | request form data exception |
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
| 400 | request form data exception |
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
| 400 | request form data exception |
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

#### PUT
##### Summary:

更新用户

##### Description:

更新用户信息，Auth

| code | message |
| --- | --- |
| 401 | authorization failed |
| 401 | token has expired |
| 404 | user not found |
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

### /user/sub

#### POST
##### Summary:

关注用户

##### Description:

关注某一用户，Auth

| code | message |
| --- | --- |
| 400 | request query param exception |
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

### /user/unsub

#### POST
##### Summary:

取消关注用户

##### Description:

取消关注某一用户，Auth

| code | message |
| --- | --- |
| 400 | request query param exception |
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

普通用户查询用户信息，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param exception |
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
| 400 | request route param exception |
| 404 | user not found | 

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| uid | path | 所查询的用户 id | Yes | integer |
| page | query | 分页 | No | integer |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ```json {   "code": 200,   "message": "Success",   "data": {     "count": 2,     "page": 1,     "data": [       {         "uid": 1,         "username": "User1",         "sex": "male",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "admin"       },       {         "uid": 2,         "username": "User2",         "sex": "unknown",         "profile": "",         "avatar_url": "",         "birth_time": "2000-01-01",         "authority": "normal"       }     ]   } } ``` |

### /user/{uid}/subscribing

#### GET
##### Summary:

用户关注的人

##### Description:

查询用户所有关注，返回分页数据，Non-Auth

| code | message |
| --- | --- |
| 400 | request route param exception |
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
