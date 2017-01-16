---
category: Team
apiurl: '/api/v1/team/#{team_id}'
title: "Delete Team By Id"
type: 'DELETE'
sample_doc: 'team.html'
layout: default
---

新增使用者群組
* [Session](#/authentication) Required
* ex. /api/v1/team/107
* 除Admin外, 使用者只能更新自己创建的team

### Response

```Status: 200```
```{"message":"team 107 is deleted. Affect row: 1 / refer delete: 4"}```
