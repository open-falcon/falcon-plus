---
category: Admin
apiurl: '/api/v1/admin/login'
title: 'Admin Login'
type: 'POST'
sample_doc: 'admin.html'
layout: default
---

SSO 登入

* [Session](#/authentication) Required
* `Admin` usage

### Request
```{
  "name": "test2",
}```

### Response

```Status: 200```
```{
  "sig": "9d791331c0ea11e690c5001500c6ca5a",
  "name": "test2",
  "admin": false
}```

For errors responses, see the [response status codes documentation](#/response-status-codes).
