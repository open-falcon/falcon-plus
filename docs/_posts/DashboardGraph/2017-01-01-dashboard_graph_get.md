---
category: DashboardGraph
apiurl: '/api/v1/dashboard/graph/:id'
title: 'Get DashboardGraph info by id'
type: 'GET'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Response

```Status: 200```
```
{
    "counters":["value/name=pfc.push.ms", "value/name=pfc.push.size"],
    "endpoints":["laiweiofficemac"],
    "falcon_tags":"",
    "graph_id":4626,
    "graph_type":"h",
    "method":"",
    "position":4626,
    "screen_id":953,
    "timespan":3600,
    "title":"test"
}```


For errors responses, see the [response status codes documentation](#/response-status-codes).
