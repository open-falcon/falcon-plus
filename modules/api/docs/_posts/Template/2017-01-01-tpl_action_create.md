---
category: Template
apiurl: '/api/v1/template/action'
title: "Create Template Action"
type: 'POST'
sample_doc: 'template.html'
layout: default
---

Create Action to a Template
* [Session](#/authentication) Required
* params:
  * url: callback url
  * uic: 需要通知的使用者群组(name)
  * callback: enable/disable

### Request
```{
  "url": "",
  "uic": "test,tt2",
  "tpl_id": 225,
  "callback": 1,
  "before_callback_sms": 0,
  "before_callback_mail": 0,
  "after_callback_sms": 0,
  "after_callback_mail": 0
}```

### Response

```Status: 200```
```{"message":"template created"}```
