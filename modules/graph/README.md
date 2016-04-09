## Introduction

对于监控系统来讲，历史数据的存储，高效率查询，快速展现，是个很重要且困难的问题。这主要表现在下面三个方面：

1. 数据量大：目前我们的监控系统，每个周期，大概有1亿次数据采集项上报（上报周期为1分钟和5分钟两种，各占50%），一天24小时里，从来不会有低峰，不管是白天和黑夜，每个周期，总会有那么多的数据要更新。
2. 写操作多：一般的业务系统，通常都是读多写少，可以方便的使用各种缓存技术，再者各类数据库，对于查询操作的处理效率远远高于写操作。而监控系统恰恰相反，写操作远远高于读。每个周期几千万次的更新操作，对于常用数据库（MySQL、PostgreSQL、MongoDB）都不是最合适和擅长的。
3. 高效率的查：我们说监控系统读操作少，是说相对写入来讲。监控系统本身对于读的要求很高，用户经常会有查询上百个metric，在过去一天、一周、一月、一年的数据。如何在秒级返回给用户并在前端展现，这是一个不小的挑战。

graph所做的事情，就是把用户每次push上来的数据，进行采样存储，并提供查询接口。

我们参考RRDtool的理念，在数据每次存入的时候，会自动进行采样、归档。在默认的归档策略，一分钟push一次的频率下，历史数据保存5年。同时为了不丢失信息量，数据归档的时候，会按照平均值采样、最大值采样、最小值采样存三份。用户在查询某个metric，在过去一年的历史数据时，Graph会选择最合适的采样频率，返回采样过后的数据，提高了数据查询速度。

## Installation

源码安装过程，如下

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/graph.git

cd graph
go get ./...
./control build

./control start
```

你可以从[这里](https://github.com/open-falcon/graph/releases)，下载最新的release安装包，避免源码编译的种种问题。

## Configuration

    {
        "debug": false, //true or false, 是否开启debug日志
        "http": {
            "enabled": true, //true or false, 表示是否开启该http端口，该端口为控制端口，主要用来对graph发送控制命令、统计命令、debug命令
            "listen": "0.0.0.0:6071" //表示监听的http端口
        },
        "rpc": {
            "enabled": true, //true or false, 表示是否开启该rpc端口，该端口为数据接收端口
            "listen": "0.0.0.0:6070" //表示监听的rpc端口
        },
        "rrd": {
            "storage": "/home/work/data/6070" //绝对路径，历史数据的文件存储路径（如有必要，请修改为合适的路）
        },
        "db": {
            "dsn": "root:@tcp(127.0.0.1:3306)/graph?loc=Local&parseTime=true", //MySQL的连接信息，默认用户名是root，密码为空，host为127.0.0.1，database为graph（如有必要，请修改)
            "maxIdle": 4  //MySQL连接池配置，连接池允许的最大连接数，保持默认即可
        },
        "callTimeout": 5000,  //RPC调用超时时间，单位ms
        "migrate": {  //扩容graph时历史数据自动迁移
            "enabled": false,  //true or false, 表示graph是否处于数据迁移状态
            "concurrency": 2, //数据迁移时的并发连接数，建议保持默认
            "replicas": 500, //这是一致性hash算法需要的节点副本数量，建议不要变更，保持默认即可（必须和transfer的配置中保持一致）
            "cluster": { //未扩容前老的graph实例列表
                "graph-00" : "127.0.0.1:6070"
            }
        }
    }

## 关于扩容时数据自动迁移

当graph集群扩容时，数据会自动迁移达到rebalance的目的。具体的操作步骤如下：

假设1：旧的graph集群为
```
graph-00 : 192.168.1.1:6070
graph-01 : 192.168.1.2:6070
```

假设2：在已有的graph集群基础上，需要再扩容两个实例，如下
```
graph-00 : 192.168.1.1:6070
graph-01 : 192.168.1.2:6070
graph-02 : 192.168.1.3:6070
graph-03 : 192.168.1.4:6070
```

####1 修改新增加的graph的cfg.json如下：
```python
"migrate": {
      "enabled": true,  //true 表示graph是否处于数据迁移状态
      "concurrency": 2,
      "replicas": 500,
      "cluster": { //未扩容前老的graph实例列表
         "graph-00" : "192.168.1.1:6070",
         "graph-01" : "192.168.1.2:6070"
      }
}
```
> 要点说明:

> 1. 只有新增加的graph实例配置需要修改，原有的graph实例配置不能变动

> 2. 修改配置文件的时候，首先要把migrate -> enabled这个开关修改为true；其次要保证migrate -> cluster这个列表为旧集群的graph列表。


####2 重启所有的graph

```cd $WORKSPACE/graph/ && bash control restart```


####3 修改所有transfer的cfg.json如下：

```python
"graph": {
    "enabled": true,
     ...,
    "replicas": 500,
    "cluster": {
        "graph-00" : "192.168.1.1:6070",
        "graph-01" : "192.168.1.2:6070",
        "graph-02" : "192.168.1.3:6070",
        "graph-03" : "192.168.1.4:6070"
    }
},
```
> 要点说明：最主要的就是修改graph -> cluster列表为扩容后完整的graph实例列表

####4 重启所有的transfer

```cd $WORKSPACE/transfer/ && bash control restart```

> 这时候，transfer就会将接收到的数据，发送给扩容后的graph实例；同时graph实例，会自动进行数据的rebalance，rebalance的过程持续时间长短，与待迁移的counter数量以及graph机器的负载、性能有关系。

####5 修改query的配置，并重启所有的query进程
```python
"graph": {
     ...,
    "replicas": 500,
    "cluster": {
        "graph-00" : "192.168.1.1:6070",
        "graph-01" : "192.168.1.2:6070",
        "graph-02" : "192.168.1.3:6070",
        "graph-03" : "192.168.1.4:6070"
    }
},
```
> 要点说明：最主要的就是修改graph -> cluster列表为扩容后完整的graph实例列表



####6 如何确认数据rebalance已经完成？

目前只能通过观察graph内部的计数器，来判断整个数据迁移工作是否完成；观察方法如下：对所有新扩容的graph实例，访问其统计接口http://127.0.0.1:6071/counter/migrate 观察到所有的计数器都不再变化，那么就意味着迁移工作完成啦。
