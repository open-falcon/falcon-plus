---
category: DashboardGraph
apiurl: '/api/v1/dashboard/graph/:id'
title: 'Update a DashboardGraph'
type: 'PUT'
sample_doc: 'dashboard.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```
{
    "counters": ["value/name=pfc.push.ms","value/name=pfc.push.size", "agent.alive"],
    "falcon_tags": "srv=falcon"
}
```

### Response

```Status: 200```

```{"message":"ok"}```

For errors responses, see the [response status codes documentation](#/response-status-codes).
