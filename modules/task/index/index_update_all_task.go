package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cron "github.com/toolkits/cron"
	nhttpclient "github.com/toolkits/http/httpclient"
	ntime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/task/g"
	"github.com/open-falcon/falcon-plus/modules/task/proc"
)

const (
	destUrlFmt = "http://%s/index/updateAll"
)

var (
	indexUpdateAllCron = cron.New()
)

// 启动 索引全量更新 定时任务
func StartIndexUpdateAllTask() {
	for graphAddr, cronSpec := range g.Config().Index.Cluster {
		ga := graphAddr
		indexUpdateAllCron.AddFuncCC(cronSpec, func() { UpdateIndexOfOneGraph(ga, "cron") }, 1)
	}

	indexUpdateAllCron.Start()
}

// 手动触发全量更新
func UpdateAllIndex() {
	for graphAddr, _ := range g.Config().Index.Cluster {
		UpdateIndexOfOneGraph(graphAddr, "manual")
	}
}

func UpdateIndexOfOneGraph(graphAddr string, src string) {
	startTs := time.Now().Unix()
	err := updateIndexOfOneGraph(graphAddr)
	endTs := time.Now().Unix()

	// statistics
	proc.IndexUpdateCnt.Incr()
	if err == nil {
		log.Printf("index update ok, %s, %s, start %s, ts %ds",
			src, graphAddr, ntime.FormatTs(startTs), endTs-startTs)
	} else {
		proc.IndexUpdateErrorCnt.Incr()
		log.Printf("index update error, %s, %s, start %s, ts %ds, reason %v",
			src, graphAddr, ntime.FormatTs(startTs), endTs-startTs, err)
	}
}

func updateIndexOfOneGraph(hostNamePort string) error {
	if hostNamePort == "" {
		return fmt.Errorf("index update error, bad host")
	}

	client := nhttpclient.GetHttpClient("index.update."+hostNamePort, 5*time.Second, 10*time.Second)

	destUrl := fmt.Sprintf(destUrlFmt, hostNamePort)
	req, _ := http.NewRequest("GET", destUrl, nil)
	req.Header.Set("Connection", "close")
	getResp, err := client.Do(req)
	if err != nil {
		log.Printf(hostNamePort+", index update error,", err)
		return err
	}
	defer getResp.Body.Close()

	body, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		log.Println(hostNamePort+", index update error,", err)
		return err
	}

	var data Dto
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(hostNamePort+", index update error,", err)
		return err
	}

	if data.Data != "ok" {
		log.Println(hostNamePort+", index update error, bad result,", data.Data)
		return err
	}

	return nil
}

type Dto struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
