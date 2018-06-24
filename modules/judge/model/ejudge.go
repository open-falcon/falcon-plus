package model

import (
	"encoding/json"
	"fmt"
	"log"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func EJudge(L *SafeELinkedList, firstItem *cmodel.EMetric, now int64) {
	CheckEExp(L, firstItem, now)
}

func sendEvent(event *cmodel.Event) {
	// update last event
	g.LastEvents.Set(event.Id, event)

	bs, err := json.Marshal(event)
	if err != nil {
		log.Printf("json marshal event %v fail: %v", event, err)
		return
	}

	// send to redis
	cfg := g.Config()
	if cfg.Alarm == nil {
		log.Println("alarm configuration not found")
		return
	}

	redisKey := fmt.Sprintf(cfg.Alarm.QueuePattern, event.Priority())
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	rc.Do("LPUSH", redisKey, string(bs))
}

func CheckEExp(L *SafeELinkedList, firstItem *cmodel.EMetric, now int64) {
	// expression可能会被多次重复处理，用此数据结构保证只被处理一次
	handledExpression := make(map[int]bool)

	m := g.EExpMap.Get()

	for _, eexps := range m {
		for _, eexp := range eexps {
			_, ok := handledExpression[(*eexp).ID]

			if (*eexp).Metric == firstItem.Metric && !ok {
				if judgeItemWithExpression(L, eexp, firstItem, now) {
					handledExpression[(*eexp).ID] = true
				}
			}
		}
	}

}

func HitFilters(eexp *cmodel.EExp, m *cmodel.EMetric) bool {
	for k, v := range eexp.Filters {
		vGot, ok := (*m).Filters[k]
		if !ok || v != vGot {
			return false
		}
	}
	return true
}

func judgeItemWithExpression(L *SafeELinkedList, eexp *cmodel.EExp, firstItem *cmodel.EMetric, now int64) (hit bool) {
	if !HitFilters(eexp, firstItem) {
		return false
	}

	var fn Function
	var historyData []*cmodel.EHistoryData
	var leftValue float64
	var isTriggered bool
	var isEnough bool
	for _, cond := range eexp.Conditions {
		switch cond.Func {
		case "all":
			fn = &AllFunction{Limit: cond.Limit, Operator: cond.Operator, Metric: cond.Metric, RightValue: cond.RightValue}
		default:
			{
				log.Println(fmt.Sprintf("not support function -%#v-", cond.Func))
				return false
			}
		}

		historyData, leftValue, isTriggered, isEnough = fn.Compute(L)
		if !isEnough {
			break
		}
		if !isTriggered {
			break
		}
	}

	if !isEnough {
		return false
	}
	if !isTriggered {
		return false
	}

	event := &cmodel.Event{
		Id:          fmt.Sprintf("e_%d_%s", eexp.ID, firstItem.PK()),
		EExp: eexp,
		Endpoint:    firstItem.Endpoint,
		LeftValue:   leftValue,
		EventTime:   firstItem.Timestamp,
		PushedTags:  firstItem.Filters,
	}

	sendEventIfNeed(historyData, isTriggered, now, event, eexp.MaxStep)

	return true

}

func sendEventIfNeed(historyData []*cmodel.EHistoryData, isTriggered bool, now int64, event *cmodel.Event, maxStep int) {
	lastEvent, exists := g.LastEvents.Get(event.Id)
	if isTriggered {
		event.Status = "PROBLEM"
		if !exists || lastEvent.Status[0] == 'O' {
			// 本次触发了阈值，之前又没报过警，得产生一个报警Event
			event.CurrentStep = 1

			// 但是有些用户把最大报警次数配置成了0，相当于屏蔽了，要检查一下
			if maxStep == 0 {
				return
			}

			sendEvent(event)
			return
		}

		// 逻辑走到这里，说明之前Event是PROBLEM状态
		if lastEvent.CurrentStep >= maxStep {
			// 报警次数已经足够多，到达了最多报警次数了，不再报警
			return
		}

		if historyData[len(historyData)-1].Timestamp <= lastEvent.EventTime {
			// 产生过报警的点，就不能再使用来判断了，否则容易出现一分钟报一次的情况
			// 只需要拿最后一个historyData来做判断即可，因为它的时间最老
			return
		}

		if now-lastEvent.EventTime < g.Config().Alarm.MinInterval {
			// 报警不能太频繁，两次报警之间至少要间隔MinInterval秒，否则就不能报警
			return
		}

		event.CurrentStep = lastEvent.CurrentStep + 1
		sendEvent(event)
	} else {
		// 如果LastEvent是Problem，报OK，否则啥都不做
		if exists && lastEvent.Status[0] == 'P' {
			event.Status = "OK"
			event.CurrentStep = 1
			sendEvent(event)
		}
	}
}
