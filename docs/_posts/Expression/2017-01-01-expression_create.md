---
category: Expression
apiurl: '/api/v1/expression'
title: "Create Expression"
type: 'POST'
sample_doc: 'expression.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```{
  "right_value": "0",
  "priority": 2,
  "pause": 0,
  "op": "==",
  "note": "this is a test exp",
  "max_step": 3,
  "func": "all(#3)",
  "expression": "each(metric=agent.alive endpoint=docker-agent)",
  "action": {
    "url": "http://localhost:1234/callback",
    "uic": [
      "test"
    ],
    "callback": 1,
    "before_callback_sms": 1,
    "before_callback_mail": 0,
    "after_callback_sms": 1,
    "after_callback_mail": 0
  }
}```

### Response

```Status: 200```
```{"message":"expression created"}```
