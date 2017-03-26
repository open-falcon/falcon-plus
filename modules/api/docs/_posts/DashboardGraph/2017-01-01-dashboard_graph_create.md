---
category: DashboardGraph
apiurl: '/api/v1/dashboard/graph'
title: 'Create User'
type: 'POST'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```
{
    "screen_id": 953,
    "title": "laiwei-test-graph1",
    "endpoints": ["laiweiofficemac"],
    "counters": ["value/name=pfc.push.ms","value/name=pfc.push.size"],
    "timespan": 1800, 
    "graph_type": "h", 
    "method": "AVG",
    "position": 0
}
```

### Response

```Status: 200```

```{"message":"ok"}```

For errors responses, see the [response status codes documentation](#/response-status-codes).
