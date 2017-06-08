---
category: User
apiurl: '/api/v1/user/u/:uid/teams'
title: 'Get user teams'
type: 'GET'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/user/u/4/teams

### Response

```Status: 200```
```{"teams":
  [{
    "id":3,
    "name":"root",
    "resume":"",
    "creator":5},
   {"id":32,
    "name":"testteam",
    "resume":"test22",
    "creator":5
   }]
} ```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
