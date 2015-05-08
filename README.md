falcon-hbs
==========

Heartbeat Server. 所有Agent都会连到hbs，每分钟发一次心跳，汇报自己的hostname、ip、agent version、plugin version，hbs据此
填充host表。agent还会通过hbs拿到应该监控的端口、进程，应该执行的插件等信息。

hbs要能够处理agent的上述请求，就需要与portal的数据库打交道，这是无论如何无法避免的，那就为hbs再赋予一个功能：DB的缓存器。judge
也需要通过portal的DB拿到策略列表，在一个大点的公司，judge实例可能比较多，几十个、甚至上百个，有了hbs这个DB缓存器在这了，judge就
无需直接访问DB了，从hbs获取策略列表即可。如此一来，hbs可以每分钟从DB读取一次数据，这一分钟内所有judge的请求都可以直接读取内存。
另外，DB存的是关系型数据，需要做一些转换才能被judge使用，hbs从DB中读取到数据之后顺便把转换也做了，这样所有judge就无需再做转换了。

所以hbs的逻辑就变成了：每分钟从DB中load各种数据，处理后放到内存里，静待agent、judge的请求。

## install

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/hbs.git
cd hbs
go get ./...
./control build
./control start
```

## configuration

- database: portal的db连接地址
- maxIdle: 数据库连接池的MaxIdle配置
- listen: 监听的rpc端口，judge要通过这个端口拿到策略列表
- trustable: 可信ip列表，安全起见留空即可
- http: 监听的http地址，主要是做调试
