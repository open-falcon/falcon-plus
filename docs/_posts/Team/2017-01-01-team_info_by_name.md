---
category: Team
apiurl: '/api/v1/team/name/#{team_name}'
title: "Get Team Info by name"
type: 'GET'
sample_doc: 'team.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/team/name/plus-dev

### Response

```Status: 200```
```
{
    "creator": 6,
    "creator_name": "",
    "id": 10,
    "name": "plus-dev",
    "resume": "test intro",
    "users": [
        {
            "cnname": "laiwei",
            "email": "laiwei@xxx.com",
            "id": 1,
            "im": "yyyyx",
            "name": "laiwei1",
            "phone": "15011518472",
            "qq": "3805112124444455",
            "role": 2
        }
    ]
}```
