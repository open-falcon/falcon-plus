---
category: Strategy
apiurl: '/api/v1/strategy'
title: "Get Strategy List"
type: 'GET'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required

### Response

```Status: 200```
```[
  {
    "id": 893,
    "metric": "process.num",
    "tags": "name=redis",
    "max_step": 3,
    "priority": 2,
    "func": "all(#2)",
    "op": "<",
    "right_value": "1",
    "note": "Redis异常",
    "run_begin": "",
    "run_end": "",
    "tpl_id": 221
  },
  {
    "id": 894,
    "metric": "process.num",
    "tags": "name=smtp",
    "max_step": 3,
    "priority": 2,
    "func": "all(#3)",
    "op": "<",
    "right_value": "1",
    "note": "Smtp异常",
    "run_begin": "",
    "run_end": "",
    "tpl_id": 221
  },
  {
    "id": 895,
    "metric": "process.num",
    "tags": "cmdline=logger",
    "max_step": 3,
    "priority": 3,
    "func": "all(#5)",
    "op": "<",
    "right_value": "2",
    "note": "logger异常",
    "run_begin": "",
    "run_end": "",
    "tpl_id": 221
  },
]```
