package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	cron "github.com/niean/cron"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/proc"
	TSemaphore "github.com/toolkits/concurrent/semaphore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	srcUrlFmt = "http://%s/statistics/all"
	destUrl   = "http://127.0.0.1:1988/v1/push"
)

var (
	monitorCron     = cron.New()
	monitorCronSpec = "0 * * * * ?"
	monitorSema     = TSemaphore.NewSemaphore(1)
)

func Start() {
	if g.Config().Monitor.Enabled {
		go startMonitorCron()
		log.Println("monitor.Start, ok")
	} else {
		log.Println("monitor.Start, not enable")
	}

}

func startMonitorCron() {
	monitorCron.AddFunc(monitorCronSpec, func() {
		monitor()
	})
	monitorCron.Start()
}

func monitor() {
	if monitorSema.AvailablePermits() <= 0 {
		log.Println("monitor.monitor, concurrent not available")
		return
	}
	monitorSema.Acquire()
	defer monitorSema.Release()

	startTs := time.Now().Unix()
	_monitor()
	endTs := time.Now().Unix()
	log.Printf("monitor, startTs %s, time-consuming %d sec\n", proc.FmtUnixTs(startTs), endTs-startTs)

	// statistics
	proc.MonitorCronCnt.Incr()
	proc.MonitorCronCnt.PutOther("lastStartTs", proc.FmtUnixTs(startTs))
	proc.MonitorCronCnt.PutOther("lastTimeConsumingInSec", endTs-startTs)
}
func _monitor() {
	tags := "type=statistics,pdl=falcon"
	client := http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	for _, host := range g.Config().Monitor.Cluster {
		ts := time.Now().Unix()
		jsonList := make([]model.MetricValue, 0)

		// get statistics by http-get
		hostInfo := strings.Split(host, ",") // "module,hostname:port"
		if len(hostInfo) != 2 {
			continue
		}
		hostModule := hostInfo[0]
		hostNamePort := hostInfo[1]

		hostNamePortList := strings.Split(hostNamePort, ":")
		if len(hostNamePortList) != 2 {
			continue
		}
		hostName := hostNamePortList[0]
		hostPort := hostNamePortList[1]

		myTags := tags + ",module=" + hostModule + ",port=" + hostPort
		srcUrl := fmt.Sprintf(srcUrlFmt, hostNamePort)
		getResp, err := client.Get(srcUrl)
		if err != nil {
			log.Printf(hostNamePort+", get statistics error,", err)
			continue
		}
		defer getResp.Body.Close()

		body, err := ioutil.ReadAll(getResp.Body)
		if err != nil {
			log.Println(hostNamePort+", get statistics error,", err)
			continue
		}

		var data Dto
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println(hostNamePort+", get statistics error,", err)
			continue
		}

		for _, item := range data.Data {
			if item["Name"] == nil {
				continue
			}
			itemName := item["Name"].(string)

			if item["Cnt"] != nil {
				var jmdCnt model.MetricValue
				jmdCnt.Endpoint = hostName
				jmdCnt.Metric = itemName
				jmdCnt.Timestamp = ts
				jmdCnt.Step = 60
				jmdCnt.Value = int64(item["Cnt"].(float64))
				jmdCnt.Type = "GAUGE"
				jmdCnt.Tags = myTags
				jsonList = append(jsonList, jmdCnt)
			}

			if item["Qps"] != nil {
				var jmdQps model.MetricValue
				jmdQps.Endpoint = hostName
				jmdQps.Metric = itemName + ".Qps"
				jmdQps.Timestamp = ts
				jmdQps.Step = 60
				jmdQps.Value = int64(item["Qps"].(float64))
				jmdQps.Type = "GAUGE"
				jmdQps.Tags = myTags
				jsonList = append(jsonList, jmdQps)
			}
		}

		if len(jsonList) < 1 { //没取到数据
			log.Println("get null from ", hostNamePort)
			continue
		}

		// format result
		jsonBody, err := json.Marshal(jsonList)
		if err != nil {
			log.Println(hostNamePort+", format body error,", err)
			continue
		}

		// send by http-post
		req, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		postResp, err := client.Do(req)
		if err != nil {
			log.Println(hostNamePort+", post to dest error,", err)
			continue
		}
		defer postResp.Body.Close()

		if postResp.StatusCode/100 != 2 {
			log.Println(hostNamePort+", post to dest, bad response,", postResp.StatusCode)
		}
	}
}

type Dto struct {
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}
