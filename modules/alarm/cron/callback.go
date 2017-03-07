package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/redis"
	"github.com/toolkits/net/httplib"
)

func HandleCallback(event *model.Event, action *api.Action) {

	// falcon,dinp
	teams := action.Uic
	phones := []string{}
	mails := []string{}

	if teams != "" {
		phones, mails = api.ParseTeams(teams)
		smsContent := GenerateSmsContent(event)
		mailContent := GenerateMailContent(event)
		if action.BeforeCallbackSms == 1 {
			redis.WriteSms(phones, smsContent)
		}

		if action.BeforeCallbackMail == 1 {
			redis.WriteMail(mails, smsContent, mailContent)
		}
	}

	message := Callback(event, action)

	if teams != "" {
		if action.AfterCallbackSms == 1 {
			redis.WriteSms(phones, message)
		}

		if action.AfterCallbackMail == 1 {
			redis.WriteMail(mails, message, message)
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
		success = fmt.Sprintf("fail:%s", e.Error())
	}
	message := fmt.Sprintf("curl %s %s. resp: %s", action.Url, success, resp)

	return message
}
