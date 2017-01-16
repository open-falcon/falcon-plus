---
category: HostGroup
apiurl: '/api/v1/hostgroup'
title: "Create HostGroup"
type: 'POST'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```{"name":"testhostgroup"}```

### Response

```Status: 200```
```{
  "id": 343,
  "grp_name": "testhostgroup",
  "create_user": "root"
}```
