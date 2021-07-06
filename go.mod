module github.com/open-falcon/falcon-plus

go 1.15

require (
	github.com/astaxie/beego v1.8.3
	github.com/denisenkom/go-mssqldb v0.10.0 // indirect
	github.com/emirpasic/gods v1.9.0
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.4.0
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/influxdata/influxdb v1.9.0
	github.com/jinzhu/gorm v0.0.0-20170703134954-2a1463811ee1
	github.com/jinzhu/inflection v0.0.0-20170102125226-1c35d901db3d // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f
	github.com/juju/testing v0.0.0-20210324180055-18c50b0c2098 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/masato25/resty v0.4.2-0.20161209040832-927c0e7d74a0
	github.com/masato25/yaag v0.0.0-20170704095552-00862ec4db8e
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989
	github.com/niean/go-metrics-lite v0.0.0-20151230091537-b5d30971b578 // indirect
	github.com/niean/goperfcounter v0.0.0-20160108100052-24860a8d3fac
	github.com/niean/gotools v0.0.0-20151221085310-ff3f51fc5c60 // indirect
	github.com/open-falcon/rrdlite v0.0.0-20170412122036-7d8646c85cc5
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/toolkits/cache v0.0.0-20190218093630-cfb07b7585e5
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3
	github.com/toolkits/conn_pool v0.0.0-20170512061817-2b758bec1177
	github.com/toolkits/consistent v0.0.0-20150827090850-a6f56a64d1b1
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/core v0.0.0-20141116054942-0ebf14900fe2
	github.com/toolkits/cron v0.0.0-20150624115642-bebc2953afa6
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/http v0.0.0-20150609122824-f3ac6e6c24be
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	github.com/toolkits/nux v0.0.0-20200401110743-debb3829764a
	github.com/toolkits/proc v0.0.0-20170520054645-8c734d0eb018
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484
	github.com/toolkits/str v0.0.0-20160913030958-f82e0f0498cb
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf
	github.com/toolkits/time v0.0.0-20160524122720-c274716e8d7f
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect

)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
