---
category: Metric
apiurl: '/api/v1/metric/default_list'
title: "Get Default Builtin Metric List"
type: 'GET'
sample_doc: 'template.html'
layout: default
---

* [Session](#/authentication) Required
* Metric Suggestion for create strategy
* base on `./data/metric` file

### Response

```Status: 200```
```[
  "cpu.busy",
  "cpu.cnt",
  "cpu.guest",
  "cpu.idle",
  "cpu.iowait",
  "cpu.irq",
  "cpu.nice",
  "cpu.softirq",
  "cpu.steal",
  "cpu.system",
  "cpu.user",
  "df.bytes.free",
  "df.bytes.free.percent"
  ....
]```
