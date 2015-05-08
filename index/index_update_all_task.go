package index

import (
	"encoding/json"
	"fmt"
	cron "github.com/niean/cron"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/proc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	destUrlFmt = "http://%s/index/updateAll"
)

var (
	indexUpdateAllCron     = cron.New()
	indexUpdateAllCronSpec = "0 0 0 ? * 0-5" // 每周6晚上22:00执行一次
)

// 启动 索引全量更新 定时任务
func StartIndexUpdateAllTask() {
	indexUpdateAllCron.AddFunc(indexUpdateAllCronSpec, func() {
		UpdateAllIndex()
	})
	indexUpdateAllCron.Start()
}

func UpdateAllIndex() {
	startTs := time.Now().Unix()
	updateAllIndex()
	endTs := time.Now().Unix()
	log.Printf("index, update all, startTs %s, time-consuming %d sec\n", proc.FmtUnixTs(startTs), endTs-startTs)

	// statistics
	proc.IndexUpdateAll.Incr()
	proc.IndexUpdateAll.PutOther("lastStartTs", proc.FmtUnixTs(startTs))
	proc.IndexUpdateAll.PutOther("lastTimeConsumingInSec", endTs-startTs)
}

func updateAllIndex() {
	client := http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	for _, hostNamePort := range g.Config().Index.Cluster {
		if hostNamePort == "" {
			continue
		}

		destUrl := fmt.Sprintf(destUrlFmt, hostNamePort)
		getResp, err := client.Get(destUrl)
		if err != nil {
			log.Printf(hostNamePort+", index update all error,", err)
			continue
		}
		defer getResp.Body.Close()

		body, err := ioutil.ReadAll(getResp.Body)
		if err != nil {
			log.Println(hostNamePort+", index update all error,", err)
			continue
		}

		var data Dto
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(hostNamePort+", index update all error,", err)
			continue
		}

		if data.Data != "ok" {
			log.Println(hostNamePort+", index update all error, bad result,", data.Data)
			continue
		}
	}
}

type Dto struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
