---
category: HostGroup
apiurl: '/api/v1/hostgroup/template'
title: "Unbind A Template of A HostGroup"
type: 'PUT'
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
```{"message":"template: 5 is unbind of HostGroup: 3"}```
