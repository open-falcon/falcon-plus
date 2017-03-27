---
category: HostGroup
apiurl: '/api/v1/hostgroup/#{hostgroup_id}'
title: "Delete HostGroup"
type: 'DELETE'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
ex. /api/v1/hostgroup/343
* 如果使用者不是 Admin 只能对创建的hostgroup做操作

### Response

```Status: 200```
```{"message":"hostgroup:343 has been deleted"}```
