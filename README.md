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
最后一步会pack出一个`falcon-gateway-$vsn.tar.gz`的安装包，拿着这个包去部署服务即可。我们也提供了编译好的安装包，在[这里](https://github.com/nieanan/gateway/releases/tag/v0.0.7)。

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
**注意: 从v0.0.4版以后，配置文件格式发生了变更。**主要变更项，为

1. 开关控制符更名为enabled，原来为enable
2. transfer地址配置改为集群形式cluster，原来为单个地址addr
3. transfer添加重试次数retry，默认1、不重试


```python
{
    "debug": true,
    "http": {
        "enabled": true,
        "listen": "0.0.0.0:6060" //http服务的监听端口
    },
    "rpc": {
        "enabled": true,
        "listen": "0.0.0.0:8433" //go-rpc服务的监听端口
    },
    "socket": { //即将被废弃,请避免使用
        "enabled": true,
        "listen": "0.0.0.0:4444", //telnet服务的监听端口
        "timeout": 3600
    },
    "transfer": {
        "enabled": true, //true/false, 表示是否开启向tranfser转发数据
        "batch": 200, //数据转发的批量大小，可以加快发送速度，建议保持默认值
        "retry": 2, //重试次数，默认1、不重试
        "connTimeout": 1000, //毫秒，与后端建立连接的超时时间，可以根据网络质量微调，建议保持默认
        "callTimeout": 5000, //毫秒，发送数据给后端的超时时间，可以根据网络质量微调，建议保持默认
        "maxConns": 32, //连接池相关配置，最大连接数，建议保持默认
        "maxIdle": 32, //连接池相关配置，最大空闲连接数，建议保持默认
        "cluster": { //transfer服务器集群，支持多条记录
            "t1": "127.0.0.1:8433" //一个transfer实例，形如"node":"$hostname:$port"
        }
    }
}
```

从版本**v0.0.11**后，gateway组件引入了golang业务监控组件[GoPerfcounter](https://github.com/niean/goperfcounter)。GoPerfcounter会主动将gateway的内部状态数据，push给本地的falcon-agent，其配置文件`perfcounter.json`内容如下，含义见[这里](https://github.com/niean/goperfcounter/blob/master/README.md#配置)

```python
{
    "tags": "service=gateway", // 业务监控数据的标签
    "bases": ["debug","runtime"], // 开启gvm基础信息采集
    "push": { // 开启主动推送,数据将被推送至本机的falcon-agent
        "enabled": true
    },
    "http": { // 开启http调试，并复用gateway的http端口
        "enabled": true
    }
}
```

## Debug
可以通过调试脚本./test/debug查看服务器的内部状态数据，含义如下

```bash
# bash ./test/debug
{
    "data": {
        "gauge": {
            "SendQueueSize": { // size of cached items
                "value": 0
            }
        },
        "meter": {
            "Recv": { // counter of items received
                "rate": 954.88407253945127,
                "rate.15min": 938.12973764674587,
                "rate.1min": 892.82060496256759,
                "rate.5min": 889.51059449035426,
                "sum": 2460636
            },
            "Send": { // counter of items sent to transfer
                "rate": 950.21411950079619,
                "rate.15min": 918.55392627259835,
                "rate.1min": 886.32981239416608,
                "rate.5min": 888.16132862191205,
                "sum": 2458708
            },
            "SendFail": { // counter of items sent to transfer failed
                "rate": 0,
                "rate.15min": 0,
                "rate.1min": 0,
                "rate.5min": 0,
                "sum": 0
            },  
            "SendDrop": { // counter of items sent to transfer drop
                "rate": 0,
                "rate.15min": 0,
                "rate.1min": 0,
                "rate.5min": 0,
                "sum": 0
            },    
        }
    },
    "msg": "success"
}
```

## TODO
+ 加密gateway经过公网传输的数据
