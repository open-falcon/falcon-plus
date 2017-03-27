---
category: DashboardScreen
apiurl: '/api/v1/dashboard/screens/pid/:screen_pid'
title: 'Gets DashboardScreens by pid'
type: 'GET'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Response

```Status: 200```
```
[
    {
        "id": 952,
        "name": "a1",
        "pid": 0
    },
    {
        "id": 961,
        "name": "laiwei-sceen1",
        "pid": 0
    }
]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
