---
category: User
apiurl: '/api/v1/user/current'
title: 'Current User info'
type: 'GET'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required
* 当前使用者资讯

### Response

```Status: 200```
```{
  "id": 2,
  "name": "root",
  "cnname": "",
  "email": "",
  "phone": "",
  "im": "",
  "qq": "",
  "role": 2
}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
