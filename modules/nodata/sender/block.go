package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	cmodel "github.com/open-falcon/common/model"
	tcron "github.com/toolkits/cron"
	thttpclient "github.com/toolkits/http/httpclient"
	ttime "github.com/toolkits/time"

	"github.com/open-falcon/nodata/g"
)

var (
	lock              = sync.RWMutex{}
	avg       float64 = -100
	dev       float64 = 0
	gaussCron         = tcron.New()
)

func startGaussCron() {
	if !g.Config().Sender.Block.EnableGauss {
		log.Println("sender.StartGaussCron warning, not enabled")
		return
	}

	// start gauss cron
	gaussCron.AddFuncCC("40 */20 * * * ?", func() {
		start := time.Now().Unix()
		cnt := calcGaussOnce()
		end := time.Now().Unix()
		if g.Config().Debug {
			log.Printf("gause cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))
		}
	}, 1)
	gaussCron.Start()
	log.Println("sender.StartGaussCron ok")
}

func getThreshold() int32 {
	lock.RLock()
	defer lock.RUnlock()

	cfg := g.Config().Sender.Block
	if avg < 0 { // gauss not inited
		return cfg.Threshold
	}

	// 3-sigma
	dt_gauss_3sigma_min := cfg.Gauss3SigmaMin
	dt_gauss_3sigma_max := cfg.Gauss3SigmaMax
	threeSigma := 3 * dev
	if threeSigma < dt_gauss_3sigma_min {
		threeSigma = dt_gauss_3sigma_min
	} else if threeSigma > dt_gauss_3sigma_max {
		threeSigma = dt_gauss_3sigma_max
	}

	// threshold
	gaussThreshold := int32(math.Ceil(avg + threeSigma))
	if gaussThreshold < cfg.Threshold {
		gaussThreshold = cfg.Threshold
	}

	return gaussThreshold
}

func CalcGaussOnce() int {
	return calcGaussOnce()
}

func calcGaussOnce() int {
	values := fetchRawItems()
	size := len(values)
	if size < 100 {
		return size
	}

	// gauss
	myavg, mydev := gaussDistribution(values)

	// filter
	nvals := make([]float64, 0)
	filter := mydev
	if filter < g.Config().Sender.Block.GaussFilter { //防止过度的过滤
		filter = g.Config().Sender.Block.GaussFilter
	}
	for _, val := range values {
		if (val-myavg) > filter || val-myavg < (-filter) {
			continue
		}
		nvals = append(nvals, val)
	}
	if len(nvals) < 100 {
		return size
	}

	// gauss
	myavg, mydev = gaussDistribution(nvals)

	lock.Lock()
	defer lock.Unlock()
	avg = myavg
	dev = mydev
	log.Printf("gause status, avg %f, dev %f\n", avg, dev)

	return size
}

func gaussDistribution(values []float64) (avg float64, dev float64) {
	size := len(values)
	if size < 1 {
		return
	}

	// avg
	myavg := float64(0.0)
	for _, val := range values {
		myavg += val
	}
	myavg /= float64(size)

	// dev
	mydev := float64(0.0)
	for _, val := range values {
		mydev += (val - myavg) * (val - myavg)
	}
	mydev = math.Sqrt(mydev / float64(size))

	return myavg, mydev
}

type GraphHistoryParam struct {
	Start            int64                   `json:"start"`
	End              int64                   `json:"end"`
	CF               string                  `json:"cf"`
	EndpointCounters []cmodel.GraphInfoParam `json:"endpoint_counters"`
}

func fetchRawItems() (values []float64) {
	cfg := g.Config()
	queryUlr := fmt.Sprintf("http://%s/graph/history", cfg.Query.QueryAddr)
	hcli := thttpclient.GetHttpClient("nodata.gauss",
		time.Millisecond*time.Duration(cfg.Query.ConnectTimeout),
		time.Millisecond*time.Duration(cfg.Query.RequestTimeout))

	// form request args
	nowTs := time.Now().Unix()
	endTs := nowTs - nowTs%1200
	startTs := endTs - 24*3600*5 //用5天的数据,做高斯拟合
	hostname, _ := os.Hostname()
	if len(cfg.Sender.Block.Hostname) > 0 {
		hostname = cfg.Sender.Block.Hostname
	}
	fcounter := "FloodRate/module=nodata,pdl=falcon,port=6090,type=statistics"
	if len(cfg.Sender.Block.FloodCounter) > 0 {
		fcounter = cfg.Sender.Block.FloodCounter
	}
	endpointCounters := make([]cmodel.GraphInfoParam, 0)
	endpointCounters = append(endpointCounters, cmodel.GraphInfoParam{Endpoint: hostname, Counter: fcounter})
	args := GraphHistoryParam{Start: startTs, End: endTs, CF: "AVERAGE", EndpointCounters: endpointCounters}

	argsBody, err := json.Marshal(args)
	if err != nil {
		log.Println(queryUlr+", format body error,", err)
		return
	}

	// fetch items
	req, err := http.NewRequest("POST", queryUlr, bytes.NewBuffer(argsBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "close")
	postResp, err := hcli.Do(req)
	if err != nil {
		log.Println(queryUlr+", post to dest error,", err)
		return
	}
	defer postResp.Body.Close()

	if postResp.StatusCode/100 != 2 {
		log.Println(queryUlr+", post to dest, bad response,", postResp.Body)
		return
	}

	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		log.Println(queryUlr+", read response error,", err)
		return
	}

	resp := make([]*cmodel.GraphQueryResponse, 0)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Println(queryUlr+", unmarshal error,", err)
		return
	}
	if len(resp) != 1 || resp[0] == nil {
		return
	}

	// store items
	values = make([]float64, 0)
	for _, glr := range resp[0].Values {
		if glr == nil || math.IsNaN(float64(glr.Value)) {
			continue
		}
		values = append(values, float64(glr.Value))
	}

	return values
}
