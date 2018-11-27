---
category: AlarmManager
apiurl: '/api/v1/fault/:id'
title: "Get fault by id"
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
        "Id": 93,
        "CreatedAt": "2018-11-12T17:52:56+08:00",
        "Title": "mipush service down",
        "Note": "hello test",
        "Creator": "root",
        "Owner": "George",
        "State": "CLOSED",
        "Tags": [
            "micloud",
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
            },
            {
                "event_id": 1,
                "eventcase_id": "s_13_3e251836d0e4992f5fc39d8b65864ee6",
                "endpoint": "c3-op-mon-test01.bj",
                "counter": "mem.memfree.percent",
                "func": "all(#1)",
                "cond": "92.59889718652626 > 0",
                "note": "",
                "max_step": 5000,
                "current_step": 971,
                "priority": 3,
                "status": "PROBLEM",
                "event_ts": "2018-11-03T17:44:00+08:00",
                "template_creator": "jingtaoli",
                "expression_id": 0,
                "strategy_id": 13,
                "template_id": 7
            }
        ],
        "Followers": [
            "root",
            "David",
            "Jack"
        ],
        "Comments": []
    },
    "Message": "Get fault successfully"
}```

