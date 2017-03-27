---
category: DashboardGraph
apiurl: '/api/v1/dashboard/tmpgraph'
title: 'Create a tmpgraph'
type: 'POST'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Request

``` {"endpoints":["e1", "e2"], "counters":["c1", "c2"]} ```

### Response

```Status: 200```

```
{
    "ck": "68c07419dbd7ac65977c97d05d99440d",
    "counters": "c1|c2",
    "endpoints": "e1|e2",
    "id": 365195
}
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
