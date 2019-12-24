## Introduction

多IDC时，可能面对 "分区到中心的专线网络质量较差&公网ACL不通" 等问题。这时，可以在分区内部署一套数据路由服务，接收本分区内的所有流量(包括所有的agent流量)，然后通过公网(开通ACL)，将数据push给中心的Transfer。如下图，
![gateway.png](https://raw.githubusercontent.com/niean/niean.github.io/master/images/20150806/gateway.png)

站在client端的角度，gateway和transfer提供了完全一致的功能和接口。**只有遇到网络分区的情况时，才有必要使用gateway组件**。

### Build

请参考[falcon-plus README](https://github.com/open-falcon/falcon-plus)

### Deploy
服务部署，包括配置修改、启动服务、检验服务、停止服务等。这之前，需要将安装包解压到服务的部署目录下。

```
# modify config
mv cfg.example.json cfg.json
vim cfg.json

# start service
./open-falcon start gateway

# check, you should get 'ok'
./open-falcon monitor gateway

...
# stop service
./open-falcon stop gateway

```
服务启动后，可以通过日志查看服务的运行状态，日志文件地址为./var/app.log。可以通过调试脚本./test/debug查看服务器的内部状态数据，如 运行 bash ./test/debug 可以得到服务器内部状态的统计信息。

gateway组件，部署于分区中。单个gateway实例的转发能力，为 {1核, 500MB内存, Qps不小于1W/s}；但我们仍然建议，一个分区至少部署两个gateway实例，来实现高可用。

## Configuration


```
{
    "debug": true,
    "minStep": 30,
    "http": {
        "enabled": true,
        "listen": "0.0.0.0:16060"
    },
    "rpc": {
        "enabled": true,
        "listen": "0.0.0.0:18433"
    },
    "socket": {
        "enabled": true,
        "listen": "0.0.0.0:14444",
        "timeout": 3600
    },
    "judge": {
        "enabled": false,
        "batch": 200,
        "connTimeout": 1000,
        "callTimeout": 5000,
        "maxConns": 32,
        "maxIdle": 32,
        "replicas": 500,
        "cluster": {
            "judge-00" : "%%JUDGE_RPC%%"
        }
    },
    "graph": {
        "enabled": false,
        "batch": 200,
        "connTimeout": 1000,
        "callTimeout": 5000,
        "maxConns": 32,
        "maxIdle": 32,
        "replicas": 500,
        "cluster": {
            "graph-00" : "%%GRAPH_RPC%%"
        }
    },
   "tsdb": {
        "enabled": false,
        "batch": 200,
        "connTimeout": 1000,
        "callTimeout": 5000,
        "maxConns": 32,
        "maxIdle": 32,
        "retry": 3,
        "address": "127.0.0.1:8088"
    },
    "transfer": {
        "enabled": true,
        "batch": 200,
        "connTimeout": 1000,
        "callTimeout": 5000,
        "maxConns": 32,
        "maxIdle": 32,
        "retry": 3,
        "cluster": {
            "t1": "%%TRANSFER_RPC%%"
        }
    }
```

从版本**v0.0.11**后，gateway组件引入了golang业务监控组件[GoPerfcounter](https://github.com/niean/goperfcounter)。GoPerfcounter会主动将gateway的内部状态数据，push给本地的falcon-agent，其配置文件`perfcounter.json`内容如下，含义见[这里](https://github.com/niean/goperfcounter/blob/master/README.md#配置)

```
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

```
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
