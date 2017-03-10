## Introduction


## Installation

源码安装过程，如下

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus

cd modules/ctrl
go get ./...
go build -o ctrl

./ctrl -config ./etc/falcon.conf.example
```


## Configuration


## 关于扩容时数据自动迁移
