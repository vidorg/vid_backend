# Vid Server Api

## Revision

+ UnComplete

|Date|Remark|
|--|--|
|`2019/09/01`|Complete user group & auth group|
|`2019/09/12`|Complete video group|

## SubRoute

|Uri|Description|
|--|--|
|`/auth/`|User register & login & modify password|
|`/user/`|User information & subscribe information|
|`/video/`|Video information & upload|

## URI

|Method|Uri|Description|
|--|--|--|
|`POST`|`/auth/login`|Login as an exist user <a href="#id1"><sup>[1]</sup></a>|
|`POST`|`/auth/register`|Register and create an unexist user <a href="#id1"><sup>[1]</sup></a>|
|`POST`|`/auth/modifypass`|Modify an exist user's password <a href="#id1"><sup>[1]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`GET`|`/user/all`|Query all users' information|
|`GET`|`/user/uid/:uid`|Query user's information <a href="#id2"><sup>[2]</sup></a>|
|`POST`|`/user/update`|Update user's information <a href="#id1"><sup>[1]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`DELETE`|`/user/delete`|Delete the current user <a href="#id4"><sup>[4]</sup></a>|
|`GET`|`/user/subscriber/:uid`|Query user's subscribers <a href="#id2"><sup>[2]</sup></a>|
|`GET`|`/user/subscribing/:uid`|Query user's subscribing users <a href="#id2"><sup>[2]</sup></a>|
|`POST`|`/user/sub?uid`|Subscribe the user <a href="#id3"><sup>[3]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`POST`|`/user/unsub?uid`|Unsubscribe the user <a href="#id3"><sup>[3]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`GET`|`/video/all`|Query all videos|
|`GET`|`/video/uid/:uid`|Query user upload video <a href="#id2"><sup>[2]</sup></a>|
|`GET`|`/video/vid/:vid`|Query video <a href="#id2"><sup>[2]</sup></a>|
|`POST`|`/video/new`|Upload new video <a href="#id1"><sup>[1]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`POST`|`/video/update`|Update video information <a href="#id1"><sup>[1]</sup></a> <a href="#id4"><sup>[4]</sup></a>|
|`DELETE`|`/video/delete?vid`|Delete current user's video <a href="#id2"><sup>[2]</sup></a> <a href="#id4"><sup>[4]</sup></a>|

