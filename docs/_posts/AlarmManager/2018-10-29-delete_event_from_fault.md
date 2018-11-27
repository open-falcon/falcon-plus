---
category: AlarmManager 
apiurl: '/api/v1/fault/:id/event'
title: "Delete event from fault"
type: 'DELETE'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required.
* The eventids in query string is divided by comma.

### Response

```Status:200```
```{
    "Code": 200,
    "Data": {
        "Id": 93,
        "CreatedAt": "2018-11-12T17:52:56+08:00",
        "Title": "mipush service down",
        "Note": "test",
        "Creator": "root",
        "Owner": "Bob",
        "State": "PROCESSING",
        "Tags": [
            "miphone",
            "miui"
        ],
        "Events": [
            {
                "event_id": 3,
                "eventcase_id": "s_11_d206cd6cbaa6b158c6894426e34d7f4e",
                "endpoint": "c3-op-mon-test01.bj",
                "counter": "agent.alive",
                "func": "all(#1)",
                "cond": "1 == 1",
                "note": "agent.alive测试",
                "max_step": 10000,
                "current_step": 1625,
                "priority": 0,
                "status": "PROBLEM",
                "event_ts": "2018-11-03T17:53:00+08:00",
                "template_creator": "jingtaoli",
                "expression_id": 0,
                "strategy_id": 11,
                "template_id": 7
            },
            {
                "event_id": 2,
                "eventcase_id": "s_11_d206cd6cbaa6b158c6894426e34d7f4e",
                "endpoint": "c3-op-mon-test01.bj",
                "counter": "agent.alive",
                "func": "all(#1)",
                "cond": "1 == 1",
                "note": "agent.alive测试",
                "max_step": 10000,
                "current_step": 1624,
                "priority": 0,
                "status": "PROBLEM",
                "event_ts": "2018-11-03T17:48:00+08:00",
                "template_creator": "jingtaoli",
                "expression_id": 0,
                "strategy_id": 11,
                "template_id": 7
            }
        ],
        "Followers": [
            "root"
        ],
        "Comments": []
    },
    "Message": "Delete event from fault successfully"
}```

