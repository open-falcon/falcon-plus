---
category: HostGroup
apiurl: '/api/v1/hostgroup/template'
title: "Bind A Template to HostGroup"
type: 'POST'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```{
  "tpl_id": 5,
  "grp_id": 3
}```

### Response

```Status: 200```
```{"grp_id":3,"tpl_id":5,"bind_user":2}```
