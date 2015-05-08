falcon-agent
===

This is a linux monitor agent. Just like zabbix-agent and tcollector.


## install

It is a golang classic project

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/agent.git
cd agent
go get ./...
./control build
./control start

# goto http://localhost:1988
```

I use [linux-dash](https://github.com/afaqurk/linux-dash) as the page theme.

## configuration

- heartbeat: heartbeat server rpc address
- transfer: transfer rpc address
- ignore: the metrics should ignore

# deploy

http://ulricqin.com/project/ops-updater/

