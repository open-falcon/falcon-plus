falcon-alarm
============

judge把报警event写入redis，alarm从redis读取event，做相应处理，可能是发报警短信、邮件，可能是callback某个http地址。
生成的短信、邮件写入queue，sender模块专门负责来发送。


## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus/modules/alarm.git
cd alarm
go get ./...
./control build
./control start
```

## Configuration

- uicToken: 留空即可
- http: 监听的http端口
- queue: 要发送的短信、邮件写入的队列，需要与sender配置一致
- redis: highQueues和lowQueues区别是是否做报警合并，默认配置是P0/P1不合并，收到之后直接发出；>=P2做报警合并
- api: 其他各个组件的地址
