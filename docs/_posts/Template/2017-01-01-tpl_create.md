---
category: Template
apiurl: '/api/v1/template'
title: "Create Template"
type: 'POST'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* parent_id: 继承现有Template

### Request
```{"parent_id":0,"name":"AtmpForTesting"}```

### Response

```Status: 200```
```{
  "id": 2,
  "parent_id": 0,
  "tpl_name": "AtmpForTesting",
  "action_id": 0,
  "create_user": "root"
}```
