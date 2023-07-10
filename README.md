# Falcon+

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/open-falcon/falcon-plus.svg?branch=plus-dev)](https://travis-ci.org/open-falcon/falcon-plus)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Backers on Open Collective](https://opencollective.com/falcon-plus/backers/badge.svg)](#backers)
[![Sponsors on Open Collective](https://opencollective.com/falcon-plus/sponsors/badge.svg)](#sponsors)

# Notice

If you are heavily using k8s and adopting microservices architecture, we recommend you to upgrade to use [Nightingale](https://github.com/ccfos/nightingale) building your modern monitoring system.

- [Nightingale](https://github.com/ccfos/nightingale) is an enterprise-level cloud-native monitoring tool, which can be used as drop-in replacement of Prometheus for alerting and management.
- [Categraf](https://github.com/flashcatcloud/categraf) is one-stop telemetry collector for Nightingale / Prometheus / M3DB / VictoriaMetrics / Thanos / Influxdb / TDengine.

It is recommended that you use [FlashDuty](https://flashcat.cloud/product/flashduty?from=categraf) as the OnCall system to realize alarm aggregation convergence, claiming, upgrading, scheduling, and coordination, so that the alarm can be reached efficiently and ensure that the alarm processing is not missed.


# Documentations

- [Usage](http://book.open-falcon.org)
- [Open-Falcon API](http://open-falcon.org/falcon-plus)

# Prerequisite

- Git >= 1.7.5
- Go >= 1.6

# Getting Started

## Docker

Please refer to ./docker/[README.md](https://github.com/open-falcon/falcon-plus/blob/master/docker/README.md).

## Build from source
**before start, please make sure you prepared this:**

```
yum install -y redis
yum install -y mysql-server

```

*NOTE: be sure to check redis and mysql-server have successfully started.*

And then

```
# Please make sure that you have set `$GOPATH` and `$GOROOT` correctly.
# If you have not golang in your host, please follow [https://golang.org/doc/install] to install golang.

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus.git

```

**And do not forget to init the database first (if you have not loaded the database schema before)**

```
cd $GOPATH/src/github.com/open-falcon/falcon-plus/scripts/mysql/db_schema/
mysql -h 127.0.0.1 -u root -p < 1_uic-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 2_portal-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 3_dashboard-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 4_graph-db-schema.sql
mysql -h 127.0.0.1 -u root -p < 5_alarms-db-schema.sql
```

**NOTE: if you are upgrading from v0.1 to v0.2.0(or above),then**. [More upgrading instruction](http://www.jianshu.com/p/6fb2c2b4d030)

    mysql -h 127.0.0.1 -u root -p < 5_alarms-db-schema.sql

# Compilation

```
cd $GOPATH/src/github.com/open-falcon/falcon-plus/

# make all modules
make all

# make specified module
make agent

# pack all modules
make pack
```

* *after `make pack` you will got `open-falcon-vx.x.x.tar.gz`*
* *if you want to edit configure file for each module, you can edit `config/xxx.json` before you do `make pack`*

#  Unpack and Decompose

```
export WorkDir="$HOME/open-falcon"
mkdir -p $WorkDir
tar -xzvf open-falcon-vx.x.x.tar.gz -C $WorkDir
cd $WorkDir
```

# Start all modules in single host
```
cd $WorkDir
./open-falcon start

# check modules status
./open-falcon check

```

# Run More Open-Falcon Commands

for example:

```
# ./open-falcon [start|stop|restart|check|monitor|reload] module
./open-falcon start agent

./open-falcon check
        falcon-graph         UP           53007
          falcon-hbs         UP           53014
        falcon-judge         UP           53020
     falcon-transfer         UP           53026
       falcon-nodata         UP           53032
   falcon-aggregator         UP           53038
        falcon-agent         UP           53044
      falcon-gateway         UP           53050
          falcon-api         UP           53056
        falcon-alarm         UP           53063
```

* For debugging , You can check `$WorkDir/$moduleName/logs/xxx.log`

# Install Frontend Dashboard
- Follow [this](https://github.com/open-falcon/dashboard).

**NOTE: if you want to use grafana as the dashboard, please check [this](https://github.com/open-falcon/grafana-openfalcon-datasource).**

# Package Release

```
make clean all pack
```

# API Standard
- [API Standard](https://github.com/open-falcon/falcon-plus/blob/master/api-standard.md)


# Q&A

- Any issue or question is welcome, Please feel free to open [github issues](https://github.com/open-falcon/falcon-plus/issues) :)
- [FAQ](http://book.open-falcon.org/zh_0_2/faq/)


## Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].
<a href="https://github.com/open-falcon/falcon-plus/contributors"><img src="https://opencollective.com/falcon-plus/contributors.svg?width=890&button=false" /></a>
