---
category: Host
apiurl: '/api/v1/host/hostname'
title: "Get the host message by hostname"
type: 'POST'
sample_doc: 'host.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```
{
    "hostname": ["host.a","host.b"]
}
```
or
```
{"hostname": ["host.a"]}
```

### Response

```Status: 200```
```
[
    {"id":81,"hostname":"host.a","ip":"192.168.1.1","agent_version":"5.1.2","plugin_version":"Error:exit status 128","maintain_begin":0,"maintain_end":0},
    {"id":85,"hostname":"host.b","ip":"192.168.1.2","agent_version":"5.1.2","plugin_version":"Error:exit status 128","maintain_begin":0,"maintain_end":0}
]
```
or
```
[
    {"id":81,"hostname":"host.a","ip":"192.168.1.1","agent_version":"5.1.2","plugin_version":"Error:exit status 128","maintain_begin":0,"maintain_end":0}
]
```
