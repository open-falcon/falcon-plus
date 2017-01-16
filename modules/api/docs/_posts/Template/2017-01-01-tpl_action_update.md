---
category: Template
apiurl: '/api/v1/template/action'
title: "Update Template Action"
type: 'PUT'
sample_doc: 'template.html'
layout: default
---

Update Action
* [Session](#/authentication) Required
* params:
  * url: callback url
  * uic: 需要通知的使用者群组(name)
  * callback: enable/disable

### Request
```{
  "url": "",
  "uic": "test,tt2,tt3",
  "id": 175,
  "callback": 1,
  "before_callback_sms": 0,
  "before_callback_mail": 0,
  "after_callback_sms": 0,
  "after_callback_mail": 0
}```

### Response

```Status: 200```
```{"message":"action is updated, row affected: 1"}```
