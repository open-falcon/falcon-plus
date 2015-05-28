package monitor

import (
	"bytes"
	"fmt"
	ncron "github.com/niean/cron"
	nhttp "github.com/niean/go-httpclient"
	nsema "github.com/niean/gotools/concurrent/semaphore"
	nmap "github.com/niean/gotools/container/nmap"
	ntime "github.com/niean/gotools/time"
	"github.com/open-falcon/task/g"
	"github.com/open-falcon/task/proc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	alarmInterval = time.Duration(300) * time.Second
)

var (
	monitorCron = ncron.New()
	sema        = nsema.NewSemaphore(1)
	statusCache = nmap.NewSafeMap()
	alarmCache  = nmap.NewSafeMap()
	cronSpec    = "0 * * * * ?"
)

func Start() {
	if g.Config().Monitor.Enabled {
		monitorCron.AddFunc(cronSpec, func() {
			monitor()
		})
		monitorCron.Start()
		go alarmJudge()
		log.Println("monitor.Start, ok")
	} else {
		log.Println("monitor.Start, not enable")
	}
}

// alarm judge
func alarmJudge() {
	interval := time.Duration(10) * time.Second
	for {
		time.Sleep(interval)
		var content bytes.Buffer

		keys := alarmCache.Keys()
		if len(keys) == 0 {
			continue
		}
		for _, key := range keys {
			aitem, found := alarmCache.GetAndRemove(key)
			if !found {
				continue
			}
			content.WriteString(aitem.(*Alarm).String() + "\n")
		}

		if content.Len() > 5 {
			mailContent := formAlarmMailContent(g.Config().Monitor.MailTos, "Self-Monitor.Alarm",
				content.String(), "Falcon")
			err := sendMail(g.Config().Monitor.MailUrl, mailContent)
			if err != nil {
				log.Println("alarm send mail error, mail:", mailContent, "", err)
			} else {
				// statistics
				proc.MonitorAlarmMailCnt.Incr()
			}
		}
	}
}

func formAlarmMailContent(tos string, subject string, content string, from string) string {
	return fmt.Sprintf("tos=%s;subject=%s;content=%s;user=%s", tos, subject, content, from)
}

func sendMail(mailUrl string, content string) error {
	transport := &nhttp.Transport{
		ConnectTimeout: time.Duration(5) * time.Second,
		RequestTimeout: time.Duration(20) * time.Second,
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}

	// send by http-post
	req, err := http.NewRequest("POST", mailUrl, bytes.NewBufferString(content))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	postResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	if postResp.StatusCode/100 != 2 {
		return fmt.Errorf("Http-Post Error, Code %d", postResp.StatusCode)
	}
	return nil
}

// status calc
func monitor() {
	if sema.AvailablePermits() <= 0 {
		proc.MonitorConcurrentErrorCnt.Incr()
		log.Println("monitor.collect, concurrent not available")
		return
	}
	sema.Acquire()
	defer sema.Release()

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
	transport := &nhttp.Transport{
		ConnectTimeout: time.Duration(1) * time.Second,
		RequestTimeout: time.Duration(10) * time.Second,
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}

	for _, host := range g.Config().Monitor.Cluster {
		hostInfo := strings.Split(host, ",") // "module,hostname:port/health"
		if len(hostInfo) != 2 {
			continue
		}
		//hostType := hostInfo[0]
		hostUrl := hostInfo[1]
		if !strings.Contains(hostUrl, "http://") {
			hostUrl = "http://" + hostUrl
		}

		getResp, err := client.Get(hostUrl)
		if err != nil {
			log.Printf(host+", monitor error,", err)
			onMonitorErr(host)
			continue
		}
		defer getResp.Body.Close()

		body, err := ioutil.ReadAll(getResp.Body)                        // body=['o','k',...]
		if !(err == nil && len(body) >= 2 && string(body[:2]) == "ok") { // err
			log.Println(host, ", error,", err)
			onMonitorErr(host)
		} else { // get "ok"
			onMonitorOk(host)
		}
	}
}

func onMonitorErr(host string) {
	// change status
	s, found := statusCache.Get(host)
	if !found {
		s = NewStatus()
		statusCache.Put(host, s)
	}
	ss := s.(*Status)
	ss.OnErr()

	// alarm
	errCnt := ss.GetErrCnt()
	if errCnt >= 4 && errCnt <= 16 {
		for i := 4; i <= errCnt; i *= 2 {
			if errCnt == i {
				a := NewAlarm(host, "err", ss.GetErrCnt())
				alarmCache.Put(host, a)
				break
			}
		}
	}
}

func onMonitorOk(host string) {
	// change status
	s, found := statusCache.Get(host)
	if !found {
		s = NewStatus()
		statusCache.Put(host, s)
	}
	ss := s.(*Status)
	errCnt := ss.GetErrCnt()
	ss.OnOk()

	if ss.IsTurnToOk() {
		if errCnt >= 4 { //有过alarm, 才能turnOk
			// alarm
			a := NewAlarm(host, "ok", ss.GetErrCnt())
			alarmCache.Put(host, a)
		}
	}
}

// Status Struct
type Status struct {
	sync.RWMutex
	Status     string
	LastStatus string
	ErrCnt     int
	OkCnt      int
}

func NewStatus() *Status {
	return &Status{Status: "ok", LastStatus: "ok", ErrCnt: 0, OkCnt: 0}
}

func (s *Status) GetErrCnt() int {
	s.RLock()
	cnt := s.ErrCnt
	s.RUnlock()
	return cnt
}

func (s *Status) OnErr() {
	s.Lock()
	s.LastStatus = s.Status
	s.Status = "err"
	s.OkCnt = 0
	s.ErrCnt += 1
	s.Unlock()
}

func (s *Status) OnOk() {
	s.Lock()
	s.LastStatus = s.Status
	s.Status = "ok"
	s.OkCnt += 1
	s.ErrCnt = 0
	s.Unlock()
}

func (s *Status) IsTurnToOk() bool {
	s.RLock()
	ret := false
	if s.LastStatus == "err" && s.Status == "ok" {
		ret = true
	}
	s.RUnlock()
	return ret
}

// AlarmItem Struct
type Alarm struct {
	ObjName   string
	AlarmType string
	AlarmCnt  int
	Ts        int64
}

func NewAlarm(obj string, atype string, cnt int) *Alarm {
	return &Alarm{AlarmType: atype, ObjName: obj, AlarmCnt: cnt, Ts: time.Now().Unix()}
}

func (a *Alarm) String() string {
	return fmt.Sprintf("[%s][%s][%d][%s]", ntime.FormatTs(a.Ts), a.AlarmType, a.AlarmCnt, a.ObjName)
}