+ <span id="id1">[1]</span> [Need request body](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-body)
+ <span id="id2">[2]</span> [Need route param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-route-param)
+ <span id="id3">[3]</span> [Need query param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-query-param)
+ <span id="id4">[4]</span> [Need authorization](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-header)
+ [Response](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#response-header)

---

## Request Header

+ Routes needed authorization

|Key|Is Required|Description|
|--|--|--|
|`Authorization`|Required|User login token (Start with `Bearer`)|

## Request Query Param

+ `POST /user/sub?uid`
+ `POST /user/unsub?uid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`uid`|`int`|Required|Subscribe/UnSubscribe user uid||

+ `DELETE /video/delete?vid`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`vid`|`int`|Required|Deleted user's video vid|Author must be the current user|

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
|`vid`|`int`|Required|Query video vid||

## Request Body

+ `POST /auth/login` (Raw-Json)
+ `POST /auth/register` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`username`|`string`|Required|User's username|Must in `[4, 30]`, can't contain blankspace|
|`password`|`string`|Required|User's password|Must in `[8, 20]`, can't contain blankspace|

Example:

```json
{
    "username": "TestUsername",
    "password": "TestPassword"
}
```

+ `POST /auth/modifypass` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`password`|`string`|Required|User's password|Must in `[8, 20]`, can't contain blankspace|

Example:

```json
{
    "password": "NewPassword"
}
```

+ `POST /user/update` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`username`|`string`|Not Required|New username|Must in `[4, 30]`, can't contain blankspace|
|`profile`|`string`|Not Required|User's profile|Must in `[0, 255]`|
|`sex`|`char`|Not Required|User sex|`M` or `F` or `X`|
|`avatar_url`|`date`|Not Required|User avatar url||
|`birth_time`|`date`|Not Required|User birth time||

Example:

```json
{
    "username": "NewUsername",
    "profile": "NewProfile",
    "sex": "M",
    "avatar_url": "https://github.com/fluidicon.png",
    "birth_time": "2008-01-01T00:00:00+08:00",
}
```

+ `POST /video/new` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`title`|`string`|Required|New video title||
|`description`|`string`|Not Required|New video description||
|`video_url`|`string`|Required|New video resource url|Need to be unique url|

Example:

```json
{
    "title": "New Title",
    "description": "New Description",
    "video_url": "https://xxxxxx",
}
```

+ `POST /video/update` (Raw-Json)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`vid`|`int`|Required|Update video vid||
|`title`|`string`|Not Required|Update video title||
|`description`|`string`|Not Required|Update video description||
|`video_url`|`string`|Not Required|Update video resource url|Need to be unique url|

Example:

```json
{
    "vid": 1,
    "title": "New Title",
    "description": "New Description",
    "video_url": "https://xxxxxx",
}
```

---

## Response Header

+ `POST /auth/login`

|Field|Type|Description|Remark|
|--|--|--|--|
|`Authorization`|`string`|User login token|Default expired time is `600s`|

## Response Body

+ `POST /auth/login` (Json)
+ `POST /auth/register` (Json)
+ `POST /auth/modifypass` (Json)
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
|`sex`|`char`|User sex|`M` or `F` or `X`|
|`avatar_url`|`date`|User avatar url||
|`birth_time`|`date`|User birth time||
|`register_time`|`date`|User register time||

Example:

```json
{
    "uid": 5,
    "username": "TestUser",
    "profile": "Test Profile",
    "sex": "X",
    "avatar_url": "https://github.com/fluidicon.png",
    "birth_time": "2000-01-01T00:00:00+08:00",
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
|`video_cnt`|`int`|User upload video count||

Example:

```json
{
    "user": {
        "uid": 5,
        "username": "TestUser",
        "profile": "Test Profile",
        "sex": "X",
        "avatar_url": "https://github.com/fluidicon.png",
        "birth_time": "2000-01-01T00:00:00+08:00",
        "register_time": "2019-09-01T14:48:08+08:00"
    },
    "info": {
        "subscriber_cnt": 2,
        "subscribing_cnt": 3,
        "video_cnt": 6
    }
}
```

+ `POST /user/sub?uid` (Json)
+ `POST /user/unsub?uid` (Json)

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
+ `POST /video/new` (Array)
+ `POST /video/update` (Array)
+ `DELETE /video/delete` (Array)

|Field|Type|Description|Remark|
|--|--|--|--|
|`vid`|`int`|Video vid||
|`title`|`int`|Video title||
|`description`|`string`|Video description||
|`video_url`|`string`|Video url||
|`upload_time`|`datetime`|Video upload time||
|`author`|`User`|Video author|`null` for deleted author|

Example:

```json
{
    "vid": 1,
    "title": "Title",
    "description": "Desctiption",
    "video_url": "https://xxxxxx",
    "upload_time": "2019-09-02T17:05:00+08:00",
    "author": {
        "uid": 5,
        "username": "TestUser",
        "profile": "Test Profile",
        "sex": "X",
        "avatar_url": "https://github.com/fluidicon.png",
        "birth_time": "2000-01-01T00:00:00+08:00",
        "register_time": "2019-09-01T14:48:08+08:00"
    }
}
```

---

## Error Type

+ See [Exception.go](https://github.com/vidorg/Vid_Backend/blob/master/exceptions/Exception.go)

|Error type|Description|
|--|--|
|||

Example:

```json
{
    "message": "Token has expired"
}
```

## Error Code

+ See [controllers module](https://github.com/vidorg/Vid_Backend/tree/master/controllers)