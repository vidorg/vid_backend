# Vid Server Api

## Revision

+ UnComplete

|Date|Remark|
|--|--|
|`2019/09/01`|Complete user group & auth group|
|`2019/09/12`|Complete video group|
|`2019/09/19`|Complete search group & raw group|

## SubRoute

|Uri|Description|
|--|--|
|`/auth/`|User register & login & modify password|
|`/user/`|User information & subscribe information|
|`/video/`|Video information|
|`/search/`|User & video search|
|`/raw/`|Image & video upload & download|

## URI

|Method|Uri|Description|
|--|--|--|
|`POST`|`/auth/login`|Login as an exist user <sup>[1]</sup>|
|`POST`|`/auth/register`|Register and create an unexist user <sup>[1]</sup>|
|`POST`|`/auth/modifypass`|Modify an exist user's password <sup>[1]</sup> <sup>[4]</sup>|
|`GET`|`/user/all`|Query all users' information <sup>[5]</sup>|
|`GET`|`/user/uid/:uid`|Query user's information <sup>[2]</sup>|
|`POST`|`/user/update`|Update user's information <sup>[1]</sup> <sup>[4]</sup>|
|`DELETE`|`/user/delete`|Delete the current user <sup>[4]</sup>|
|`GET`|`/user/subscriber/:uid`|Query user's subscribers <sup>[2]</sup>|
|`GET`|`/user/subscribing/:uid`|Query user's subscribing users <sup>[2]</sup>|
|`POST`|`/user/sub?uid`|Subscribe the user <sup>[3]</sup> <sup>[4]</sup>|
|`POST`|`/user/unsub?uid`|Unsubscribe the user <sup>[3]</sup> <sup>[4]</sup>|
|`GET`|`/video/all`|Query all videos <sup>[5]</sup>|
|`GET`|`/video/uid/:uid`|Query user upload video <sup>[2]</sup>|
|`GET`|`/video/vid/:vid`|Query video <sup>[2]</sup>|
|`POST`|`/video/new`|Upload new video <sup>[1]</sup> <sup>[4]</sup>|
|`POST`|`/video/update`|Update video information <sup>[1]</sup> <sup>[4]</sup>|
|`DELETE`|`/video/delete?vid`|Delete current user's video <sup>[3]</sup> <sup>[4]</sup>|
|`GET`|`/search/user?keyword`|Search user by keyword or uid <sup>[3]</sup>|
|`GET`|`/search/video?keyword`|Search video by keyword or vid <sup>[3]</sup>|
|`POST`|`/raw/upload/image`|Upload user's image <sup>[1]</sup> <sup>[4]</sup>|
|`POST`|`/raw/upload/video`|Upload user's video <sup>[1]</sup> <sup>[4]</sup>|
|`GET`|`/raw/image/:user/:filename`|Download image <sup>[2]</sup>|
|`GET`|`/raw/video/:user/:filename`|Download video <sup>[2]</sup>|

+ [1] [Need request body](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-body)
+ [2] [Need route param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-route-param)
+ [3] [Need query param](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-query-param)
+ [4] [Need authorization](https://github.com/vidorg/Vid_Backend/blob/master/docs/api.md#request-header)
+ [5] Need admin authority
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

+ `GET /search/user?keyword`
+ `GET /search/video?keyword`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`keyword`|`string`|Required|Search keyword|Can start with `uid:` or `vid:`|

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

+ `GET /raw/image/:user/:filename`
+ `GET /raw/video/:user/:filename`

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`user`|`int`|Required|Resource's author uid||
|`filename`|`string`|Required|Resource's filename||

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

+ `POST /raw/upload/image` (Form-Data)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`image`|`File`|Required|Upload image file|Only support `jpg` `png` `bmp`|

+ `POST /raw/upload/video` (Form-Data)

|Field|Type|Is Required|Description|Remark|
|--|--|--|--|--|
|`video`|`File`|Required|Upload video file|Only support `mp4`|

---

## Response Header

+ `POST /auth/login`

|Field|Type|Description|Remark|
|--|--|--|--|
|`Authorization`|`string`|User login token|Default expired time is `600s`|

+ `GET /raw/image/:user/:filename`
+ `GET /raw/video/:user/:filename`

|Field|Type|Description|Remark|
|--|--|--|--|
|`Content-Type`|`string`|Resource content type|`image/png` or `video/mpeg4`|

## Response Body

+ `POST /auth/login` (Json)
+ `POST /auth/register` (Json)
+ `POST /auth/modifypass` (Json)
+ `GET /user/all` (Array)
+ `POST /user/update` (Json)
+ `DELETE /user/delete` (Json)
+ `GET /user/subscriber/:uid` (Array)
+ `GET /user/subscribing/:uid` (Array)
+ `GET /search/user?keyword` (Array)

|Field|Type|Description|Remark|
|--|--|--|--|
|`uid`|`int`|User uid||
|`username`|`string`|User name||
|`profile`|`string`|User profile||
|`sex`|`char`|User sex|`M` or `F` or `X`|
|`avatar_url`|`date`|User avatar url||
|`birth_time`|`date`|User birth time||
|`register_time`|`date`|User register time||
|`authority`|`string`|User authority|`ENUM('admin', 'normal')`|

Example:

```json
{
    "uid": 5,
    "username": "TestUser",
    "profile": "Test Profile",
    "sex": "X",
    "avatar_url": "https://github.com/fluidicon.png",
    "birth_time": "2000-01-01T00:00:00+08:00",
    "register_time": "2019-09-01T14:48:08+08:00",
    "authority": "admin"
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
        "register_time": "2019-09-01T14:48:08+08:00",
        "authority": "admin"
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
+ `GET /search/video?keyword` (Array)

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
        "register_time": "2019-09-01T14:48:08+08:00",
        "authority": "admin"
    }
}
```

+ `POST /raw/upload/image`
+ `POST /raw/upload/video`

|Field|Type|Description|Remark|
|--|--|--|--|
|`type`|`string`|Upload type|`Image` or `Video`|
|`url`|`string`|Upload resource url|New filename is modified as the current time|

Example:

```json
{
    "type": "Image",
    "url": "http://127.0.0.1:1234/raw/image/2/20190919172804.png"
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