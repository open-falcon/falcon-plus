---
category: Template
apiurl: '/api/v1/template/#{template_id}'
title: "Get Template Info by id"
type: 'GET'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/template/178

### Response

```Status: 200```
```{
  "action": {
    "id": 141,
    "uic": "admin,mm1",
    "url": "",
    "callback": 0,
    "before_callback_sms": 0,
    "before_callback_mail": 0,
    "after_callback_sms": 0,
    "after_callback_mail": 0
  },
  "stratges": [
    {
      "id": 686,
      "metric": "xxx.check",
      "tags": "name=xxx",
      "max_step": 3,
      "priority": 2,
      "func": "all(#2)",
      "op": "<",
      "right_value": "1",
      "note": "xxx服务异常",
      "run_begin": "",
      "run_end": "",
      "tpl_id": 178
    },
    {
      "id": 687,
      "metric": "xxx.sync",
      "tags": "",
      "max_step": 3,
      "priority": 2,
      "func": "all(#3)",
      "op": "!=",
      "right_value": "0",
      "note": "XXX同步异常",
      "run_begin": "",
      "run_end": "",
      "tpl_id": 178
    },
    {
      "id": 688,
      "metric": "bbb.check.mq",
      "tags": "",
      "max_step": 3,
      "priority": 2,
      "func": "all(#3)",
      "op": "==",
      "right_value": "1",
      "note": "bbb连接MQ异常",
      "run_begin": "",
      "run_end": "",
      "tpl_id": 178
    },
    {
      "id": 793,
      "metric": "aaa.proc.num",
      "tags": "",
      "max_step": 3,
      "priority": 2,
      "func": "all(#3)",
      "op": "==",
      "right_value": "1",
      "note": "aaaa 进程大于5",
      "run_begin": "",
      "run_end": "",
      "tpl_id": 178
    }
  ],
  "template": {
    "id": 178,
    "tpl_name": "TemplateA",
    "parent_id": 0,
    "action_id": 141,
    "create_user": "root"
  }
}```
