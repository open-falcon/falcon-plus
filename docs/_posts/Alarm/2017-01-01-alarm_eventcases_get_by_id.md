---
category: Alarm
apiurl: '/api/v1/alarm/eventcases'
title: 'Get EventCases by id'
type: 'GET'
sample_doc: 'alarm.html'
layout: default
---

* [Session](#/authentication) Required

### Request
Content-type: application/x-www-form-urlencoded

```event_id=s_165_cef145900bf4e2a4a0db8b85762b9cdb ```

### Response

```Status: 200```
```
    [
        {
            "closed_at": null,
            "closed_note": "",
            "cond": "0 != 66",
            "current_step": 3,
            "endpoint": "agent2",
            "expression_id": 0,
            "func": "all(#1)",
            "id": "s_165_cef145900bf4e2a4a0db8b85762b9cdb",
            "metric": "cpu.idle",
            "note": "\u5a13\ue103\u2502\u935b\u5a45\ue11f\u9477\ue044\u5aca\u93c7\u5b58\u67ca",
            "priority": 0,
            "process_note": 56603,
            "process_status": "ignored",
            "status": "PROBLEM",
            "step": 300,
            "strategy_id": 165,
            "template_id": 45,
            "timestamp": "2017-03-23T15:51:11+08:00",
            "tpl_creator": "root",
            "update_at": "2016-06-23T05:00:00+08:00",
            "user_modified": 0
        }
    ]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
