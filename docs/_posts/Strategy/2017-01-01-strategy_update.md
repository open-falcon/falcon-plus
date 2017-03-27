---
category: Strategy
apiurl: '/api/v1/strategy'
title: "Update Strategy"
type: 'PUT'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required


### Request
```{
  "tags": "",
  "run_end": "",
  "run_begin": "",
  "right_value": "1",
  "priority": 2,
  "op": "==",
  "note": "this is a test",
  "metric": "agent.alive",
  "max_step": 3,
  "id": 904,
  "func": "all(#3)"
}```

### Response

```Status: 200```
```{"message":"stragtegy:904 has been updated"}```
