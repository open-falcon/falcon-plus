## Introduction

多IDC时，可能面对 "分区到中心的专线网络质量较差&公网ACL不通" 等问题。这时，可以在分区内部署一套数据路由服务，接收本分区内的所有流量(包括所有的agent流量)，然后通过公网(开通ACL)，将数据push给中心的Transfer。如下图，
![gateway.png](https://raw.githubusercontent.com/niean/niean.github.io/master/images/20150806/gateway.png)

站在client端的角度，gateway和transfer提供了完全一致的功能和接口。**只有遇到网络分区的情况时，才有必要使用gateway组件**。

## Installation

首先，通过github仓库的源码，编译出可执行的二进制文件。然后，将二进制文件部署到服务器上，并提供服务。

### Build

```bash
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/gateway.git

cd gateway
go get ./...
./control build
./control pack
```
最后一步会pack出一个`falcon-gateway-$vsn.tar.gz`的安装包，拿着这个包去部署服务即可。我们也提供了编译好的安装包，在[这里](https://github.com/open-falcon/gateway/releases/tag/v0.0.5)。

### Deploy
服务部署，包括配置修改、启动服务、检验服务、停止服务等。这之前，需要将安装包解压到服务的部署目录下。

```bash
# download 'falcon-gateway-$vsn.tar.gz'
# tar -zxf falcon-gateway-$vsn.tar.gz && rm -f falcon-gateway-$vsn.tar.gz

# modify config
mv cfg.example.json cfg.json
vim cfg.json

# start service
./control start

# check, you should get 'ok'
curl -s "127.0.0.1:6060/health"

...
# stop service
./control stop

```
服务启动后，可以通过日志查看服务的运行状态，日志文件地址为./var/app.log。可以通过调试脚本./test/debug查看服务器的内部状态数据，如 运行 bash ./test/debug 可以得到服务器内部状态的统计信息。

gateway组件，部署于分区中。单个gateway实例的转发能力，为 {1核, 500MB内存, Qps不小于1W/s}；但我们仍然建议，一个分区至少部署两个gateway实例，来实现高可用。


## Usage
send items via transfer's http-api

```bash
#!/bin/bash
e="test.endpoint.1" 
m="test.metric.1"
t="t0=tag0,t1=tag1,t2=tag2"
ts=`date +%s`
curl -s -X POST -d "[{\"metric\":\"$m\", \"endpoint\":\"$e\", \"timestamp\":$ts,\"step\":60, \"value\":9, \"counterType\":\"GAUGE\",\"tags\":\"$t\"}]" "127.0.0.1:6060/api/push" | python -m json.tool
```

## Configuration

    debug: true/false, 如果为true，日志中会打印debug信息

    http
        - enable: true/false, 表示是否开启该http端口，该端口为控制端口，主要用来对transfer发送控制命令、统计命令、debug命令等
        - listen: 表示监听的http端口

    rpc
        - enable: true/false, 表示是否开启该jsonrpc数据接收端口, Agent发送数据使用的就是该端口
        - listen: 表示监听的jsonrpc端口

    socket #即将被废弃,请避免使用
        - enable: true/false, 表示是否开启该telnet方式的数据接收端口，这是为了方便用户一行行的发送数据给transfer
        - listen: 表示监听的socket端口

    transfer
        - enable: true/false, 表示是否开启向tranfser转发数据
        - batch: 数据转发的批量大小，可以加快发送速度，建议保持默认值
        - connTimeout: 单位是毫秒，与后端建立连接的超时时间，可以根据网络质量微调，建议保持默认
        - callTimeout: 单位是毫秒，发送数据给后端的超时时间，可以根据网络质量微调，建议保持默认
        - maxConns: 连接池相关配置，最大连接数，建议保持默认
        - maxIdle: 连接池相关配置，最大空闲连接数，建议保持默认
        - addr: transfer服务的地址，形如 $hostname:$port

注意: `transfer.addr`应该指向一台transfer服务器，或者，指向一个挂在了transfer服务器集群的 域名或vip。这里，不提供指向transfer集群列表的功能，有此需求的用户可以进行二次开发。

## Debug
可以通过调试脚本./test/debug查看服务器的内部状态数据，含义如下

```bash
# bash ./test/debug
{
    "data": [
        { // counter of items received
            "Cnt": 0,
            "Name": "RecvCnt",
            "Other": {},
            "Qps": 0,
            "Time": "2015-08-14 06:41:57"
        },
        { // counter of items sent to transfer
            "Cnt": 0,
            "Name": "SendCnt",
            "Other": {},
            "Qps": 0,
            "Time": "2015-08-14 06:41:57"
        },
        { // counter of items drop(gateway would drop items when caching too many)
            "Cnt": 0,
            "Name": "SendDropCnt",
            "Other": {},
            "Qps": 0,
            "Time": "2015-08-14 06:41:57"
        },
        { // counter of items sent to transfer failed
            "Cnt": 0,
            "Name": "SendFailCnt",
            "Other": {},
            "Qps": 0,
            "Time": "2015-08-14 06:41:57"
        },
        { // size of cached items
            "Cnt": 0,
            "Name": "SendQueuesCnt",
            "Other": {},
            "Time": "2015-08-14 06:41:54"
        }
    ],
    "msg": "success"
}

```

## TODO
+ 加密gateway经过公网传输的数据
