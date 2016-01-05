# Graph

graph组件，负责存储监控绘图数据。

## 功能描述
----
graph组件，使用rrdtool，来存储监控历史数据。主要功能，如下：

0. 提供Golang-RPC接口，接收格式化的监控数据。监控数据的格式，见[这里](http://blog.niean.name/2015/08/06/falcon-intro/#数据模型)。
1. 支持存储原始数据。默认保存12H的原始数据。
2. 支持对原始数据进行归档存储。支持5min、20min、3h、12h四个归档粒度。
3. 提供Http-Get和Golang-Rpc两种接口，用于查询历史数据。查询数据时，graph会自适应选择归档粒度，使返回的数据点不过多。


## 模块结构
----
graph主要的模块结构（含数据流）如下。

![falcon.graph.arch](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/graph/graph.arch.png)

graph的模块，主要分为两个功能：数据接收、数据查询。

#### 数据接收

graph对外提供了golang-rpc形式的接口，用于接收外部push过来的监控数据。收到的数据，被复制到索引生成器`index`一份、被复制到接收缓存`recv cache`一份。索引生成器，会为新增的监控数据建立一条索引信息（形如`endpoint/metric/tags`），并存储到外部的mysql数据库中、供dashboard等系统使用。接收缓存，是监控数据被写入磁盘`rrdfile`前的一个缓存，目的是优化磁盘写入性能。

#### 数据查询
graph对外提供了golang-rpc形式的接口，用于接收外部的查询请求。query模块，会使用索引信息来验证外部请求的counter合法性，验证通过后，读取磁盘上`rrdfile`内容、并merge接收缓存中尚未flush到磁盘的监控数据，得到完整的结果后 返回给外部。


## 部署实践
----

graph消耗的主要资源是disk，同时会消耗部分mem和cpu。可以根据监控指标个数，来预估graph集群的容量。默认的，graph资源消耗的一个参考值，如下。

|监控指标量|磁盘空间|磁盘写|磁盘读|内存|CPU|
|:----|:----|:----|:----|:----|:----|
|10K条|1GB|100KB/s|150KB/s|125MB|0.05核|

建议，部署graph组建的服务器，配置较大容量的SSD硬盘。


## 扩容
----
系统扩容时，只需求改transfer和query的graph集群配置即可，graph组件会自动完成历史数据迁移的工作。扩容期间，graph可以正常的对外提供服务。


## 缺点
----

1. 对磁盘资源消耗严重。rrdtool自带的归档功能，会消耗大量的磁盘IO。
2. 精确的历史数据保存时间短，不利于历史的现场回放。默认只保存12H的原始数据。
3. 绘图数据的高可用，实现成本较高。冷备、热备绘图数据都多多少少存在一些问题，灾难恢复也可能需要较多的时间。


