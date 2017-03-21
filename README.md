# Falcon+

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/open-falcon/falcon-plus.svg?branch=plus-dev)](https://travis-ci.org/open-falcon/falcon-plus)
[![codecov](https://codecov.io/gh/open-falcon/falcon-plus/branch/plus-dev/graph/badge.svg)](https://codecov.io/gh/open-falcon/falcon-plus)
[![GoDoc](https://godoc.org/github.com/open-falcon/falcon-plus?status.svg)](https://godoc.org/github.com/open-falcon/falcon-plus)
[![Code Issues](https://www.quantifiedcode.com/api/v1/project/5035c017b02c4a4a807ebc4e9f153e6f/badge.svg)](https://www.quantifiedcode.com/app/project/5035c017b02c4a4a807ebc4e9f153e6f)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-falcon/falcon-plus)](https://goreportcard.com/report/github.com/open-falcon/falcon-plus)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# Documentations

- [Installation and Usage](http://book.open-falcon.org)
- [Open-Faclon API](http://api.open-falcon.org)

# Get Started

*before start , please make sure you prepared [this](https://book.open-falcon.org/zh/install_from_src/prepare.html) ready.*

And then

```
git clone https://github.com/open-falcon/falcon-plus.git
cd falcon-plus/
```

# Compilation

```
# all modules
make all

# specified module
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

# Run Open-Falcon Commands

for example:

```
# ./open-falcon [start|stop|restart|check|monitor|reload] module
./open-falcon start hbs

./open-falcon check
        falcon-graph         UP           79292 
          falcon-hbs         UP           79295 
       falcon-sender         UP           79298 
        falcon-judge         UP           79301 
     falcon-transfer         UP           79304 
       falcon-nodata         UP           79307 
         falcon-task         UP           79310 
   falcon-aggregator         UP           79313 
        falcon-agent         UP           79316 
      falcon-gateway         UP           79319 
          falcon-api         UP           98421 
        falcon-alarm         UP           86388

```

* For debugging , You can check `$WorkDir/$moduleName/log/logs/xxx.log`

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

Any issue or question is welcome, Please feel free to open github issues:)
