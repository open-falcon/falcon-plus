# Falcon+

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/Cepave/open-falcon-backend.svg?branch=develop)](https://travis-ci.org/Cepave/open-falcon-backend)
[![codecov](https://codecov.io/gh/Cepave/open-falcon-backend/branch/develop/graph/badge.svg)](https://codecov.io/gh/Cepave/open-falcon-backend)
[![GoDoc](https://godoc.org/github.com/Cepave/open-falcon-backend?status.svg)](https://godoc.org/github.com/Cepave/open-falcon-backend)
[![Join the chat at https://gitter.im/goappmonitor/Lobby](https://badges.gitter.im/goappmonitor/Lobby.svg)](https://gitter.im/goappmonitor/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Code Health](https://landscape.io/github/Cepave/open-falcon-backend/master/landscape.svg?style=flat)](https://landscape.io/github/Cepave/open-falcon-backend/master)
[![Code Issues](https://www.quantifiedcode.com/api/v1/project/df24b20e9c504ad0a2ac9fa3e99936f5/badge.svg)](https://www.quantifiedcode.com/app/project/df24b20e9c504ad0a2ac9fa3e99936f5)
[![Go Report Card](https://goreportcard.com/badge/github.com/Cepave/open-falcon-backend)](https://goreportcard.com/report/github.com/Cepave/open-falcon-backend)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# Documentations

- http://book.open-falcon.org
- http://docs.openfalcon.apiary.io

# Get Started

    git clone https://github.com/open-falcon/open-falcon.git
    cd open-falcon

# Compilation

```bash
# all modules
make all

# specified module
make agent
```

# Run Open-Falcon Commands

Agent for example:

    ./open-falcon agent [build|pack|start|stop|restart|status|tail]

# Package Management
## How-to

Make sure you're using Go 1.5+ and **GO15VENDOREXPERIMENT=1** env var is exported. (You can ignore GO15VENDOREXPERIMENT using Go 1.6+.)

 0. Install `trash` by `go get github.com/rancher/trash`.
 1. Edit `trash.yml` file to your needs. See the example as follow.
 2. Run `trash --keep` to download the dependencies.

```yaml
package: github.com/open-falcon/open-falcon

import:
- package: github.com/open-falcon/common              # package name
  version: origin/develop                        # tag, commit, or branch
  repo:    https://github.com/open-falcon/common.git  # (optional) git URL
```

# Package Release

	make clean all pack
