package graph

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/api/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func init() {
	viper.AddConfigPath("../")
	viper.SetConfigName("cfg.example")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = config.InitLog("debug")
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitDB(viper.GetBool("db.db_bug"))
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}

	Start(viper.GetStringMapString("graphs.cluster"))
}

func TestGraphAPI(t *testing.T) {
	Convey("testing delete item from index cache", t, func() {
		p := &cmodel.GraphCacheParam{
			Endpoint: "0.0.0.0",
			Metric:   "CollectorCronCnt.Qps",
			Step:     60,
			DsType:   "GAUGE",
			Tags:     "module=task,pdl=falcon,port=8002,type=statistics",
		}
		params := []*cmodel.GraphCacheParam{p}
		DeleteIndexCache(params)
	})
}
