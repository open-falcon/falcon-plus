---
category: User
apiurl: '/api/v1/user/name/#{user_name}'
title: 'Get User info by name'
type: 'GET'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required
* `Admin` usage
* ex. /api/v1/user/name/laiwei

### Response

```Status: 200```
```
{
    "cnname": "laiwei8",
    "email": "laiwei8@xx",
    "id": 8,
    "im": "",
    "name": "laiwei8",
    "phone": "",
    "qq": "",
    "role": 0
}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
