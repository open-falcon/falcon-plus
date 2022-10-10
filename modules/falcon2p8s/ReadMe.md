# falcon2p8s

## 该组件是一个prometheus exporter，实现了将falcon监控数据转储到prometheus


## 源码编译

```bash
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus.git
cd falcon-plus/modules/falcon2p8s
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
    "log_level": "debug",
    "concurrent": 100, # 处理metric的最大并发数
	"http": {
		"listen": "0.0.0.0:9090"
	},
    "rpc": {
        "listen": "0.0.0.0:8080"
    }
}
```

sm.example.yaml是使用prometheus operator中servicemonitor做服务发现时，对应的yaml，作为参考