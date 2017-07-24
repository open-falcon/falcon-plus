---
category: Host
apiurl: '/api/v1/host/maintain'
title: "Set host maintain by ids or hostnames"
type: 'POST'
sample_doc: 'host.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```
{
    "ids": [1,2,3,4],
    "maintain_begin": 1497951907,
    "maintain_end": 1497951907
}
```
or
```
{
    "hosts": ["host.a","host.b"],
    "maintain_begin": 1497951907,
    "maintain_end": 1497951907
}
```

### Response

```Status: 200```
```{ "message": "Through: hosts, Affect row: 2" }```
