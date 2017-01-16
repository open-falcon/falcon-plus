---
category: Host
apiurl: '/api/v1/host/#{host_id}/template'
title: "Get bind Template List of Host"
type: 'GET'
sample_doc: 'host.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/host/1647/template
* tpl_name: template name

### Response

```Status: 200```
```[
  {
    "id": 125,
    "tpl_name": "tplA",
    "parent_id": 0,
    "action_id": 99,
    "create_user": "root"
  },
  {
    "id": 142,
    "tpl_name": "tplB",
    "parent_id": 0,
    "action_id": 111,
    "create_user": "root"
  },
  {
    "id": 180,
    "tpl_name": "tplC",
    "parent_id": 0,
    "action_id": 142,
    "create_user": "root"
  }
]```
