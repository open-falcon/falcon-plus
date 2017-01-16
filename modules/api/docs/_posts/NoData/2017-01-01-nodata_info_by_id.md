---
category: NoData
apiurl: '/api/v1/nodata/#{nodata_id}'
title: "Get Nodata Info by id"
type: 'GET'
sample_doc: 'nodata.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/nodata/2

### Response

```Status: 200```
```{
  "id": 2,
  "name": "owl_nodate",
  "obj": "docker-agent",
  "obj_type": "host",
  "metric": "test.metric",
  "tags": "",
  "dstype": "GAUGE",
  "step": 60,
  "mock": -2,
  "creator": "root"
}```
