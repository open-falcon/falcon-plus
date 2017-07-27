---
category: HostGroup
apiurl: '/api/v1/hostgroup/update/#{hostgroup_id}'
title: "Update HostGroup info by id"
type: 'PUT'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup

### Request
```{
  "id" : 343,
  "grp_name": "test1"
}```

### Response

```Status: 200```
```{
  "message":"hostgroup profile updated"
}```
