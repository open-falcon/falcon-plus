---
category: HostGroup
apiurl: '/api/v1/hostgroup/#{hostgroup_id}'
title: "Get HostGroup info by id"
type: 'GET'
sample_doc: 'hostgroup.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup/343

### Response

```Status: 200```
```{
  "hostgroup": {
    "id": 343,
    "grp_name": "testhostgroup",
    "create_user": "root"
  },
  "hosts": [
    {
      "id": 9313,
      "hostname": "agent_test",
      "ip": "",
      "agent_version": "",
      "plugin_version": "",
      "maintain_begin": 0,
      "maintain_end": 0
    }
  ]
}```
