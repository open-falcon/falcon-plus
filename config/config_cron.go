package config

import (
	"log"
	"time"

	"github.com/toolkits/container/nmap"
	tcron "github.com/toolkits/cron"
	ttime "github.com/toolkits/time"

	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/nodata/config/service"
	"github.com/open-falcon/nodata/g"
)

var (
	ndconfigCron     = tcron.New()
	ndconfigCronSpec = "50 */2 * * * ?"
)

func StartNdConfigCron() {
	ndconfigCron.AddFuncCC(ndconfigCronSpec, func() {
		start := time.Now().Unix()
		cnt, _ := syncNdConfig()
		end := time.Now().Unix()
		if g.Config().Debug {
			log.Printf("config cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))
		}

		// statistics
		g.ConfigCronCnt.Incr()
		g.ConfigLastTs.SetCnt(end - start)
		g.ConfigLastCnt.SetCnt(int64(cnt))
	}, 1)
	ndconfigCron.Start()
}

func SyncNdConfigOnce() int {
	cnt, _ := syncNdConfig()
	return cnt
}

func syncNdConfig() (cnt int, errt error) {
	// get configs
	configs := service.GetMockCfgFromDB()
	// restruct
	nm := nmap.NewSafeMap()
	for _, ndc := range configs {
		endpoint := ndc.Endpoint
		metric := ndc.Metric
		tags := ndc.Tags
		if endpoint == "" {
			log.Printf("bad config: %+v\n", ndc)
			continue
		}
		pk := cutils.PK(endpoint, metric, tags)
		nm.Put(pk, ndc)
	}

	// cache
	SetNdConfigMap(nm)

	return nm.Size(), nil
}
