---
category: Template
apiurl: '/api/v1/template/#{template_id}/hostgroup'
title: "Get hostgroups list by id"
type: 'GET'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/template/178/hostgroup

### Response

```Status: 200```
```{
  "hostgroups": [{
      "id":33,
      "grp_name":"HostGroup",
      "create_user":"root"
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
