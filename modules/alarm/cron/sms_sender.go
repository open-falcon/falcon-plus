package cron

import (
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
	"time"
)

func ConsumeSms() {
	for {
		L := redi.PopAllSms()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*model.Sms) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
		go SendChat(sms)
	}
}

func SendSms(sms *model.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config().Api.Sms
	resp, err := post(url, sms)
	if err != nil {
		log.Errorf("send sms fail, tos:%s, cotent:%s, error:%v", sms.Tos, sms.Content, err)
	}

	log.Debugf("send sms:%v, resp:%v, url:%s", sms, resp, url)
}

func SendChat(sms *model.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config().Api.Chat
	resp, err := post(url, sms)
	if err != nil {
		log.Errorf("send chat message fail, tos:%s, cotent:%s, error:%v", sms.Tos, sms.Content, err)
	}

	log.Debugf("send chat message:%v, resp:%v, url:%s", sms, resp, url)
}

func post(url string, sms *model.Sms) (string, error) {
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Param("tos", sms.Tos)
	r.Param("content", sms.Content)
	return r.String()
}

