---
category: Team
apiurl: '/api/v1/team'
title: "Team Create"
type: 'POST'
sample_doc: 'team.html'
layout: default
---

新增使用者群組
* [Session](#/authentication) Required
* users: 属於该群组的user id list
* resume: team的描述

### Request
```{"team_name": "ateamname","resume": "i'm descript", "users": [1]}```

### Response

```Status: 200```
```
{
  "team": {
    "id": 6,
    "name": "ateamname",
    "resume": "i'm descript",
    "creator": 3
  },
  "message": "team created! Afftect row: 1, Affect refer: 1"
}
```
