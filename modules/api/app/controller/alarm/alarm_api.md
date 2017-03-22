# OWL alarm 文件

### `GET` / `POST` http://localhost:8088/api/v1/alarm/eventcases

| prams | description | example | 备注 |  
| ----- | ----------- | ------- | --- |
| startTime | 开始时间 | 1466956800 | 与endTime必须一起设置 |
| endTime | 结束时间 | 1467043200 | 与startTime必须一起设置 |
| priority | 告警等级 | 0 | 0 ~ 4 |
| status | 告警状态 | "OK" OR "OK,PROBLEM" | 可查询参数: "OK", "PROBLEM" |
| process_status | 人工设置状态 | "unresolved" OR "ignored,unresolved" | 可查询参数 "unresolved", "in progress", "recovered", "resolved" |
| metric | 监控项过滤 | "cpu.+"| 支援regexp查询 |
| id | 特定告警id | "s_104_0060710cedc48126d38aaa99447adda6" | 选取特定告警 |
| limit | 返回上限值 | 50 | 当page為空or-1时代表不分页返回, 上限為2000; 如page>=0时,上限為50.数值超过会被取代 |
| page | 分页 | 0 | 选定返回页面页数, -1表示不分页 |

ex.
```json
{
	"startTime": 1466956800,
	"endTime": 1467043200,
	"priority": 0,
	"status": "OK,PROBLEM",
	"process_status": "ignored,unresolved",
	"metrics": "cpu.idle",
	"id": "",
	"limit": 10,
	"page": 0
}
```
