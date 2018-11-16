---
category: AlarmManager
apiurl: '/api/v1/fault/:id/timeline'
title: "Get timeline of fault"
type: 'GET'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required.

### Response

```Status:200```
```{
    "Code": 200,
    "Data": {
        "FirstEventTime": "2018-11-03T17:44:00+08:00",
        "LastEventTime": "2018-11-03T17:53:00+08:00",
        "FaultCreatedAt": "2018-11-12T17:52:56+08:00",
        "FaultClosedAt": "2018-11-13T15:11:13+08:00"
    },
    "Message": "Get fault timeline successfully"
}```

