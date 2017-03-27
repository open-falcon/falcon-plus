---
category: NoData
apiurl: '/api/v1/nodata/'
title: "Create Nodata"
type: 'POST'
sample_doc: 'nodata.html'
layout: default
---

* [Session](#/authentication) Required

### Request

```{
  "tags": "",
  "step": 60,
  "obj_type": "host",
  "obj": "docker-agent",
  "name": "testnodata",
  "mock": -1,
  "metric": "test.metric",
  "dstype": "GAUGE"
}```

### Response

```Status: 200```
```{
  "id": 4,
  "name": "testnodata",
  "obj": "docker-agent",
  "obj_type": "host",
  "metric": "test.metric",
  "tags": "",
  "dstype": "GAUGE",
  "step": 60,
  "mock": -1,
  "creator": "root"
}```
