---
category: User
apiurl: '/api/v1/user/u/:uid/in_teams'
title: 'Check user in teams or not'
type: 'GET'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/user/u/4/in_teams?team_names=team1,team4

### Request
Content-type: application/x-www-form-urlencoded
```team_names=team1,team2```

### Response

```Status: 200```
```{"message":"true"} ```

For more example, see the [user](/doc/user.html).

For errors responses, see the [response status codes documentation](#/response-status-codes).
