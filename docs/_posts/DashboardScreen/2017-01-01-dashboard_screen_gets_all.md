---
category: DashboardScreen
apiurl: '/api/v1/dashboard/screens'
title: 'Gets all DashboardScreens'
type: 'GET'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Request
Content-type: application/x-www-form-urlencoded

* params:
  * limit: 【选填】查询最大数据量，默认值：500，如：limit=10 只查询最多10条数据

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
        "id": 953,
        "name": "aa1",
        "pid": 952
    },
    {
        "id": 968,
        "name": "laiwei-screen2",
        "pid": 1
    },
    {
        "id": 972,
        "name": "laiwei-sceen1",
        "pid": 0
    },
    {
        "id": 991,
        "name": "xnew",
        "pid": 972
    },
    {
        "id": 993,
        "name": "clone3",
        "pid": 972
    },
    {
        "id": 995,
        "name": "op",
        "pid": 0
    }
]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
