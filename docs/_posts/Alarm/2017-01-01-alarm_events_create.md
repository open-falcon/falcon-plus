---
category: Alarm
apiurl: '/api/v1/alarm/events'
title: 'Create Events'
type: 'POST'
sample_doc: 'alarm.html'
layout: default
---

* [Session](#/authentication) Required

### Request

Content-type: application/x-www-form-urlencoded

Key|Value
---|---
endTime|1466628960
event_id|s_165_cef145900bf4e2a4a0db8b85762b9cdb
startTime|1466611200

### Response

```Status: 200```
```
    [
        {
            "cond": "10.649350649350648 != 66",
            "event_caseId": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
            "id": 635166,
            "status": 0,
            "step": 0,
            "timestamp": "2016-06-23T04:55:00+08:00"
        },
        {
            "cond": "13.486005089058525 != 66",
            "event_caseId": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
            "id": 635149,
            "status": 0,
            "step": 0,
            "timestamp": "2016-06-23T04:50:00+08:00"
        }
    ]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
