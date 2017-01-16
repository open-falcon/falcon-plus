---
category: Strategy
apiurl: '/api/v1/strategy/#{strategy_id}'
title: "Get Strategy info by id"
type: 'GET'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/strategy/904

### Response

```Status: 200```
```{
  "id": 904,
  "metric": "agent.alive",
  "tags": "",
  "max_step": 3,
  "priority": 1,
  "func": "all(#3)",
  "op": "==",
  "right_value": "1",
  "note": "this is a test",
  "run_begin": "00:00",
  "run_end": "24:00",
  "tpl_id": 221
}```
