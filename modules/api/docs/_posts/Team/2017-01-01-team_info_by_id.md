---
category: Team
apiurl: '/api/v1/team/#{team_id}'
title: "Get Team Info By Id"
type: 'GET'
sample_doc: 'team.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/team/107

### Response

```Status: 200```
```{
  "id": 107,
  "name": "ateamname",
  "resume": "i'm descript",
  "creator": 1,
  "users": [
    {
      "id": 4,
      "name": "test1",
      "cnname": "翱鶚Test",
      "email": "root123@cepave.com",
      "phone": "99999999999",
      "im": "44955834958",
      "qq": "904394234239",
      "role": 0
    },
    {
      "id": 7,
      "name": "cepave1",
      "cnname": "",
      "email": "",
      "phone": "",
      "im": "",
      "qq": "904394234239",
      "role": 0
    }
  ]
}```
