## Introduction

task是监控系统一个必要的辅助模块。有些功能，不适合与监控的核心业务耦合、无高可用要求、但又必不可少，我们把这部分功能拿出来 放到定时任务task模块中。

定时任务，要求单机部署，整套falcon系统中应该只有一个定时任务的服务实例。部署定时任务的服务器上，应该安装了falcon-agent、开放了1988的数据推送接口。

定时任务，实现了如下几个功能：

+ index更新。包括图表索引的全量更新。
+ falcon服务组件的自身状态数据采集。当前，定时任务了采集了 transfer、graph、task这三个服务的内部状态数据。
+ falcon自检控任务。

部署时，index更新、falcon自身状态采集只能部署单实例, falcon自监控任务 建议至少部署2个实例。通过cfg.json中的enable来开关某个任务。


## Installation

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/task.git

cd task
go get ./...
./control build

./control start
```

## Configuration

    debug: true/false, 如果为true，日志中会打印debug信息

    http
        - enable: true/false, 表示是否开启该http端口，该端口为控制端口，主要用来对task发送控制命令、统计命令、debug命令等
        - listen: 表示http-server监听的端口

    index
        - enable: true/false, 表示是否开启索引更新任务
        - dsn: 索引服务的MySQL的连接信息，默认用户名是root，密码为空，host为127.0.0.1，database为graph（如有必要，请修改）
        - maxIdle: MySQL连接池配置，连接池允许的最大空闲连接数，保持默认即可
        - cluster: 后端graph列表，用具体的hostname:port表示

    monitor
        - enable: true/false, 表示是否开启falcon的自监控任务
        - mailUrl: 邮件服务的http接口,用于发送自监控报警邮件
        - mainTos: 接收自监控报警邮件的邮箱地址,多个邮箱地址用逗号隔开
        - cluster: falcon后端服务列表，用具体的"module,hostname:port"表示，module取值可以为graph、transfer、judge、task等任意falcon组件
        
    collector
        - enable: true/false, 表示是否开启falcon的自身状态采集任务
        - destUrl: 监控数据的push地址,默认为本机的1988接口
        - srcUrlFmt: 监控数据采集的url格式, %s将由机器名或域名替换
        - cluster: falcon后端服务列表，用具体的"module,hostname:port"表示，module取值可以为graph、transfer、task等

### 如何清除过期索引
监控数据停止上报后，该数据对应的索引也会停止更新、变为过期索引。过期索引，影响视听，部分用户希望删除之。

我们原来的方案，是: 通过task模块，有数据上报的索引、每天被更新一次，7天未被更新的索引、清除之。但是，很多用户不能正确配置graph实例的http接口，导致正常上报的监控数据的索引 无法被更新；7天后，合法索引被task模块误删除。

为了解决上述问题，我们停掉了task模块自动删除过期索引的功能、转而提供了过期索引删除的接口。用户按需触发索引删除操作，具体步骤为:

1.运行task模块，并正确配置graph集群及其http端口，即task配置文件中index.cluster的内容。此处配置不正确，不应该进行索引删除操作，否则将导致索引数据的误删除。

2.进行一次索引数据的全量更新。方法为 ``` curl -s "$Hostname.Of.Task:$Http.Port/index/updateAll" ```。这里，"$Hostname.Of.Task:$Http.Port"是task的http接口地址。
PS:索引数据存放在graph实例上，这里，只是通过task，触发了各个graph实例的索引全量更新。更直接的办法，是，到每个graph实例上，运行```curl -s "127.0.0.1:6071/index/updateAll"```，直接触发graph实例 进行索引全量更新(这里假设graph的http监听端口为6071)。

3.待索引全量更新完成后，发起过期索引删除 ``` curl -s "$Hostname.Of.Task:$Http.Port/index/delete" ```。运行索引删除前，请务必**确保索引全量更新已完成**。典型的做法为，周六运行一次索引全量更新，周日运行一次索引删除；索引更新和删除之间，留出足够的时间。

在此，建议您: **若无必要，请勿删除索引**；若确定要删除索引，请确保删除索引之前，对所有的graph实例进行一次索引全量更新。
