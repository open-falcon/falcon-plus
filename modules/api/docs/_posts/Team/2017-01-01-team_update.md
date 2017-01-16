---
category: Team
apiurl: '/api/v1/team'
title: "Team Update"
type: 'PUT'
sample_doc: 'team.html'
layout: default
---

更新使用者群組
* [Session](#/authentication) Required
* users: 属於该群组的user id list
* resume: team的描述
* 除Admin外, 使用者只能更新自己创建的team

### Request
```{
  "team_id": 107,
  "resume": "i'm descript update",
  "users": [4,5,6,7]
}```

### Response

```Status: 200```
```{"message":"team updated!"}```
