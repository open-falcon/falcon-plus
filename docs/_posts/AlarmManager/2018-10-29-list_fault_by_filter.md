---
category: AlarmManager
apiurl: '/api/v1/fault'
title: "List fault by filter"
type: 'GET'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required.
* Filter field in url
    * start: Start  which in the format of unix timestamp in second is used to filter fault by created time of fault and is Now()-24h by default. 
    * end: End  which in the format of unix timestamp in second is used to filter fault by created time of fault and is Now() by default. 
    * creator: Creator is used to filter fault by creator of fault.
    * owner: Owner is used to filter fault by owner of fault.
    * state: State is used to filter fault by state of fault.
    * title: Title is used to filter fault by title of fault.
    * follower: Follower is used to filter fault by follower of fault.
    * tag: Tag is used to filter fault by tag of fault.
    * limit: Limit is used to pagination and is in the range of 0 to 50. By default, limit is 10.
    * offset: Offset is used to pagination and is 0 by default.
* Params in response body
    * Count is the number of fault which meets the filter rather than the number of fault in response body.

### Request url

```/api/v1/fault?start=1542010140&end=1542103740&creator=root&tag=miui&limit=1&offset=0```

### Response

```Status:200```
```{
    {
    "Code": 200,
    "Data": {
        "Faults": [
            {
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
            }
        ],
        "Count": 2
    },
    "Message": "List fault successfully"
}
}```

