---
category: HostGroup
apiurl: '/api/v1/hostgroup/host'
title: "Unbind a Host on HostGroup"
type: 'PUT'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* 如果使用者不是 Admin 只能对创建的hostgroup做操作

### Request
```{
  "hostgroup_id": 343,
  "host_id": 9312
}```

### Response

```Status: 200```
```{"message":"unbind host:9312 of hostgroup: 343"}```
