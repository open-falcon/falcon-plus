---
category: Expression
apiurl: '/api/v1/expression/#{expression_id}'
title: "Get Expression Info by id"
type: 'GET'
sample_doc: 'expression.html'
layout: default
---

* [Session](#/authentication) Required
ex. /api/v1/expression/5


### Response

```Status: 200```
```{
  "action": {
    "id": 5,
    "uic": "taipei",
    "url": "",
    "callback": 0,
    "before_callback_sms": 0,
    "before_callback_mail": 0,
    "after_callback_sms": 0,
    "after_callback_mail": 0
  },
  "expression": {
    "id": 5,
    "expression": "each(metric=agent.alive endpoint=docker-agent)",
    "func": "all(#3)",
    "op": "==",
    "right_value": "0",
    "max_step": 3,
    "priority": 2,
    "note": "this is a test exp",
    "action_id": 177,
    "create_user": "root",
    "pause": 1
  }
}```
