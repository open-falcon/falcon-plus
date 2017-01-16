---
category: HostGroup
apiurl: '/api/v1/hostgroup/#{hostgroup_id}/template'
title: "Get Template List of  HostGroup"
type: 'GET'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup/3/template

### Response

```Status: 200```
```{
  "hostgroup": {
    "id": 3,
    "grp_name": "hostgroupA",
    "create_user": "root"
  },
  "templates": [
    {
      "id": 5,
      "tpl_name": "TplA",
      "parent_id": 0,
      "action_id": 12,
      "create_user": "root"
    },
    {
      "id": 91,
      "tpl_name": "TplB",
      "parent_id": 0,
      "action_id": 59,
      "create_user": "userA"
    },
    {
      "id": 94,
      "tpl_name": "TplB",
      "parent_id": 0,
      "action_id": 62,
      "create_user": "userA"
    },
    {
      "id": 103,
      "tpl_name": "TplC",
      "parent_id": 0,
      "action_id": 74,
      "create_user": "root"
    },
    {
      "id": 104,
      "tpl_name": "TplD",
      "parent_id": 0,
      "action_id": 75,
      "create_user": "root"
    },
    {
      "id": 105,
      "tpl_name": "TplE",
      "parent_id": 0,
      "action_id": 76,
      "create_user": "root"
    },
    {
      "id": 116,
      "tpl_name": "TplG",
      "parent_id": 0,
      "action_id": 87,
      "create_user": "root"
    },
    {
      "id": 125,
      "tpl_name": "TplH",
      "parent_id": 0,
      "action_id": 99,
      "create_user": "root"
    },
    {
      "id": 126,
      "tpl_name": "http",
      "parent_id": 0,
      "action_id": 100,
      "create_user": "root"
    },
    {
      "id": 127,
      "tpl_name": "TplJ",
      "parent_id": 0,
      "action_id": 101,
      "create_user": "root"
    }
  ]
}```
