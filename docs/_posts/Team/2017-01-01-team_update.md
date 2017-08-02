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
* name: team的名字
* 除Admin外, 使用者只能更新自己创建的team

### Request
```{
  "team_id": 107,
  "name": "new_name",
  "resume": "i'm descript update",
  "users": [4]
}```

### Response

```Status: 200```
```
{
  "id": 107,
  "name": "new_name",
  "resume":"i'm descript update",
  "creator":3,
  "users":[
    {"id":4,"name":"testuser99","cnname":"testuser99","email":"","phone":"","im":"","qq":"","role":0}
  ],
  "creator_name": "",
  "message": "team updated!"
}
```
