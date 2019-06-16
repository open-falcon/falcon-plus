---
category: Graph
apiurl: '/api/v1/graph/endpoint'
title: "Endpoint List"
type: 'GET'
sample_doc: 'graph.html'
layout: default
---

* [Session](#/authentication) Required
* params:
    * q: 使用 regex 查询字符
      * option 参数
    * page: 【选填】分页查询的页码，默认值：1,如： page=2 表示第2页
    * limit: 【选填】分页查询的页大小，默认值：500，如：limit=10 表示每页10条数据

### Response

```Status: 200```
```[
  {
    "endpoint": "docker-agent",
    "id": 7
  },
  {
    "endpoint": "docker-task",
    "id": 6
  },
  {
    "endpoint": "graph",
    "id": 3
  },
  {
    "endpoint": "nodata",
    "id": 920
  },
  {
    "endpoint": "task",
    "id": 5
  }
]```
