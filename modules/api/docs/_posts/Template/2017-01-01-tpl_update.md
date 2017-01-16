---
category: Template
apiurl: '/api/v1/template/'
title: "Update Template"
type: 'PUT'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* parent_id: 继承现有Template

### Request
```{
  "tpl_id": 225,
  "parent_id": 0,
  "name": "AtmpForTesting2"
}```

### Response

```Status: 200```
```{"message":"template updated"}```
