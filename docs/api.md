# Vid Server Api

## Revision

+ UnComplete

|Date|Remark|
|--|--|
|`2019/09/01`|Complete user group & auth group|

## URI

|Method|Uri|Description|
|--|--|--|
|`POST`|`/auth/login`|Login as an exist user <sup>[1]</sup>|
|`POST`|`/auth/register`|Register and create an unexist user <sup>[1]</sup>|
|`GET`|`/user/all`|Query all users' information|
|`GET`|`/user/uid/:uid`|Query user's information <sup>[2]</sup>|
|`POST`|`/user/update`|Update user's information <sup>[1]</sup> <sup>[4]</sup>|
|`DELETE`|`/user/delete`|Delete the current user <sup>[4]</sup>|
|`GET`|`/user/subscriber/:uid`|Query user's subscribers <sup>[2]</sup>|
|`GET`|`/user/subscribing/:uid`|Query user's subscribing users <sup>[2]</sup>|
|`POST`|`/user/sub/?uid`|Subscribe the user <sup>[3]</sup> <sup>[4]</sup>|
|`POST`|`/user/unsub/?uid`|Unsubscribe the user <sup>[3]</sup> <sup>[4]</sup>|
|`GET`|`/video/all`|Query all videos|
|`GET`|`/video/uid/:uid`|Query user upload video <sup>[2]</sup>|
|`GET`|`/video/vid/:vid`|Query video <sup>[2]</sup>|

+ [1] [Need request body](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-body)
+ [2] [Need route param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-route-param)
+ [3] [Need query param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-query-param)
+ [4] [Need authorization](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-header)
+ [Response](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#response-header)

---

## Request Header

+ Routes needed authorization

|Key|Is Required|Description|
|--|--|--|
|`Authorization`|Required|User login token (Start with `Bearer`)|

## Request Query Param

+ `POST /user/sub/?uid`
+ `POST /user/unsub/?uid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`uid`|`int`|Required|Subscribe/UnSubscribe user uid||

## Request Route Param

+ `GET /user/uid/:uid`
+ `GET /user/subscriber/:uid`
+ `GET /user/subscribing/:uid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`uid`|`int`|Required|Query user uid||

+ `GET /video/uid/:uid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`uid`|`int`|Required|Query video author uid||

+ `GET /video/vid/:vid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`uid`|`int`|Required|Query video vid||

## Request Body

+ `POST /auth/login` (Raw-Json)
+ `POST /auth/register` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`username`|`string`|Required|User's username|Must in [4, 30], can't contain blankspace|
|`password`|`string`|Required|User's password|Length must in [8, 20]|

Example:

```json
{
    "username": "TestUsername",
    "password": "TestPassword"
}
```

+ `POST /user/update` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`username`|`string`|Required|New username|Must in [4, 30], can't contain blankspace|
|`profile`|`string`|Not required|User's profile|Length must smaller than 255|

```json
{
    "username": "NewUsername",
    "profile": "NewProfile"
}
```
---

## Response Header

+ `POST /auth/login`

|Field|Type|Description|Remark|
|--|--|--|--|
|`Authorization`|`string`|User login token|Default expired time is 600s|

## Response Body

+ `POST /auth/login` (Json)
+ `POST /auth/register` (Json)
+ `GET /user/all` (Array)
+ `POST /user/update` (Json)
+ `DELETE /user/delete` (Json)
+ `GET /user/subscriber/:uid` (Array)
+ `GET /user/subscribing/:uid` (Array)

|Field|Type|Description|Remark|
|--|--|--|--|
|`uid`|`int`|User uid||
|`username`|`string`|User name||
|`profile`|`string`|User profile||
|`register_time`|`date`|User register time||

Example:

```json
{
    "uid": 5,
    "username": "TestUser",
    "profile": "Test Profile",
    "register_time": "2019-09-01T14:48:08+08:00"
}
```

+ `GET /user/uid/:uid` (Json)

|Field|Type|Description|Remark|
|--|--|--|--|
|`user`|`User`|User basic information|See route `POST /auth/login` response body|
|`info`|`UserExtraInfo`|User extra information|See the next chart|

|Field|Type|Description|Remark|
|--|--|--|--|
|`subscriber_cnt`|`int`|User subscribers count||
|`subscribing_cnt`|`int`|User subscribing count||

Example:

```json
{
    "user": {
        "uid": 5,
        "username": "TestUser",
        "profile": "Test Profile",
        "register_time": "2019-09-01T14:48:08+08:00"
    },
    "info": {
        "subscriber_cnt": 2,
        "subscribing_cnt": 3
    }
}
```

+ `POST /user/sub/?uid` (Json)
+ `POST /user/unsub/?uid` (Json)

|Field|Type|Description|Remark|
|--|--|--|--|
|`me_uid`|`int`|Subscriber uid||
|`up_uid`|`int`|Subscribee uid||
|`action`|`string`|Subscribe or unsubscribe||

Example:

```json
{
    "me_uid": 8,
    "up_uid": 4,
    "action": "UnSubscribe"
}
```

+ `GET /video/all` (Array)
+ `GET /video/vid/:vid` (Json)
+ `GET /video/uid/:uid` (Array)

|Field|Type|Description|Remark|
|--|--|--|--|
|`vid`|`int`|Video vid||
|`title`|`int`|Video title||
|`description`|`string`|Video description||
|`video_url`|`string`|Video url||
|`upload_time`|`datetime`|Video upload time||
|`author`|`User`|Video author|When author is deleted, this field is `null`|

Example:

```json
{
    "vid": 1,
    "title": "Title",
    "description": "Desctiption",
    "video_url": "",
    "upload_time": "2019-09-02T17:05:00+08:00",
    "author": {
        "uid": 1,
        "username": "Username",
        "profile": "Profile",
        "register_time": "2019-09-02T16:59:53+08:00"
    }
}
```

---

## Error Type

+ See [Exception.go](https://github.com/vidorg/Vid_Backend/blob/master/exceptions/Exception.go)

## Error Code

+ See [controllers module](https://github.com/vidorg/Vid_Backend/tree/master/controllers)