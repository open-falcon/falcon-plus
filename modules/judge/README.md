falcon-judge
============

Judge是用于判断是否触发报警条件的组件。

Transfer的数据不但要打到Graph来存储并绘图，还要打到Judge用于报警判断。Judge先从hbs获取所有策略列表，静等Transfer的数据转发。
每收到一条Transfer转发过来的数据，立即找到这条数据关联的Strategy、Expression，然后做阈值判断。

**如何找到关联的Strategy**
push上来的数据带有一个endpoint，endpoint通常都是hostname，hostname隶属于多个HostGroup，HostGroup可以关联多个Template，各个
Teamplate下面就是Strategy，层层顺藤摸瓜可得。但是，如果endpoint不是hostname，并没有被HostGroup管理，那就找不到了。

**如何找到关联的Expression**
这是一种更通用的方案，主要针对endpoint不是hostname的情况。push上来的数据通常带有多个tag，比如project=falcon,module=judge，
假如我们要针对所有打了project=falcon这个tag的数据做qps的阈值判断，那我们可以配置一个这样的表达式：

```
each(metric=qps project=falcon)
```

如上配置之后，push上来的数据如果发现metric=qps，并且带有project=falcon这个tag，那就说明与这个expression相关，要做相关阈值判断

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/judge.git
cd judge
go get ./...
./control build
./control start
```

## Configuration

配置文件中主要是一些连接地址和监听的端口，没啥好说的，看一下alarm的配置，judge报警判断完毕之后会产生报警event，这些event会写入
alarm的redis队列中，不同优先级（配置策略的时候每个策略会配置一个优先级，0-5）写入不同队列，alarm中除了redis地址需要修改，其他
的建议维持默认。

alarm中有一个minInterval的配置，单位是秒，默认是300秒，表示同一个event，如果配置报警多次，那么两个报警之间至少间隔300秒。
这是个经验值，我们觉得报警太频繁没有意义，对工程师来说是干扰。收到报警之后拿出电脑、开机、连上vpn就差不多要3分钟了……

