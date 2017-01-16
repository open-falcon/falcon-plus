---
category: NoData
apiurl: '/api/v1/nodata/'
title: "Update Nodata"
type: 'PUT'
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
  "mock": -2,
  "metric": "test.metric",
  "id": 4,
  "dstype": "GAUGE"
}```

### Response

```Status: 200```
```{
  "id": 0,
  "name": "",
  "obj": "docker-agent",
  "obj_type": "host",
  "metric": "test.metric",
  "tags": "",
  "dstype": "GAUGE",
  "step": 60,
  "mock": -2,
  "creator": ""
}```
