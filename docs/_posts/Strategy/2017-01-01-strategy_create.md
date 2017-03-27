---
category: Strategy
apiurl: '/api/v1/strategy'
title: "Create Strategy"
type: 'POST'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required


### Request
```{
  "tpl_id": 221,
  "tags": "",
  "run_end": "24:00",
  "run_begin": "00:00",
  "right_value": "1",
  "priority": 1,
  "op": "==",
  "note": "this is a test",
  "metric": "agent.alive",
  "max_step": 3,
  "func": "all(#3)"
}```

### Response

```Status: 200```
```{"message":"stragtegy created"}```
