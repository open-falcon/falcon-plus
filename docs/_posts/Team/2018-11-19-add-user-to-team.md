---
category: Team
apiurl: '/api/v1/team/user'
title: "Add users to team"
type: 'POST'
sample_doc: ''
layout: default
---

添加用户到某个组
* [Session](#/authentication) Required
* users: 待添加到 team 中的用户名列表
* team_id: 目标 team 的 team_id

### Request
```{
  "team_id": 107,
  "users": ["root", "test1"]
}```

### Response

```Status: 200```
```{"message":"add successful"}```
