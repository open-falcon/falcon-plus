falcon-alarm
============

judge把报警event写入redis，alarm从redis读取event，做相应处理，可能是发报警短信、邮件，可能是callback某个http地址。


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
- redis: highQueues和lowQueues区别是是否做报警合并，默认配置是P0/P1不合并，收到之后直接发出；>=P2做报警合并
- api: 其他各个组件的地址, 注意plus_api_token要和falcon-plus api组件配置文件中的default_token一致 
其中chat的配置 可参考 https://github.com/yanjunhui/chat