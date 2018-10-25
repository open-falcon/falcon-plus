---
category: Aggregator
apiurl: '/api/v1/hostgroup/#{hostgroup_id}/aggregators'
title: "Get Aggregator List of HostGroup"
type: 'GET'
sample_doc: 'aggregator.html'
layout: default
---

* [Session](#/authentication) Required
* ex. /api/v1/hostgroup/343/aggregators
* numerator: 分子
* denominator: 分母
* step: 汇报周期（秒为单位）

### Response

```Status: 200```
```[
  {
    "id": 13,
    "grp_id": 343,
    "numerator": "$(cpu.idle)",
    "denominator": "2",
    "endpoint": "testenp",
    "metric": "test.idle",
    "tags": "",
    "ds_type": "GAUGE",
    "step": 60,
    "creator": "root"
  },
  {
    "id": 14,
    "grp_id": 343,
    "numerator": "$(cpu.idle)",
    "denominator": "2",
    "endpoint": "testenp",
    "metric": "test.idle",
    "tags": "",
    "ds_type": "GAUGE",
    "step": 60,
    "creator": "root"
  }
]```
