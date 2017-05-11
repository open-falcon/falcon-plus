package cron

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
)

func HandleCallback(event *model.Event, action *api.Action) {

	teams := action.Uic
	phones := []string{}
	mails := []string{}
	ims := []string{}

	if teams != "" {
		phones, mails, ims = api.ParseTeams(teams)
		smsContent := GenerateSmsContent(event)
		mailContent := GenerateMailContent(event)
		chatContent := GenerateChatContent(event)
		if action.BeforeCallbackSms == 1 {
			redi.WriteSms(phones, smsContent)
			redi.WriteChat(ims, chatContent)
		}

		if action.BeforeCallbackMail == 1 {
			redi.WriteMail(mails, smsContent, mailContent)
		}
	}

	message := Callback(event, action)

	if teams != "" {
		if action.AfterCallbackSms == 1 {
			redi.WriteSms(phones, message)
			redi.WriteChat(ims, message)
		}

		if action.AfterCallbackMail == 1 {
			redi.WriteMail(mails, message, message)
		}
	}

}

func Callback(event *model.Event, action *api.Action) string {
	if action.Url == "" {
		return "callback url is blank"
	}

	L := make([]string, 0)
	if len(event.PushedTags) > 0 {
		for k, v := range event.PushedTags {
			L = append(L, fmt.Sprintf("%s:%s", k, v))
		}
	}

	tags := ""
	if len(L) > 0 {
		tags = strings.Join(L, ",")
	}

	req := httplib.Get(action.Url).SetTimeout(3*time.Second, 20*time.Second)

	req.Param("endpoint", event.Endpoint)
	req.Param("metric", event.Metric())
	req.Param("status", event.Status)
	req.Param("step", fmt.Sprintf("%d", event.CurrentStep))
	req.Param("priority", fmt.Sprintf("%d", event.Priority()))
	req.Param("time", event.FormattedTime())
	req.Param("tpl_id", fmt.Sprintf("%d", event.TplId()))
	req.Param("exp_id", fmt.Sprintf("%d", event.ExpressionId()))
	req.Param("stra_id", fmt.Sprintf("%d", event.StrategyId()))
	req.Param("tags", tags)

	resp, e := req.String()

	success := "success"
	if e != nil {
		log.Errorf("callback fail, action:%v, event:%s, error:%s", action, event.String(), e.Error())
		success = fmt.Sprintf("fail:%s", e.Error())
	}
	message := fmt.Sprintf("curl %s %s. resp: %s", action.Url, success, resp)
	log.Debugf("callback to url:%s, event:%s, resp:%s", action.Url, event.String(), resp)

	return message
}
