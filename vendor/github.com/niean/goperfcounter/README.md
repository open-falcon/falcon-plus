GoPerfcounter
==========

goperfcounter用于golang应用的业务监控。goperfcounter需要和开源监控系统[Open-Falcon](http://book.open-falcon.com/zh/index.html)一起使用。

概述
-----
使用goperfcounter进行golang应用的监控，大体如下: 

1. 用户在其golang应用代码中，调用goperfcounter提供的统计函数；统计函数被调用时，perfcounter会生成统计记录、并保存在内存中
2. goperfcounter会自动的、定期的将这些统计记录push给Open-Falcon的收集器([agent](https://github.com/open-falcon/agent)或[transfer](https://github.com/open-falcon/transfer))
3. 用户在Open-Falcon中，查看统计数据的绘图曲线、设置实时报警

另外，goperfcounter提供了golang应用的基础监控，包括runtime指标、debug指标等。默认情况下，基础监控是关闭的，用户可以通过[配置文件](#配置)来开启此功能。

安装
-----

在golang项目中使用goperfcounter时，需要进行安装，操作如下

```bash
go get github.com/niean/goperfcounter

```

使用
-----

用户需要引入goperfcounter包，需要在代码片段中调用goperfcounter的[API](#API)。比如，用户想要统计函数的出错次数，可以调用`Meter`方法。

```go
package xxx

import (
	pfc "github.com/niean/goperfcounter"
)

func foo() {
	if err := bar(); err != nil {
		pfc.Meter("bar.called.error", int64(1))
	}
}

func bar() error {
	// do sth ...
	return nil
}

```

这个调用主要会产生2个Open-Falcon统计指标，如下。其中，`timestamp `和`value`是监控数据的取值；`endpoint`默认为服务器`Hostname()`，可以通过配置文件设置；`step`默认为60s，可以通过配置文件设置；`tags`中包含一个`name=bar.called.error`的标签(`bar.called.error`为用户自定义的统计器名称)，其他`tags`标签可以通过配置文件设置；`counterType `和`metric`由goperfcounter决定。

```python
{
    "counterType": "GAUGE",
    "endpoint": "git",
    "metric": "rate",
    "step": 20,
    "tags": "module=perfcounter,name=bar.called.error",
    "timestamp": 1451397266,
    "value": 13.14
},
{
    "counterType": "GAUGE",
    "endpoint": "git",
    "metric": "sum",
    "step": 20,
    "tags": "module=perfcounter,name=bar.called.error",
    "timestamp": 1451397266,
    "value": 1023
}

```


配置
----
默认情况下，goperfcounter不需要进行配置。如果用户需要定制goperfcounter的行为，可以通过配置文件来进行。配置文件需要满足以下的条件:

+ 配置文件必须和golang二进制文件应用文件，在同一目录
+ 配置文件命名，必须为```perfcounter.json```

配置文件的内容，如下

```go
{
    "debug": false, // 是否开启调制，默认为false
    "hostname": "", // 机器名(也即endpoint名称)，默认为本机名称
    "tags": "", // tags标签，默认为空。一个tag形如"key=val"，多个tag用逗号分隔；name为保留字段，因此不允许设置形如"name=xxx"的tag。eg. "cop=xiaomi,module=perfcounter"
    "step": 60, // 上报周期，单位s，默认为60s
    "bases":[], // gvm基础信息采集，可选值为"debug"、"runtime"，默认不采集
    "push": { // push数据到Open-Falcon
        "enabled":true, // 是否开启自动push，默认开启
        "api": "" // Open-Falcon接收器地址，默认为本地agent，即"http:// 127.0.0.1:1988/v1/push"
    },
    "http": { // http服务，为了安全考虑，当前只允许本地访问
        "enabled": false, // 是否开启http服务，默认不开启
        "listen": "" // http服务监听地址，默认为空。eg. "0.0.0.0:2015"表示在2015端口开启http监听
    }
}

```



API
----

几个常用接口，如下。

|接口名称|例子|使用场景|
|:----|:----|:---|
|Meter|`// 统计页面访问次数，每来一次请求，pv加1`<br/>`Meter("pageView", int64(1)) `|Meter用于累加计数。输出累加求和、变化率|
|Gauge|`// 统计队列长度` <br/>`Gauge("queueSize", int64(len(myQueueList))) ` <br/> `GaugeFloat64("queueSize", float64(len(myQueueList)))`|Gauge用于记录瞬时值。支持int64、float64类型|
|Histogram|`// 统计线程并发度` <br/>`Histogram("processNum", int64(326)) `| Histogram用于计算统计分布。输出最大值、最小值、平均值、75th、95th、99th等|

更详细的API介绍，请移步到[这里](https://github.com/niean/goperfcounter/blob/master/doc/API.md)。



数据上报
----

goperfcounter会将各种统计器的统计结果，定时发送到Open-Falcon。每种统计器，会被转换成不同的Open-Falcon指标项，转换关系如下。每条数据，至少包含一个```name=XXX```的tag，```XXX```是用户定义的统计器名称。

<table>
<tr>
  <th>统计器类型</th>
  <th>输出指标的名称</th>
  <th>输出指标的含义</th>
</tr>
<tr>
  <th rowspan="1">Gauge</th>
  <td>value</td>
  <td>最后一次的记录值(float64)</td>
</tr>
<tr>
  <th rowspan="2">Meter</th>
  <td>sum</td>
  <td>事件发生的总次数(即所有计数的累加和)</td>
</tr>
<tr>
  <td>rate</td>
  <td>一个Open-Falcon上报周期(默认60s)内，事件发生的频率，单位CPS</td>
</tr>
<tr>
  <th rowspan="6">Histogram</th>
  <td>max</td>
  <td>采样数据的最大值</td>
</tr>
<tr>
  <td>min</td>
  <td>采样数据的最小值</td>
</tr>
<tr>
  <td>mean</td>
  <td>采样数据的平均值</td>
</tr>
<tr>
  <td>75th</td>
  <td>所有采样数据中，处于75%处的数值</td>
</tr>
<tr>
  <td>95th</td>
  <td>所有采样数据中，处于95%处的数值</td>
</tr>
<tr>
  <td>99th</td>
  <td>所有采样数据中，处于99%处的数值</td>
</tr>
</table>


Bench
----

请移步到[这里](https://github.com/niean/goperfcounter/blob/master/doc/BENCH.md)


TODO
----

+ 支持本地缓存统计数据及UI展示
