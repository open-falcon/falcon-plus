---
category: User
apiurl: '/api/v1/user/create'
title: 'Create User'
type: 'POST'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```{"name": "test","password": "test", "email":"xxx@xxx.com", "cnname": "翱鹗"}```

### Response

```Status: 200```
```{
  "name": "owltester",
  "password": "mypassword",
  "cnname": "翱鹗",
  "email": "root123@cepave.com",
  "im": "44955834958",
  "phone": "99999999999",
  "qq": "904394234239"
}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
