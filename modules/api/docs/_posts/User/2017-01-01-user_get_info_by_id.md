---
category: User
apiurl: '/api/v1/user/u/#{user_id}'
title: 'Get User info by id'
type: 'GET'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required
* `Admin` usage
* ex. /api/v1/user/u/4

### Response

```Status: 200```
```{
  "id": 4,
  "name": "userA",
  "cnname": "tear",
  "email": "",
  "phone": "",
  "im": "",
  "qq": "",
  "role": 0
}```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
