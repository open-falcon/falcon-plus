---
category: AlarmManager
apiurl: '/api/v1/event/count'
title: 'Get the Events, Or Contains the Faults Count By the User'
type: 'POST'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required
* 接口说明：
    * 通过报警接收者、或者报警组、或者接收者所在的告警组组合其它条件获取告警事件数量
    * 通过报警接收者、或者报警组、或者接收者所在的告警组组合其它条件获取已添加至故障的告警事件数量
* 参数说明：
    *   请求body包括: username、uic、start_time、end_time、priority、endpoint、counter、status、now_event_status、have_fault、limit、offset.
    * username、uic: 分别表示告警接收者、告警组；至少一个不为空
    * start_time、end_time: 查询告警信息的时间范围，默认最近一小时
    * priority: 告警优先级，支持多级别
    * endpoint: 告警主体
    * counter: 告警指标，metric+tags组合
    * status(ok/problem)、now_event_status(true/false): status表示告警事件状态，now_event_status表示是否当前告警事件状态。意义：now_event_status用于区分历史告警状态(告警事件追溯)和目前告警事件状态(例如目前未恢复的报警、已经恢复的报警)，当status为problem，now_event_status为true，表示当前未恢复的告警信息。
    * have_fault: 表示已添加至故障，表示当前那些告警event已经被添加至故障
    * limit、offset: 用于查询分页，默认limit 100，offset 0

### Request

```
{
        "username":  "xiaoming",
        "uic":  [
            "plus-dev01",
            "plus-dev02"
        ],
        "endpoint": "agent01",
        "counter":  "mem.memused",
        "priority": [
            "0"
        ],
        "start_time": 1535760000,
        "end_time": 1535767200,
        "status": "PROBLEM",
        "now_event_status": false，
        "have_fault": false
}
```

### Response
```Status: 200```

```
{
    "Code": 200,
    "Data": {
        "count": 100
    },
    "Message": ""
}
```


200 status code is returned if succeed,otherwise 400 or 500 status code is returned.
