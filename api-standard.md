# HTTP API建议规范

## 说明
本规范参考阮一峰的《[RESTful API设计指南](http://www.ruanyifeng.com/blog/2014/05/restful_api.html)》编写，部分内容做了适应性的调整和修改。

## 协议
本文档目前只针对HTTP与HTTPS协议，提供参考规范和建议。

## 域名
应尽量将API模块部署在专用域名之下。
```
http://api.open-falcon.org
```

## 版本(Versioning)
应将API版本号放入URL。
TODO：此处v1高亮
```
http://api.open-falcon.org/v1/
```

## 路径(Endpoint)
路径又称``“终点”(endpoint)``，表示API的具体网址。

 路径中以``“/”``区分的每一项都应该是``小写英文单词``，若包含多个单词，需要用``下划线``分隔。
在RESTful架构中，每个网址代表一种``资源（resource）``，所以网址中不应包含动词，如show、delete等。只能有名词，而且所用的名词往往与数据库的表格名对应。

一般来说，数据库中的表都是同种记录的``“集合”（collection）``，所以API中的名词也应该使用``复数``。

以Open-Falcon举例，API模块提供各种关于监控系统的信息，包括数据及报警：
```
· http://api.open-falcon.org/v1/counters
· http://api.open-falcon.org/v1/strategies
· http://api.open-falcon.org/v1/points
```

## HTTP动词
对于资源的具体操作类型，由HTTP动词表示。

常用的HTTP动词有下面五个（括号里是对应的SQL命令）。

```
· GET（SELECT）：从服务器取出资源（一项或多项）
· POST（CREATE）：在服务器新建一个资源
· PUT（UPDATE）：在服务器更新资源（客户端提供改变后的完整资源）
· PATCH（UPDATE）：在服务器更新资源（客户端提供改变的属性）
· DELETE（DELETE）：从服务器删除资源
```
PATCH一般不常用，更新资源建议使用```PUT```。

还有两个不常用的动词。
```
HEAD：获取资源的元数据。
OPTIONS：获取信息，关于资源的哪些属性是客户端可以改变的。
```

以下是一些例子
```
· GET /v1/users : 获取所有的用户
· GET /v1/users/{int:ID} : 获取指定的用户信息
· POST /v1/users : 新增一个用户
· PUT /v1/users/{int:ID} : 更新某个指定用户的信息(提供用户的全部信息)
· PATCH /v1/users/{int:ID}  : 更新某个指定用户的信息（提供用户的部分信息）
```

## 请求参数、请求体
如果记录数很多，服务器不可能都将他们返回给用户。API应该提供参数，支持过滤结果返回。
下面是一些常见的参数。
```
· ?limit=10 : 指定返回记录的数量
· ?offset=10 : 指定返回记录的开始位置
· ?page=2&per_page=100 : 指定第几页，以及每页的记录数
· ?sortby=name&order=asc : 指定返回结果按照哪个属性排序，以及排序顺序
· ?type_id=1 : 指定筛选条件
```
其中，参数的命名规则为``snake_case``，即单词全部小写，单词间使用下划线分隔。 

参数的设计允许存在``冗余``，即允许API路径和URL参数偶尔有重复。比如，GET /services/ID/domains 与 GET /services?domain_id=ID 的含义是相同的。

另外，参考HTTP协议(1.0版和1.1版)的主要设计者Roy Thomas Fielding所言，任何HTTP请求都允许包含请求体，但发送一个请求体非空的GET请求是没有意义的。因此GET请求请求体应该为``空``。

对于POST、PUT、DELETE请求，若请求体不为空，则请求体格式应该为JSON，如: 
```json
{
	"id" : 1,
	"name" : ZhangSan
}
```
需要注意的是，JSON中key的命名规则同样需要遵循``snake_case``。 

## 状态码（Status Code）
服务器向用户返回的状态码和提示信息，常见的有以下一些（方括号中是该状态码对应的HTTP动词）。
```
· 200 OK - [GET]：服务器成功返回用户请求的数据，该操作是幂等的（Idempotent）。
· 201 CREATED - [POST/PUT/PATCH]：用户新建或修改数据成功。
· 202 Accepted - [*]：表示一个请求已经进入后台排队（异步任务）
· 204 NO CONTENT - [DELETE]：用户删除数据成功。
· 400 INVALID REQUEST - [POST/PUT/PATCH]：用户发出的请求有错误，服务器没有进行新建或修改数据的操作，该操作是幂等的。
· 401 Unauthorized - [*]：表示用户没有权限（令牌、用户名、密码错误）。
· 403 Forbidden - [*] 表示用户得到授权（与401错误相对），但是访问是被禁止的。
· 404 NOT FOUND - [*]：用户发出的请求针对的是不存在的记录，服务器没有进行操作，该操作是幂等的。
· 406 Not Acceptable - [GET]：用户请求的格式不可得（比如用户请求JSON格式，但是只有XML格式）。
· 410 Gone -[GET]：用户请求的资源被永久删除，且不会再得到的。
· 422 Unprocesable entity - [POST/PUT/PATCH] 当创建一个对象时，发生一个验证错误。
· 500 INTERNAL SERVER ERROR - [*]：服务器发生错误，用户将无法判断发出的请求是否成功。
```
状态码的完整列表，参见[这里](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html)。

## 返回结果
服务器返回的数据格式，统一为``JSON``。针对资源的不同操作，服务器向用户返回的结果

应该符合以下规范：
```
· GET /collection : 返回资源对象的列表(数组)
· DELETE /collection/resource : 返回如下的json响应，其中data为被删除对象.
· GET /collection/resource : 返回单个资源对象
· POST /collection : 返回新生成的资源对象
· PUT /collection/resource : 返回更改后的完整的资源对象 
```

响应格式统一如下：
```
{
    "code": 200,    // http status code
    "data": data,    // return data if success
    "ret_msg": "return message if failed",    // return message if failed, when success, it's empty
    "debug_msg": "used for debug",    // used for debug when failed,when success, it's emptuy
    "url_manual": "manual url to show if failed",    // manual url to show if failed
    "url_redirect": "used for 3xx results",    // used for 3xx results
    “response_type": "application/json"    // response body type
}

```
