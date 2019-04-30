---
category: Graph
apiurl: '/api/v1/graph/endpoint_counter'
title: "Get Counter of Endpoint"
type: 'GET'
sample_doc: 'graph.html'
layout: default
---

* [Session](#/authentication) Required
* params:
  * eid: endpoint id list
  * metricQuery: 查询counter参数【选填】。如：metricQuery=device=sda
  * page: 【选填】分页查询的页码，默认值：1,如： page=2 表示第2页
  * limit: 【选填】分页查询的页大小，默认值：500，如：limit=10 表示每页10条数据
  * q: 使用 regex 查询字符
    * option 参数

### Response

```Status: 200```
```[
  "disk.io.avgqu-sz/device=sda",
  "disk.io.ios_in_progress/device=sda",
  "disk.io.msec_read/device=sda",
  "disk.io.read_requests/device=sda",
  ...
]```
