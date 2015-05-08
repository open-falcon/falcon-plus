## intro

Query面向终端用户，收到查询请求后，根据一致性哈希算法，会去相应的Graph里面，查询不同metric的数据，汇总后统一返回给用户。

## python query 历史数据的例子

```python
#!-*- coding:utf8 -*-

import requests
import time
import json

end = int(time.time())
start = end - 3600  #查询过去一小时的数据

d = {
        "start": start,
        "end": end,
        "cf": "AVERAGE",
        "endpoint_counters": [
            {
                "endpoint": "host1",
                "counter": "cpu.idle",
            },
            {
                "endpoint": "host1",
                "counter": "load.1min",
            },
        ],
}

url = "http://127.0.0.1:9966/graph/history"
r = requests.post(url, data=json.dumps(d))
print r.text

```

其中cf的值可以为：AVERAGE、MAX、MIN ，具体可以参考RRDtool的相关概念

## install

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/query.git

cd query
go get ./...
./control build

./control start
```

## configuration

    log_level: 可选 error/warn/info/debug/trace，默认为info

    slow_log: 单位是毫秒，query的时候，较慢的请求，会被打印到日志中，默认是2000ms

    debug: true/false, 如果为true，日志中会打印debug信息

    http
        - enable: true/false, 表示是否开启该http端口，该端口为数据查询接口（即提供http api给用户查询数据）
        - listen: 表示监听的http端口

    graph
        - backends: 后端graph列表文件，格式参考下面的介绍，该文件默认是./graph_backends.txt
        - reload_interval: 单位是秒，表示每隔多久自动reload一次backends列表文件中的内容
        - timeout: 单位是毫秒，表示和后端graph组件交互的超时时间，可以根据网络质量微调，建议保持默认
        - pingMethod: 后端提供的ping接口，用来探测连接是否可用，必须保持默认
        - max_conns: 连接池相关配置，最大连接数，建议保持默认
        - max_idle: 连接池相关配置，最大空闲连接数，建议保持默认
        - replicas: 这是一致性hash算法需要的节点副本数量，建议不要变更，保持默认即可

## backends 文件格式
1. 每行由空格分割的两列组成，第一列表示graph的名字，第二列表示具体的ip:port
2. 该文件需要和transfer配置文件中的cluster的配置项，保持一致

```bash
$ cat ./graph_backends.txt

graph-00 127.0.0.1:6070
graph-01 127.0.0.2:6070
```
