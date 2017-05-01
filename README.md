# Falcon+

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/open-falcon/falcon-plus.svg?branch=plus-dev)](https://travis-ci.org/open-falcon/falcon-plus)
[![codecov](https://codecov.io/gh/open-falcon/falcon-plus/branch/plus-dev/graph/badge.svg)](https://codecov.io/gh/open-falcon/falcon-plus)
[![GoDoc](https://godoc.org/github.com/open-falcon/falcon-plus?status.svg)](https://godoc.org/github.com/open-falcon/falcon-plus)
[![Code Issues](https://www.quantifiedcode.com/api/v1/project/5035c017b02c4a4a807ebc4e9f153e6f/badge.svg)](https://www.quantifiedcode.com/app/project/5035c017b02c4a4a807ebc4e9f153e6f)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-falcon/falcon-plus)](https://goreportcard.com/report/github.com/open-falcon/falcon-plus)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# Documentations

- [Usage](http://book.open-falcon.org)
- [Open-Faclon API](http://open-falcon.org/falcon-plus)

# Getting Started

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
mysql -h 127.0.0.1 -u root -p < uic-db-schema.sql
mysql -h 127.0.0.1 -u root -p < portal-db-schema.sql
mysql -h 127.0.0.1 -u root -p < graph-db-schema.sql
mysql -h 127.0.0.1 -u root -p < dashboard-db-schema.sql
mysql -h 127.0.0.1 -u root -p < alarms-db-schema.sql
```

**NOTE: if you are upgrading from v0.1 to current version v0.2.0,then**

    mysql -h 127.0.0.1 -u root -p < alarms-db-schema.sql

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
* if you got dependencies fetch error on `golang.org/x/*`, please run `sh fix_gosrouce_fetch_error.sh` before you compile it. 

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

* For debugging , You can check `$WorkDir/$moduleName/log/logs/xxx.log`

# Install Front-end Dashboard
- Follow [this](https://github.com/open-falcon/dashboard).

**NOTE: if you want to use grafana as the dashboard, please check [this](https://github.com/open-falcon/grafana-openfalcon-datasource).**

# Package Management
## How-to

Make sure you're using Go 1.5+ and **GO15VENDOREXPERIMENT=1** env var is exported. (You can ignore GO15VENDOREXPERIMENT using Go 1.6+.)

 0. Install `trash` by `go get github.com/rancher/trash`.
 1. Edit `trash.yml` file to your needs. See the example as follow.
 2. Run `trash --keep` to download the dependencies.

trash file example:

```
package: github.com/open-falcon/falcon-plus

import:
- package: github.com/open-falcon/common              # package name
  version: origin/develop                        # tag, commit, or branch
  repo:    https://github.com/open-falcon/common.git  # (optional) git URL
```

# Package Release

```
make clean all pack
```

# Q&A

Any issue or question is welcome, Please feel free to open [github issues](https://github.com/open-falcon/falcon-plus/issues) :)
