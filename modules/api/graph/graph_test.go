// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	err = config.InitDB(viper.GetBool("db.db_bug"), viper.GetViper())
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
