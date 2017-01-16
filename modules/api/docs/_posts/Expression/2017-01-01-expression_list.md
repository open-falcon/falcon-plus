---
category: Expression
apiurl: '/api/v1/expression'
title: "Expression List"
type: 'GET'
sample_doc: 'expression.html'
layout: default
---

* [Session](#/authentication) Required

### Response

```Status: 200```
```[
  {
    "id": 2,
    "expression": "each(metric=? xx=yy)",
    "func": "all(#3)",
    "op": "==",
    "right_value": "0",
    "max_step": 3,
    "priority": 0,
    "note": "",
    "action_id": 18,
    "create_user": "root",
    "pause": 0
  },
  {
    "id": 3,
    "expression": "each(metric=ss.close.wait endpoint=docker-A)",
    "func": "all(#1)",
    "op": "!=",
    "right_value": "0",
    "max_step": 1,
    "priority": 4,
    "note": "boss docker-A 连接数大于10",
    "action_id": 91,
    "create_user": "root",
    "pause": 0
  },
  {
    "id": 4,
    "expression": "each(metric=agent.alive endpoint=docker-agent)",
    "func": "all(#3)",
    "op": "==",
    "right_value": "0",
    "max_step": 3,
    "priority": 2,
    "note": "this is a test exp",
    "action_id": 176,
    "create_user": "root",
    "pause": 1
  }
]```
