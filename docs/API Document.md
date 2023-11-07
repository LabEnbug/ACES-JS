# ACES#JS



# Authentication

- HTTP Authentication, scheme: bearer

# 用户

## POST 登录

POST /v1/user/login

> Body 请求参数

```yaml
username: user1
password: user1

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "exp": "string",
    "token": "string",
    "user": {
      "nickname": "string",
      "username": "string",
      "reg_time": "string",
      "user_id": 0,
      "follow_count": 0,
      "be_followed_count": 0,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "be_followed": true,
      "is_self": true,
      "avatar_url": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» exp|string|true|none||token过期时间|
|»» token|string|true|none||JWT token|
|»» user|object|true|none||用户信息|
|»»» nickname|string|true|none||用户名|
|»»» username|string|true|none||昵称|
|»»» reg_time|string|true|none||注册时间|
|»»» user_id|integer|true|none||用户id|
|»»» follow_count|integer|true|none||关注该用户的人数|
|»»» be_followed_count|integer|true|none||该用户被关注数|
|»»» be_liked_count|integer|true|none||该用户视频被点赞数|
|»»» be_favorite_count|integer|true|none||该用户视频被收藏数|
|»»» be_commented_count|integer|true|none||该用户视频被评论数|
|»»» be_forwarded_count|integer|true|none||该用户视频被转发数|
|»»» be_watched_count|integer|true|none||该用户视频被浏览量|
|»»» be_followed|boolean|true|none||被自己关注|
|»»» is_self|boolean|true|none||是否是自己|
|»»» avatar_url|string|true|none||用户的头像url|

## POST 充值

POST /v1/user/deposit

> Body 请求参数

```yaml
card_key: ACES-AAAA-AAAA-AAAA

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» card_key|body|string| 是 |充值卡号|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {},
  "err_msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|失败|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|» err_msg|string|true|none||none|

## GET 注销

GET /v1/user/logout

> 返回示例

> 成功

