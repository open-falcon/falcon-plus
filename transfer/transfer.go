package transfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	cron "github.com/niean/cron"
	"github.com/open-falcon/model"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/proc"
	TSemaphore "github.com/toolkits/concurrent/semaphore"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	srcUrlFmt = "http://%s:6060/statistics/all"
	destUrl   = "http://127.0.0.1:1988/v1/push"
)

var (
	transferMonitorCronSpec = "0 * * * * ?"
	transferMonitorSema     = TSemaphore.NewSemaphore(1)
	transferMonitorCron     = cron.New()
)

func Start() {
	if g.Config().Transfer.Enabled {
		go StartTransferMonitorCron()
		log.Println("transfer.Start, ok")
	} else {
		log.Println("transfer.Start, not enable")
	}

}

func StartTransferMonitorCron() {
	transferMonitorCron.AddFunc(transferMonitorCronSpec, func() {
		MonitorTransfer()
	})
	transferMonitorCron.Start()
}

func MonitorTransfer() {
	if transferMonitorSema.AvailablePermits() <= 0 {
		log.Println("transfer.monitorTransfer, concurrent not available")
		return
	}
	transferMonitorSema.Acquire()
	defer transferMonitorSema.Release()

	startTs := time.Now().Unix()
	monitorTransfer()
	endTs := time.Now().Unix()
	log.Printf("monitorTransfer, startTs %s, time-consuming %d sec\n", proc.FmtUnixTs(startTs), endTs-startTs)

	// statistics
	proc.TransferMonitorCronCnt.Incr()
	proc.TransferMonitorCronCnt.PutOther("lastStartTs", proc.FmtUnixTs(startTs))
	proc.TransferMonitorCronCnt.PutOther("lastTimeConsumingInSec", endTs-startTs)
}

func monitorTransfer() {
	tags := "type=statistics,pdl=falcon,module=transfer,owner=niean"
	client := http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	for _, tran := range g.Config().Transfer.Cluster {
		ts := time.Now().Unix()
		jsonList := make([]model.JsonMetaData, 0)

		// get statistics by http-get
		srcUrl := fmt.Sprintf(srcUrlFmt, tran)
		getResp, err := client.Get(srcUrl)
		if err != nil {
			log.Println("transfer, get statistics error,", err)
			continue
		}
		defer getResp.Body.Close()

		body, err := ioutil.ReadAll(getResp.Body)
		if err != nil {
			log.Println("transfer, get statistics error,", err)
			continue
		}

		var data Dto
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("transfer, get statistics error,", err)
			continue
		}

		for _, item := range data.Data {
			if item["Name"] == nil {
				continue
			}
			itemName := item["Name"].(string)

			var jmdCnt model.JsonMetaData
			jmdCnt.Endpoint = tran
			jmdCnt.Metric = itemName
			jmdCnt.Timestamp = ts
			jmdCnt.Step = 60
			jmdCnt.Value = int64(item["Cnt"].(float64))
			jmdCnt.CounterType = "GAUGE"
			jmdCnt.Tags = tags
			jsonList = append(jsonList, jmdCnt)

			if item["Qps"] == nil {
				continue
			}
			var jmdQps model.JsonMetaData
			jmdQps.Endpoint = tran
			jmdQps.Metric = itemName + ".Qps"
			jmdQps.Timestamp = ts
			jmdQps.Step = 60
			jmdQps.Value = int64(item["Qps"].(float64))
			jmdQps.CounterType = "GAUGE"
			jmdQps.Tags = tags
			jsonList = append(jsonList, jmdQps)
		}

		if len(jsonList) < 1 { //没取到数据
			log.Println("get null from ", tran)
			continue
		}

		// format result
		jsonBody, err := json.Marshal(jsonList)
		if err != nil {
			log.Println("transfer, format body error,", tran, err)
			continue
		}

		// send by http-post
		req, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		postResp, err := client.Do(req)
		if err != nil {
			log.Println("transfer, post to dest error,", tran, err)
			continue
		}
		defer postResp.Body.Close()

		if postResp.StatusCode/100 != 2 {
			log.Println("transfer, post to dest, bad response,", tran, postResp.StatusCode)
		}
	}
}

type Dto struct {
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}
