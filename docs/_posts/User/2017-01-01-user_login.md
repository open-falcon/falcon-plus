---
category: User
apiurl: '/api/v1/user/login'
title: 'Login'
type: 'POST'
sample_doc: 'user.html'
layout: default
---

使用者登入

### Request
```{
  "name": "test2",
  "password": "test2"
}```

### Response

```Status: 200```
```{
  "sig": "9d791331c0ea11e690c5001500c6ca5a",
  "name": "test2",
  "admin": false
}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
