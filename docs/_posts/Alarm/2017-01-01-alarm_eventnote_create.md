---
category: Alarm
apiurl: '/api/v1/alarm/event_note'
title: 'Create Event Note'
type: 'POST'
sample_doc: 'alarm.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```
    {
        "event_id": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
        "note": "test note",
        "status": "comment"
    }
```

### Response

```Status: 200```
```
    {
        "id": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
        "message": "add note to s_165_cef145900bf4e2a4a0db8b85762b9cdb successfuled"
    }
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
