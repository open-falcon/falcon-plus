package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cron "github.com/toolkits/cron"
	nhttpclient "github.com/toolkits/http/httpclient"
	ntime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/task/g"
	"github.com/open-falcon/falcon-plus/modules/task/proc"
)

var (
	collectorCron = cron.New()
	srcUrlFmt     = "http://%s/statistics/all"
	destUrl       = "http://127.0.0.1:1988/v1/push"
)

func Start() {
	if !g.Config().Collector.Enable {
		log.Println("collector.Start warning, not enable")
		return
	}

	// init url
	if g.Config().Collector.DestUrl != "" {
		destUrl = g.Config().Collector.DestUrl
	}
	if g.Config().Collector.SrcUrlFmt != "" {
		srcUrlFmt = g.Config().Collector.SrcUrlFmt
	}
	// start
	go startCollectorCron()
	log.Println("collector.Start, ok")
}

func startCollectorCron() {
	collectorCron.AddFuncCC("0 * * * * ?", func() { collect() }, 1)
	collectorCron.Start()
}

func collect() {
	startTs := time.Now().Unix()
	_collect()
	endTs := time.Now().Unix()
	log.Printf("collect, start %s, ts %ds\n", ntime.FormatTs(startTs), endTs-startTs)

	// statistics
	proc.CollectorCronCnt.Incr()
}

func _collect() {
	clientGet := nhttpclient.GetHttpClient("collector.get", 10*time.Second, 20*time.Second)
	tags := "type=statistics,pdl=falcon"
	for _, host := range g.Config().Collector.Cluster {
		ts := time.Now().Unix()
		jsonList := make([]*cmodel.JsonMetaData, 0)

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
		reqGet, _ := http.NewRequest("GET", srcUrl, nil)
		reqGet.Header.Set("Connection", "close")
		getResp, err := clientGet.Do(reqGet)
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
				var jmdCnt cmodel.JsonMetaData
				jmdCnt.Endpoint = hostName
				jmdCnt.Metric = itemName
				jmdCnt.Timestamp = ts
				jmdCnt.Step = 60
				jmdCnt.Value = int64(item["Cnt"].(float64))
				jmdCnt.CounterType = "GAUGE"
				jmdCnt.Tags = myTags
				jsonList = append(jsonList, &jmdCnt)
			}

			if item["Qps"] != nil {
				var jmdQps cmodel.JsonMetaData
				jmdQps.Endpoint = hostName
				jmdQps.Metric = itemName + ".Qps"
				jmdQps.Timestamp = ts
				jmdQps.Step = 60
				jmdQps.Value = int64(item["Qps"].(float64))
				jmdQps.CounterType = "GAUGE"
				jmdQps.Tags = myTags
				jsonList = append(jsonList, &jmdQps)
			}
		}

		// format result
		err = sendToTransfer(jsonList, destUrl)
		if err != nil {
			log.Println(hostNamePort, "send to transfer error,", err.Error())
		}
	}

	// collector.alive
	_collectorAlive()
}

func _collectorAlive() error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("get hostname failed,", err)
		return err
	}

	var jmdCnt cmodel.JsonMetaData
	jmdCnt.Endpoint = hostname
	jmdCnt.Metric = "falcon.task.alive"
	jmdCnt.Timestamp = time.Now().Unix()
	jmdCnt.Step = 60
	jmdCnt.Value = 0
	jmdCnt.CounterType = "GAUGE"
	jmdCnt.Tags = ""

	jsonList := make([]*cmodel.JsonMetaData, 0)
	jsonList = append(jsonList, &jmdCnt)
	err = sendToTransfer(jsonList, destUrl)
	if err != nil {
		log.Println("send task.alive failed,", err)
		return err
	}

	return nil
}

func sendToTransfer(items []*cmodel.JsonMetaData, destUrl string) error {
	if len(items) < 1 {
		return nil
	}

	// format result
	jsonBody, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("json.Marshal failed with %v", err)
	}

	// send by http-post
	req, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "close")
	clientPost := nhttpclient.GetHttpClient("collector.post", 5*time.Second, 10*time.Second)
	postResp, err := clientPost.Do(req)
	if err != nil {
		return fmt.Errorf("post to %s, resquest failed with %v", destUrl, err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode/100 != 2 {
		return fmt.Errorf("post to %s, got bad response %d", destUrl, postResp.StatusCode)
	}

	return nil
}

type Dto struct {
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}
