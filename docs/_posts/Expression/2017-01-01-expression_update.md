---
category: Expression
apiurl: '/api/v1/expression'
title: "Update Expression"
type: 'PUT'
sample_doc: 'expression.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```{
  "right_value": "0",
  "priority": 2,
  "pause": 1,
  "op": "==",
  "note": "this is a test exp",
  "max_step": 3,
  "id": 5,
  "func": "all(#3)",
  "expression": "each(metric=agent.alive endpoint=docker-agent)",
  "action": {
    "url": "http://localhost:1234/callback",
    "uic": [
      "test",
      "test2"
    ],
    "callback": 0,
    "before_callback_sms": 1,
    "before_callback_mail": 0,
    "after_callback_sms": 1,
    "after_callback_mail": 0
  }
}```

### Response

```Status: 200```
```{"message":"expression:5 has been updated"}```
