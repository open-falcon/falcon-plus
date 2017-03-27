---
category: Alarm
apiurl: '/api/v1/alarm/event_note'
title: 'Get Event Note by id or time range'
type: 'GET'
sample_doc: 'alarm.html'
layout: default
---

* [Session](#/authentication) Required

### Request

Content-type: application/x-www-form-urlencoded
```endTime=1466697600&startTime=1466611200```
or
```event_id=s_165_cef145900bf4e2a4a0db8b85762b9cdb```


### Response

```Status: 200```
```
    [
        {
            "case_id": "",
            "event_caseId": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
            "note": "test",
            "status": "ignored",
            "timestamp": "2016-06-23T05:39:09+08:00",
            "user": "root"
        },
        {
            "case_id": "",
            "event_caseId": "s_165_9d223f126e7ecb3477cd6806f1ee9656",
            "note": "Ignored by user",
            "status": "ignored",
            "timestamp": "2016-06-23T05:38:56+08:00",
            "user": "root"
        }
    ]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