```json
{
  "status": 200,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

## GET 获取用户信息

GET /v1/user/info

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "exp": "string",
    "token": "string",
    "user": {
      "nickname": "string",
      "username": "string",
      "reg_time": "string",
      "user_id": 0,
      "balance": 0,
      "follow_count": 0,
      "be_followed_count": 0,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "be_followed": true,
      "is_self": true,
      "avatar_url": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» exp|string|true|none||none|
|»» token|string|true|none||none|
|»» user|object|true|none||none|
|»»» nickname|string|true|none||none|
|»»» username|string|true|none||none|
|»»» reg_time|string|true|none||none|
|»»» user_id|integer|true|none||none|
|»»» balance|number|false|none||none|
|»»» follow_count|integer|true|none||none|
|»»» be_followed_count|integer|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» be_followed|boolean|true|none||none|
|»»» is_self|boolean|true|none||none|
|»»» avatar_url|string|false|none||none|

## PUT 修改用户信息-头像

PUT /v1/user/info

> Body 请求参数

```yaml
type: avatar
file: []

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» type|body|string| 是 |none|
|» file|body|string(binary)| 是 |头像图片文件|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "user": {
      "nickname": "string",
      "reg_time": "string",
      "username": "string",
      "user_id": 0,
      "follow_count": 0,
      "be_followed_count": 0,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "be_followed": true,
      "is_self": true,
      "avatar_url": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» user|object|true|none||none|
|»»» nickname|string|true|none||none|
|»»» reg_time|string|true|none||none|
|»»» username|string|true|none||none|
|»»» user_id|integer|true|none||none|
|»»» follow_count|integer|true|none||none|
|»»» be_followed_count|integer|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» be_followed|boolean|true|none||none|
|»»» is_self|boolean|true|none||none|
|»»» avatar_url|string|true|none||none|

## GET 获取其它用户信息

GET /v1/users/{username}

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|username|path|string| 是 |用户名|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "user": {
      "nickname": "string",
      "username": "string",
      "reg_time": "string",
      "user_id": 0,
      "follow_count": 0,
      "be_followed_count": 0,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "be_followed": true,
      "is_self": true
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» user|object|true|none||none|
|»»» nickname|string|true|none||none|
|»»» username|string|true|none||none|
|»»» reg_time|string|true|none||none|
|»»» user_id|integer|true|none||none|
|»»» follow_count|integer|true|none||none|
|»»» be_followed_count|integer|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» be_followed|boolean|true|none||none|
|»»» is_self|boolean|true|none||none|

## POST 注册

POST /v1/user/signup

> Body 请求参数

```yaml
username: user20
password: pass20
nickname: 用户名啊

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |用户名|
|» password|body|string| 是 |密码|
|» nickname|body|string| 否 |昵称|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "user": {
      "user_id": 0,
      "username": "string",
      "nickname": "string",
      "follow_count": 0,
      "be_followed_count": 0,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "be_followed": true,
      "reg_time": "string",
      "is_self": true
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

*用户信息*

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» user|object|true|none||none|
|»»» user_id|integer|true|none||none|
|»»» username|string|true|none||none|
|»»» nickname|string|true|none||none|
|»»» follow_count|integer|true|none||none|
|»»» be_followed_count|integer|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» be_followed|boolean|true|none||none|
|»»» reg_time|string|true|none||none|
|»»» is_self|boolean|true|none||none|

## POST 关注用户

POST /v1/users/{username}/follow

失败原因一般是已经关注，此时不会变更关注状态，仍然为已关注。

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|username|path|string| 是 |用户名|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

## DELETE 取消关注用户

DELETE /v1/users/{username}/follow

失败原因一般是未关注或者没有关注过，此时不会变更关注状态，仍然为不关注。

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|username|path|string| 是 |用户名|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

# 视频列表

## GET 获取相关视频列表

GET /v1/videos/{videoUid}/related

获取的是单个视频的相关视频列表，使用了视频侧的推荐算法。

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |none|
|limit|query|number| 否 |none|
|start|query|number| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video_list": [
      {
        "user": {
          "nickname": "string",
          "reg_time": "string",
          "user_id": 0,
          "username": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "is_self": true
        },
        "video_uid": "string",
        "type": 0,
        "content": "string",
        "keyword": "string",
        "upload_time": "string",
        "cover_url": "string",
        "play_url": "string",
        "is_user_liked": true,
        "is_user_favorite": true,
        "is_user_uploaded": true,
        "is_user_last_play": true,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "is_user_watched": true,
        "is_top": true,
        "is_private": true
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video_list|[object]|true|none||none|
|»»» user|object|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»» video_uid|string|true|none||none|
|»»» type|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» keyword|string|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» cover_url|string|true|none||none|
|»»» play_url|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## GET 获取单个视频信息

GET /v1/videos/{videoUid}

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video": {
      "video_uid": "string",
      "type": 0,
      "content": "string",
      "upload_time": "string",
      "cover_url": "string",
      "play_url": "string",
      "user": {
        "nickname": "string",
        "reg_time": "string",
        "user_id": 0,
        "username": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "is_self": true
      },
      "keyword": "string",
      "is_user_liked": true,
      "is_user_favorite": true,
      "is_user_uploaded": true,
      "is_user_last_play": true,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "is_user_watched": true,
      "is_top": true,
      "is_private": true
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video|object|true|none||none|
|»»» video_uid|string|true|none||视频uuid|
|»»» type|integer|true|none||视频类型 1为普通|
|»»» content|string|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» cover_url|string|true|none|封面url|如为空值，则说明未获取封面|
|»»» play_url|string|true|none|视频url|如为空值，则说明未转码好视频|
|»»» user|object|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»» keyword|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## GET 获取视频列表(指定类型-娱乐)

GET /v1/videos

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|limit|query|string| 是 |限制条数|
|start|query|string| 是 |开始于|
|type|query|string| 是 |视频类型，如4是娱乐|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video_list": [
      {
        "video_uid": "string",
        "type": 0,
        "content": "string",
        "upload_time": "string",
        "cover_url": "string",
        "play_url": "string",
        "user": {
          "nickname": "string",
          "reg_time": "string",
          "user_id": 0,
          "username": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "is_self": true
        },
        "keyword": "string",
        "is_user_liked": true,
        "is_user_favorite": true,
        "is_user_uploaded": true,
        "is_user_last_play": true,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "is_user_watched": true,
        "is_top": true,
        "is_private": true
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video_list|[object]|true|none||none|
|»»» video_uid|string|true|none||none|
|»»» type|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» cover_url|string|true|none||none|
|»»» play_url|string|true|none||none|
|»»» user|object|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»» keyword|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## GET 获取视频类型信息

GET /v1/video/types

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video_types": [
      {
        "id": 0,
        "type_name": "string"
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video_types|[object]|true|none||none|
|»»» id|integer|true|none||none|
|»»» type_name|string|true|none||none|

# 视频操作/不用登录也能进行的视频操作

## POST 记录播放

POST /v1/videos/{videoUid}/actions/{action}

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|action|path|string| 是 |watch|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "action": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» action|string|true|none||none|

# 视频操作/所有者才能进行的视频操作

## PUT 确认发布新视频

PUT /v1/video/upload/{videoUid}

> Body 请求参数

```yaml
video_content: 视频内容
video_keyword: 视频关键词
video_type: "4"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |上传的视频uid|
|body|body|object| 否 |none|
|» video_content|body|string| 是 |内容|
|» video_keyword|body|string| 是 |关键词|
|» video_type|body|integer| 是 |视频类型|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video": {
      "user": {
        "nickname": "string",
        "reg_time": "string",
        "user_id": 0,
        "username": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "is_self": true,
        "avatar_url": "string"
      },
      "video_uid": "string",
      "type": 0,
      "content": "string",
      "keyword": "string",
      "upload_time": "string",
      "cover_url": "string",
      "play_url": "string",
      "is_user_liked": true,
      "is_user_favorite": true,
      "is_user_uploaded": true,
      "is_user_last_play": true,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "is_user_watched": true,
      "is_top": true,
      "is_private": true
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video|object|true|none||none|
|»»» user|object|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» video_uid|string|true|none||none|
|»»» type|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» keyword|string|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» cover_url|string|true|none||none|
|»»» play_url|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## PUT 修改视频信息

PUT /v1/videos/{videoUid}

> Body 请求参数

```yaml
video_content: 视频内容
video_keyword: 视频关键词
video_type: "4"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|body|body|object| 否 |none|
|» video_content|body|string| 是 |内容|
|» video_keyword|body|string| 是 |关键词|
|» video_type|body|integer| 是 |视频类型|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video": {
      "user": {
        "nickname": "string",
        "reg_time": "string",
        "user_id": 0,
        "username": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "is_self": true,
        "avatar_url": "string"
      },
      "video_uid": "string",
      "type": 0,
      "content": "string",
      "keyword": "string",
      "upload_time": "string",
      "cover_url": "string",
      "play_url": "string",
      "is_user_liked": true,
      "is_user_favorite": true,
      "is_user_uploaded": true,
      "is_user_last_play": true,
      "be_liked_count": 0,
      "be_favorite_count": 0,
      "be_commented_count": 0,
      "be_forwarded_count": 0,
      "be_watched_count": 0,
      "is_user_watched": true,
      "is_top": true,
      "is_private": true
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video|object|true|none||none|
|»»» user|object|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» video_uid|string|true|none||none|
|»»» type|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» keyword|string|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» cover_url|string|true|none||none|
|»»» play_url|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## DELETE 删除视频

DELETE /v1/videos/{videoUid}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

## POST 上传新视频

POST /v1/video/upload

先上传文件，后进行信息更新来确认发布

> Body 请求参数

```yaml
file: []

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» file|body|string(binary)| 是 |视频文件|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video_uid": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video_uid|string|true|none||none|

# 视频操作/非所有者也能进行的视频操作

## DELETE 视频操作-取消-喜欢/收藏

DELETE /v1/videos/{videoUid}/actions/{action}

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|action|path|string| 是 |like/favorite|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "action": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» action|string|true|none||none|

# 视频评论

## POST 评论视频的评论(子评论)

POST /v1/videos/{videoUid}/comments

> Body 请求参数

```yaml
content: 小金毛真可爱3
quote_comment_id: "133"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|body|body|object| 否 |none|
|» content|body|string| 是 |评论内容|
|» quote_comment_id|body|integer| 是 |引用的评论id|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "comment": {
      "id": 0,
      "content": "string",
      "quote_comment_id": 0,
      "quote_child_comment_id": 0,
      "user": {
        "user_id": 0,
        "username": "string",
        "nickname": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "reg_time": "string",
        "is_self": true,
        "avatar_url": "string"
      },
      "quote_user": null,
      "comment_time": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» comment|object|true|none||none|
|»»» id|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» quote_comment_id|integer|true|none||none|
|»»» quote_child_comment_id|integer|true|none||none|
|»»» user|object|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» quote_user|null|true|none||none|
|»»» comment_time|string|true|none||none|

## GET 获取评论(子评论)

GET /v1/videos/{videoUid}/comments

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|limit|query|string| 是 |none|
|start|query|string| 是 |none|
|comment_id|query|integer| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "child_comment_count_left": 0,
    "child_comment_list": [
      {
        "id": 0,
        "content": "string",
        "quote_comment_id": 0,
        "quote_child_comment_id": 0,
        "user": {
          "user_id": 0,
          "username": "string",
          "nickname": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "reg_time": "string",
          "is_self": true,
          "avatar_url": "string"
        },
        "quote_user": {
          "user_id": 0,
          "username": "string",
          "nickname": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "reg_time": "string",
          "is_self": true,
          "avatar_url": "string"
        },
        "comment_time": "string"
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» child_comment_count_left|integer|true|none||none|
|»» child_comment_list|[object]|true|none||none|
|»»» id|integer|false|none||none|
|»»» content|string|false|none||none|
|»»» quote_comment_id|integer|false|none||none|
|»»» quote_child_comment_id|integer|false|none||none|
|»»» user|object|false|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» quote_user|object|false|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» comment_time|string|false|none||none|

# 视频弹幕

## POST 发弹幕

POST /v1/videos/{videoUid}/bullet_comments

> Body 请求参数

```yaml
content: 小金毛真可爱15
comment_at: "1.5"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|body|body|object| 否 |none|
|» content|body|string| 是 |弹幕内容|
|» comment_at|body|number| 是 |弹幕所在视频播放的时间|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "bullet_comment": {
      "id": 0,
      "content": "string",
      "comment_at": 0,
      "user": {
        "user_id": 0,
        "username": "string",
        "nickname": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "reg_time": "string",
        "is_self": true,
        "avatar_url": "string"
      },
      "comment_time": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» bullet_comment|object|true|none||none|
|»»» id|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» comment_at|number|true|none||none|
|»»» user|object|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» comment_time|string|true|none||none|

## GET 获取弹幕列表

GET /v1/videos/{videoUid}/bullet_comments

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|videoUid|path|string| 是 |视频uid|
|limit|query|integer| 是 |none|
|start|query|integer| 是 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "bullet_comment_count": 0,
    "bullet_comment_list": [
      {
        "id": 0,
        "content": "string",
        "comment_at": 0,
        "user": {
          "user_id": 0,
          "username": "string",
          "nickname": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "reg_time": "string",
          "is_self": true,
          "avatar_url": "string"
        },
        "comment_time": "string"
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» bullet_comment_count|integer|true|none||none|
|»» bullet_comment_list|[object]|true|none||none|
|»»» id|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» comment_at|number|true|none||none|
|»»» user|object|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»»» avatar_url|string|true|none||none|
|»»» comment_time|string|true|none||none|

# 搜索

## GET 获取视频搜索hotkeys

GET /v1/search/video/hotkeys

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|max_count|query|number| 否 |5~20|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "hotkeys": [
      "string"
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» hotkeys|[string]|true|none||none|

## GET 获取用户搜索hotkeys

GET /v1/search/user/hotkeys

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|max_count|query|number| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "hotkeys": [
      "string"
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» hotkeys|[string]|true|none||none|

## GET 视频搜索

GET /v1/search/video

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|keyword|query|string| 否 |搜索的关键词|
|limit|query|string| 否 |1~24|
|start|query|string| 否 |开始于第几条|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "video_list": [
      {
        "video_uid": "string",
        "type": 0,
        "content": "string",
        "keyword": "string",
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "user": {
          "user_id": 0,
          "username": "string",
          "nickname": "string",
          "follow_count": 0,
          "be_followed_count": 0,
          "be_liked_count": 0,
          "be_favorite_count": 0,
          "be_commented_count": 0,
          "be_forwarded_count": 0,
          "be_watched_count": 0,
          "be_followed": true,
          "reg_time": "string",
          "is_self": true
        },
        "upload_time": "string",
        "is_user_liked": true,
        "is_user_favorite": true,
        "is_user_uploaded": true,
        "is_user_watched": true,
        "is_user_last_play": true,
        "cover_url": "string",
        "play_url": "string",
        "is_top": true,
        "is_private": true
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» video_list|[object]|true|none||none|
|»»» video_uid|string|true|none||none|
|»»» type|integer|true|none||none|
|»»» content|string|true|none||none|
|»»» keyword|string|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» user|object|true|none||none|
|»»»» user_id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» follow_count|integer|true|none||none|
|»»»» be_followed_count|integer|true|none||none|
|»»»» be_liked_count|integer|true|none||none|
|»»»» be_favorite_count|integer|true|none||none|
|»»»» be_commented_count|integer|true|none||none|
|»»»» be_forwarded_count|integer|true|none||none|
|»»»» be_watched_count|integer|true|none||none|
|»»»» be_followed|boolean|true|none||none|
|»»»» reg_time|string|true|none||none|
|»»»» is_self|boolean|true|none||none|
|»»» upload_time|string|true|none||none|
|»»» is_user_liked|boolean|true|none||none|
|»»» is_user_favorite|boolean|true|none||none|
|»»» is_user_uploaded|boolean|true|none||none|
|»»» is_user_watched|boolean|true|none||none|
|»»» is_user_last_play|boolean|true|none||none|
|»»» cover_url|string|true|none||none|
|»»» play_url|string|true|none||none|
|»»» is_top|boolean|true|none||none|
|»»» is_private|boolean|true|none||none|

## GET 用户搜索

GET /v1/search/user

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|keyword|query|string| 否 |搜索的关键词|
|limit|query|string| 否 |1~24|
|start|query|string| 否 |开始于第几条|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {
    "user_list": [
      {
        "user_id": 0,
        "username": "string",
        "nickname": "string",
        "follow_count": 0,
        "be_followed_count": 0,
        "be_liked_count": 0,
        "be_favorite_count": 0,
        "be_commented_count": 0,
        "be_forwarded_count": 0,
        "be_watched_count": 0,
        "be_followed": true,
        "reg_time": "string",
        "is_self": true
      }
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|
|»» user_list|[object]|true|none||none|
|»»» user_id|integer|true|none||none|
|»»» username|string|true|none||none|
|»»» nickname|string|true|none||none|
|»»» follow_count|integer|true|none||none|
|»»» be_followed_count|integer|true|none||none|
|»»» be_liked_count|integer|true|none||none|
|»»» be_favorite_count|integer|true|none||none|
|»»» be_commented_count|integer|true|none||none|
|»»» be_forwarded_count|integer|true|none||none|
|»»» be_watched_count|integer|true|none||none|
|»»» be_followed|boolean|true|none||none|
|»»» reg_time|string|true|none||none|
|»»» is_self|boolean|true|none||none|

# 回调

## POST 七牛视频转码回调

POST /callback/qiniu/hls

> Body 请求参数

```json
{
  "a": "string",
  "version": "string",
  "id": "string",
  "reqid": "string",
  "pipeline": "string",
  "input": {
    "kodo_file": {
      "bucket": "string",
      "key": "string"
    }
  },
  "code": 0,
  "desc": "string",
  "ops": [
    {
      "id": "string",
      "fop": {
        "cmd": "string",
        "input_from": "string",
        "result": {
          "code": 0,
          "desc": "string",
          "has_output": true,
          "kodo_file": {}
        }
      },
      "depends": [
        "string"
      ]
    }
  ],
  "created_at": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» a|body|string| 是 |none|
|» version|body|string| 是 |none|
|» id|body|string| 是 |none|
|» reqid|body|string| 是 |none|
|» pipeline|body|string| 是 |none|
|» input|body|object| 是 |none|
|»» kodo_file|body|object| 是 |none|
|»»» bucket|body|string| 是 |none|
|»»» key|body|string| 是 |none|
|» code|body|integer| 是 |none|
|» desc|body|string| 是 |none|
|» ops|body|[object]| 是 |none|
|»» id|body|string| 是 |none|
|»» fop|body|object| 是 |none|
|»»» cmd|body|string| 是 |none|
|»»» input_from|body|string| 否 |none|
|»»» result|body|object| 是 |none|
|»»»» code|body|integer| 是 |none|
|»»»» desc|body|string| 是 |none|
|»»»» has_output|body|boolean| 是 |none|
|»»»» kodo_file|body|object| 是 |none|
|»»»»» bucket|body|string| 是 |none|
|»»»»» key|body|string| 是 |none|
|»»»»» hash|body|string| 是 |none|
|»» depends|body|[string]| 否 |none|
|» created_at|body|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

## POST 七牛视频截图回调

POST /callback/qiniu/screenshot

> Body 请求参数

```json
{
  "a": "string",
  "version": "string",
  "id": "string",
  "reqid": "string",
  "pipeline": "string",
  "input": {
    "kodo_file": {
      "bucket": "string",
      "key": "string"
    }
  },
  "code": 0,
  "desc": "string",
  "ops": [
    {
      "id": "string",
      "fop": {
        "cmd": "string",
        "input_from": "string",
        "result": {
          "code": 0,
          "desc": "string",
          "has_output": true,
          "kodo_file": {}
        }
      },
      "depends": [
        "string"
      ]
    }
  ],
  "created_at": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» a|body|string| 是 |none|
|» version|body|string| 是 |none|
|» id|body|string| 是 |none|
|» reqid|body|string| 是 |none|
|» pipeline|body|string| 是 |none|
|» input|body|object| 是 |none|
|»» kodo_file|body|object| 是 |none|
|»»» bucket|body|string| 是 |none|
|»»» key|body|string| 是 |none|
|» code|body|integer| 是 |none|
|» desc|body|string| 是 |none|
|» ops|body|[object]| 是 |none|
|»» id|body|string| 是 |none|
|»» fop|body|object| 是 |none|
|»»» cmd|body|string| 是 |none|
|»»» input_from|body|string| 否 |none|
|»»» result|body|object| 是 |none|
|»»»» code|body|integer| 是 |none|
|»»»» desc|body|string| 是 |none|
|»»»» has_output|body|boolean| 是 |none|
|»»»» kodo_file|body|object| 是 |none|
|»»»»» bucket|body|string| 是 |none|
|»»»»» key|body|string| 是 |none|
|»»»»» hash|body|string| 是 |none|
|»» depends|body|[string]| 否 |none|
|» created_at|body|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "status": 0,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» status|integer|true|none||none|
|» data|object|true|none||none|

