# Aggregator

集群聚合模块。聚合某集群下的所有机器的某个指标的值，提供一种集群视角的监控体验。

## 源码编译

```bash
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/sdk.git
git clone https://github.com/open-falcon/falcon-plus.git
cd falcon-plus/modules/aggregator
go get 
./control build
./control pack
```

最后一步会pack出一个tar.gz的安装包，拿着这个包去部署服务即可。

## 服务部署
服务部署，包括配置修改、启动服务、检验服务、停止服务等。这之前，需要将安装包解压到服务的部署目录下。

```bash
# 修改配置, 配置项含义见下文
mv cfg.example.json cfg.json
vim cfg.json

# 启动服务
./control start

# 校验服务，看端口是否在监听
ss -tln

# 检查log
./control tail

...
# 停止服务
./control stop

```

## 配置说明
配置文件默认为./cfg.json。默认情况下，安装包会有一个cfg.example.json的配置文件示例。各配置项的含义，如下

```bash
## Configuration
{
    "debug": true,
    "http": {
        "enabled": true,
        "listen": "0.0.0.0:6055"
    },
    "database": {
        "addr": "root:@tcp(127.0.0.1:3306)/falcon_portal?loc=Local&parseTime=true",
        "idle": 10,
        "ids": [1,-1], # aggregator模块可以部署多个实例，这个配置表示当前实例要处理的数据库中cluster表的id范围
        "interval": 55
    },
    "api": {
        "hostnames": "http://127.0.0.1:5050/api/group/%s/hosts.json", # 注意修改为你的portal的ip:port
        "push": "http://127.0.0.1:6060/api/push", # 注意修改为你的transfer的ip:port
        "graphLast": "http://127.0.0.1:9966/graph/last" # 注意修改为你的query的ip:port
    }
}
       
```
