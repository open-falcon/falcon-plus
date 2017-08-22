---
category: HostGroup
apiurl: '/api/v1/hostgroup/#{hostgroup_id}/host'
title: "Update partial hosts in HostGroup"
type: 'PATCH'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* 如果使用者不是 Admin 只能对创建的 HostGroup 做操作
* ex. /api/v1/hostgroup/1/host

### Request
```{
  "hosts": [
    "host01",
    "host02"
  ],
  "action": "add"
}```

### Response

```Status: 200```
```{
  "message": "[host01, host02] bind to hostgroup: test, [] have been exist"
}```

### Request
```{
  "hosts": [
    "host01",
    "host02"
  ],
  "action": "remove"
}```

### Response

```Status: 200```
```{
  "message": "[host01, host02] unbind to hostgroup: test"
}```

