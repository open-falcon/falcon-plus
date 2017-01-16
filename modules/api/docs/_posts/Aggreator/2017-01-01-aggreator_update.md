---
category: Aggreator
apiurl: '/api/v1/aggregators'
title: "Update Aggreator"
type: 'PUT'
sample_doc: 'aggreator.html'
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
  "id": 16,
  "endpoint": "testenp",
  "denominator": "$#"
}```

### Response

```Status: 200```
```{
  "id": 16,
  "grp_id": 343,
  "numerator": "$(cpu.idle)",
  "denominator": "$#",
  "endpoint": "testenp",
  "metric": "test.idle",
  "tags": "",
  "ds_type": "GAUGE",
  "step": 60,
  "creator": "root"
}```
