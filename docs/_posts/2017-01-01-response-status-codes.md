---
title: 'Response status codes'

layout: default
---

### Success

* `POST`, `GET`, `PUT`, `DELETE` returns `200` on success,
* 当参数使用不正确的时候会回覆 `400`

### 参考

return code 请参考 [status.go](https://golang.org/src/net/http/status.go)
