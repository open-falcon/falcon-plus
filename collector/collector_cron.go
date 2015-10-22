package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	tsema "github.com/toolkits/concurrent/semaphore"
	tcron "github.com/toolkits/cron"
	thttpclient "github.com/toolkits/http/httpclient"
	ttime "github.com/toolkits/time"

	"github.com/open-falcon/nodata/config"
	"github.com/open-falcon/nodata/g"
)

var (
	collectorCron = tcron.New()
)

func StartCollectorCron() {
	collectorCron.AddFuncCC("*/20 * * * * ?", func() {
		start := time.Now().Unix()
		cnt := collectDataOnce()
		end := time.Now().Unix()
		if g.Config().Debug {
			log.Printf("collect cron, cnt %d, time %ds, start %s\n", cnt, end-start, ttime.FormatTs(start))
		}

		// statistics
		g.CollectorCronCnt.Incr()
		g.CollectorLastTs.SetCnt(end - start)
		g.CollectorLastCnt.SetCnt(int64(cnt))
		g.CollectorCnt.IncrBy(int64(cnt))
	}, 1)
	collectorCron.Start()
}

func CollectDataOnce() int {
	return collectDataOnce()
}

func collectDataOnce() int {
	keys := config.Keys()
	keysLen := len(keys)

	// 并发+同步控制
	cfg := g.Config().Collector
	concurrent := int(cfg.Concurrent)
	if concurrent < 1 || concurrent > 50 {
		concurrent = 10
	}
	sema := tsema.NewSemaphore(concurrent)

	batch := int(cfg.Batch)
	if batch < 100 || batch > 1000 {
		batch = 200 //batch不能太小, 否则channel将会很大
	}

	batchCnt := (keysLen + batch - 1) / batch
	rch := make(chan int, batchCnt+1)

	i := 0
	for i < keysLen {
		leftLen := keysLen - i
		fetchSize := batch // 每次处理batch个配置
		if leftLen < fetchSize {
			fetchSize = leftLen
		}
		fetchKeys := keys[i : i+fetchSize]

		// 并发collect数据
		sema.Acquire()
		go func(keys []string, keySize int) {
			defer sema.Release()
			size, _ := fetchItemsAndStore(keys, keySize)
			rch <- size
		}(fetchKeys, fetchSize)

		i += fetchSize
	}

	collectCnt := 0
	for i := 0; i < batchCnt; i++ {
		select {
		case cnt := <-rch:
			collectCnt += cnt
		}
	}

	return collectCnt
}

func fetchItemsAndStore(fetchKeys []string, fetchSize int) (size int, errt error) {
	if fetchSize < 1 {
		return
	}

	cfg := g.Config()
	queryUlr := fmt.Sprintf("http://%s/graph/last", cfg.Query.QueryAddr)
	hcli := thttpclient.GetHttpClient("nodata.collector",
		time.Millisecond*time.Duration(cfg.Query.ConnectTimeout),
		time.Millisecond*time.Duration(cfg.Query.RequestTimeout))

	// form request args
	args := make([]*cmodel.GraphLastParam, 0)
	for _, key := range fetchKeys {
		ndcfg, found := config.GetNdConfig(key)
		if !found {
			continue
		}

		endpoint := ndcfg.Endpoint
		counter := cutils.Counter(ndcfg.Metric, ndcfg.Tags)
		arg := &cmodel.GraphLastParam{endpoint, counter}
		args = append(args, arg)
	}
	if len(args) < 1 {
		return
	}
	argsBody, err := json.Marshal(args)
	if err != nil {
		log.Println(queryUlr+", format body error,", err)
		errt = err
		return
	}

	// fetch items
	req, err := http.NewRequest("POST", queryUlr, bytes.NewBuffer(argsBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "close")
	postResp, err := hcli.Do(req)
	if err != nil {
		log.Println(queryUlr+", post to dest error,", err)
		errt = err
		return
	}
	defer postResp.Body.Close()

	if postResp.StatusCode/100 != 2 {
		log.Println(queryUlr+", post to dest, bad response,", postResp.Body)
		errt = fmt.Errorf("request failed, %s", postResp.Body)
		return
	}

	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		log.Println(queryUlr+", read response error,", err)
		errt = err
		return
	}

	resp := []*cmodel.GraphLastResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Println(queryUlr+", unmarshal error,", err)
		errt = err
		return
	}

	// store items
	fts := time.Now().Unix()
	for _, glr := range resp {
		//log.Printf("collect:%v\n", glr)
		if glr == nil || glr.Value == nil {
			continue
		}
		AddItem(cutils.PK2(glr.Endpoint, glr.Counter), NewDataItem(glr.Value.Timestamp, float64(glr.Value.Value), "OK", fts))
	}

	return len(resp), nil
}
