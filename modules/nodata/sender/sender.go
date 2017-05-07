package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	tsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/nmap"
	thttpclient "github.com/toolkits/http/httpclient"
	ttime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

var (
	MockMap = nmap.NewSafeMap()
	sema    = tsema.NewSemaphore(1)
)

func AddMock(key string, endpoint string, metric string, tags string, ts int64, dstype string, step int64, value interface{}) {
	item := &cmodel.JsonMetaData{metric, endpoint, ts, step, value, dstype, tags}
	MockMap.Put(key, item)
}

func SendMockOnceAsync() {
	go SendMockOnce()
}

func SendMockOnce() int {
	if !sema.TryAcquire() {
		return -1
	}
	defer sema.Release()

	// not enabled
	if !g.Config().Sender.Enabled {
		return 0
	}

	start := time.Now().Unix()
	cnt, _ := sendMock()
	end := time.Now().Unix()
	if g.Config().Debug {
		log.Printf("sender cron, cnt %d, time %ds, start %s", cnt, end-start, ttime.FormatTs(start))
	}

	// statistics
	g.SenderCronCnt.Incr()
	g.SenderLastTs.SetCnt(end - start)
	g.SenderCnt.IncrBy(int64(cnt))

	return cnt
}

func sendMock() (cnt int, errt error) {

	cfg := g.Config().Sender
	batch := int(cfg.Batch)
	connTimeout := cfg.ConnectTimeout
	requTimeout := cfg.RequestTimeout

	// send mock to transfer
	mocks := MockMap.Slice()
	MockMap.Clear()
	mockSize := len(mocks)
	i := 0
	for i < mockSize {
		leftLen := mockSize - i
		sendSize := batch
		if leftLen < sendSize {
			sendSize = leftLen
		}
		fetchMocks := mocks[i : i+sendSize]
		i += sendSize

		items := make([]*cmodel.JsonMetaData, 0)
		for _, val := range fetchMocks {
			if val == nil {
				continue
			}
			items = append(items, val.(*cmodel.JsonMetaData))
		}
		cntonce, err := sendItemsToTransfer(items, len(items), "nodata.mock",
			time.Millisecond*time.Duration(connTimeout),
			time.Millisecond*time.Duration(requTimeout))
		if err == nil {
			if g.Config().Debug {
				log.Println("send items:", items)
			}
			cnt += cntonce
		}
	}

	return cnt, nil
}

//
func sendItemsToTransfer(items []*cmodel.JsonMetaData, size int, httpcliname string,
	connT, reqT time.Duration) (cnt int, errt error) {
	if size < 1 {
		return
	}

	cfg := g.Config()
	transUlr := fmt.Sprintf("http://%s/api/push", cfg.Sender.TransferAddr)
	hcli := thttpclient.GetHttpClient(httpcliname, connT, reqT)

	// form request args
	itemsBody, err := json.Marshal(items)
	if err != nil {
		log.Println(transUlr+", format body error,", err)
		errt = err
		return
	}

	// fetch items
	req, err := http.NewRequest("POST", transUlr, bytes.NewBuffer(itemsBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Connection", "close")
	postResp, err := hcli.Do(req)
	if err != nil {
		log.Println(transUlr+", post to dest error,", err)
		errt = err
		return
	}
	defer postResp.Body.Close()

	if postResp.StatusCode/100 != 2 {
		log.Println(transUlr+", post to dest, bad response,", postResp.Body)
		errt = fmt.Errorf("request failed, %s", postResp.Body)
		return
	}

	return size, nil
}
