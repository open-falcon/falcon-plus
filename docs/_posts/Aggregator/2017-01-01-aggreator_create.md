---
category: Aggregator
apiurl: '/api/v1/aggregator'
title: "Create Aggregator to a HostGroup"
type: 'POST'
sample_doc: 'aggregator.html'
layout: default
---

* [Session](#/authentication) Required
* numerator: 分子
* denominator: 分母
* step: 汇报周期（秒为单位）

### Request

```{
  "tags": "",
  "step": 60,
  "numerator": "$(cpu.idle)",
  "metric": "test.idle",
  "hostgroup_id": 343,
  "endpoint": "testenp",
  "denominator": "2"
}```

### Response

```Status: 200```
```{
  "id": 16,
  "grp_id": 343,
  "numerator": "$(cpu.idle)",
  "denominator": "2",
  "endpoint": "testenp",
  "metric": "test.idle",
  "tags": "",
  "ds_type": "GAUGE",
  "step": 60,
  "creator": "root"
}```
