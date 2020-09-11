## 分布式部署 open-falcon 到 k8s 集群

### 构建

克隆源码：

```shell
git clone https://github.com/open-falcon/falcon-plus.git && cd falcon-plus 
```

编译源码并构建 Docker 镜像

```shell
./docker/k8s-cluster/build.sh
```

### 部署

注意点：

1. graph 为有状态组件，部署方式为每个节点一个 Deployment，同时 replica 设置为 1
2. 所有配置文件使用挂载目录的方式持久化，可以改为 ConfigMap 管理配置
3. 镜像可以使用构建脚本自己构建后推送到自己的私仓使用（推荐），或者直接使用我构建好的，仓库地址在 docker/k8s-cluster/build.sh 查看或修改