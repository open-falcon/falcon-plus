# nodata
+ [需求定位](#需求定位)
+ [系统设计](#系统设计)
+ [用户手册](#用户手册)

nodata用于检测监控数据的上报异常。nodata和实时报警judge模块协同工作，过程为: 配置了nodata的采集项超时未上报数据，nodata生成一条默认的模拟数据；用户配置相应的报警策略，收到mock数据就产生报警。采集项上报异常检测，作为judge模块的一个必要补充，能够使judge的实时报警功能更加可靠、完善。

## 需求定位
nodata只处理如下的用户需求，

1. 监测"特征采集项"的上报异常
2. 监测少量的、十分重要的采集项的上报异常

这里的特征采集项，指的是，能够表征某一监控采集服务数据上报情况的单个采集项。例如，falcon-agent的agent.alive指标就是一个特征采集项，它能够说明agent是否正常存活，进而能够说明通过agent上报的监控数据是否正常。

nodata所谓的异常，限定为 用户数据采集服务异常、falcon数据上报链路异常等，主要场景有：

**用户数据采集服务异常**

+ 用户数据采集服务，异常终止
+ 用户数据采集服务，与falcon数据收集器之间的通信链路异常，使得数据无法上报
+ 用户数据采集服务，上报的数据格式错误

**falcon数据上报链路异常**

+ agent异常，无法接收用户的数据推送、无法主动采集监控数据
+ agent与数据转发transfer之间通信异常

出现以下情况时，nodata不应该引发大面积的报警:

+ 由于网络故障，导致大部分的采集项上报异常
+ 由于falcon自身服务故障，导致大量的采集项上报异常


> 从系统边界的描述可知，nodata只是为少数重要的采集指标而设计的。nodata处理的采集项的数量，不应该多于judge的十分之一，nodata的滥用将会给falcon的运维管理带来麻烦。

## 系统设计
#### 系统流图
![nodata.flow](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/nodata/nodata.flow.png)

#### 模块结构
![nodata.module](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/nodata/nodata.module.png)

#### 部署架构
![nodata.deploy](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/nodata/nodata.deploy.png)

## 用户手册
使用Nodata，需要进行两个配置: Nodata配置 和 策略配置。下面，我们以一个例子，讲述如何使用Nodata提供的服务。

#### 用户需求
当机器分组`cop.xiaomi_owt.inf_pdl.falcon_service.task`下的所有机器，其采集指标 `agent.alive` 上报中断时，通知用户。

#### Nodata配置
进入Nodata配置主页，点击右上角的添加按钮，添加nodata配置。
![nodata.config](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/nodata/nodata.config.open.png)

进行完上述配置后，分组`cop.xiaomi_owt.inf_pdl.falcon_service.task`下的所有机器，其采集项 `agent.alive`上报中断后，nodata服务就会补发一个取值为 `-1.0`、agent.alive的监控数据给监控系统。

#### 策略配置
配置了Nodata后，如果有数据上报中断的情况，Nodata配置中的默认值就会被上报。我们可以针对这个默认值，设置报警；只要收到了默认值，就认为发生了数据上报的中断（如果你设置的默认值，可能与正常上报的数据相等，那么请修改你的Nodata配置、使默认值有别于正常值）。将此策略，绑定到分组`cop.xiaomi_owt.inf_pdl.falcon_service.task`即可。

![nodata.judge](https://raw.githubusercontent.com/niean/niean.common.store/master/images/open-falcon/nodata/ndoata.strategy.png)

#### 注意事项
1. 配置名称name，要全局唯一。这是为了方便Nodata配置的管理。
2. 监控实例endpoint, 可以是机器分组、机器名或者其他 这三种类型，只能选择其中的一种。同一类型，支持多个记录，但建议不超过5个。选择机器分组时，系统会帮忙展开成具体机器名，支持动态生效。监控实体不是机器名时，只能选择“其他”类型。
3. 监控指标metric。
4. 数据标签tags，多个tag要用逗号隔开。必须填写完整的tags串，因为nodata会按照此tags串，去完全匹配、筛选监控数指标项。
5. 数据类型type，只支持原始值类型GAUGE。因为，nodata只应该监控 "特征指标"(如agent.alive)，"特征指标"都是GAUGE类型的。
6. 采集周期step，单位是秒。必须填写 完整&真实step。该字段不完整 或者 不真实，将会导致nodata监控的误报、漏报。
7. 补发值default，必须有别于上报的真实数据。比如，`cpu.idle`的取值范围是[0,100]，那么它的nodata默认取值 只能取小于0或者大于100的值。否则，会发生误报、漏报。
