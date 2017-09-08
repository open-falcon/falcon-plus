---
category: Alarm
apiurl: '/api/v1/alarm/eventcases'
title: 'EventCases List'
type: 'POST'
sample_doc: 'alarm.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```
    {
        "endTime": 1480521600,
        "limit": 10,
        "process_status": "ignored,unresolved",
        "startTime": 1466956800,
        "status": "PROBLEM",
        "endpoints": ["agent4"],
        "strategy_id": 46,
        "template_id": 126
    }
```

### Response

```Status: 200```
```
    {
        "closed_at": null,
        "closed_note": "",
        "cond": "48.33759590792839 > 40",
        "current_step": 1,
        "endpoint": "agent4",
        "expression_id": 0,
        "func": "all(#3)",
        "id": "s_46_1ac45122afb893adc02fbd30154ac303",
        "metric": "cpu.iowait",
        "note": "CPU I/O wait\u74d2\u5470\u7e4340",
        "priority": 1,
        "process_note": 16907,
        "process_status": "ignored",
        "status": "PROBLEM",
        "step": 1,
        "strategy_id": 46,
        "template_id": 126,
        "timestamp": "2016-08-01T06:25:00+08:00",
        "tpl_creator": "root",
        "update_at": "2016-08-01T06:25:00+08:00",
        "user_modified": 0
    }
]
```

For errors responses, see the [response status codes documentation](#/response-status-codes).
