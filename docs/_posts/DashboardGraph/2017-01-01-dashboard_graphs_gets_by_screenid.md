---
category: DashboardGraph
apiurl: '/api/v1/dashboard/graphs/screen/:screen_id'
title: 'Gets graphs by screen id'
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
        "counters": [
            "value/name=pfc.push.ms"
        ],
        "endpoints": [
            "laiweiofficemac"
        ],
        "falcon_tags": "",
        "graph_id": 4640,
        "graph_type": "h",
        "method": "",
        "position": 0,
        "screen_id": 991,
        "timespan": 3600,
        "title": "dddd"
    },
    {
        "counters": [
            "aaa"
        ],
        "endpoints": [
            "xxx"
        ],
        "falcon_tags": "",
        "graph_id": 4641,
        "graph_type": "h",
        "method": "SUM",
        "position": 0,
        "screen_id": 991,
        "timespan": 3600,
        "title": "dddd"
    }
]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
