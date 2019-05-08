---
category: User
apiurl: '/api/v1/user/u/:uid'
title: 'Update Specific User'
type: 'PUT'
sample_doc: 'user.html'
layout: default
---

更新使用者
* [Session](#/authentication) Required

### Request
```{
  "cnname": "翱鶚Test",
  "email": "root123@cepave.com",
  "im": "44955834958",
  "phone": "99999999999",
  "qq": "904394234239"
}```

### Response

```Status: 200```
```{"message":"user info updated"}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
