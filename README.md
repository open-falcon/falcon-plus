## intro

定时任务，是监控系统一个必要的辅助模块。有些功能，不适合与监控的核心业务耦合、无高可用要求、但又必不可少，我们把这部分功能拿出来 放到定时任务task模块中。

定时任务，要求单机部署，整套falcon系统中应该只有一个定时任务的服务实例。部署定时任务的服务器上，应该安装了falcon-agent、开放了1988的数据推送接口。

定时任务，实现了如下几个功能：

1.图表索引更新。包括图表索引的全量更新 和 垃圾索引清理。
2.falcon服务的自监控数据采集。当前，定时任务了采集了 transfer、graph、task这三个服务的内部状态数据。


## install

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

## configuration

    debug: true/false, 如果为true，日志中会打印debug信息

    http
        - enable: true/false, 表示是否开启该http端口，该端口为控制端口，主要用来对task发送控制命令、统计命令、debug命令等
        - listen: 表示http-server监听的端口

    db
        - dsn: MySQL的连接信息，默认用户名是root，密码为空，host为127.0.0.1，database为graph（如有必要，请修改）
        - maxIdle: MySQL连接池配置，连接池允许的最大空闲连接数，保持默认即可

    index
        - enable: true/false, 表示是否开启索引更新任务
        - cluster: 后端graph列表，用具体的ip:port表示

    monitor
        - enable: true/false, 表示是否开启falcon的自监控数据采集
        - cluster: falcon后端服务列表，用具体的"module,ip:port"表示，module取值可以为graph、transfer、task等

