---
category: Host
apiurl: '/api/v1/host/maintain'
title: "Reset host maintain by ids or hostnames"
type: 'DELETE'
sample_doc: 'host.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```{"ids": [1,2,3,4]}```
or
```{"hosts": ["host.a","host.b"]}```

### Response

```Status: 200```
```{ "message": "Through: hosts, Affect row: 2" }```
