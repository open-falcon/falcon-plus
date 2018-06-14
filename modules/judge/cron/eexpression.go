package cron

import (
	"log"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SyncEExpresions() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncEExpressions()
		syncEFilter()
		time.Sleep(duration)
	}
}

func syncEExpressions() {
	var eeResp model.EExpressionResponse
	err := g.HbsClient.Call("Judge.GetEExpressions", model.NullRpcRequest{}, &eeResp)
	if err != nil {
		log.Println("[ERROR] Judge.GetEExpressions:", err)
		return
	}

	rebuildEExpressionMap(&eeResp)
}

func rebuildEExpressionMap(eeResp *model.EExpressionResponse) {
	m := make(map[string][]*model.EExpression)
	for _, exp := range eeResp.EExpressions {
		if _, exists := m[exp.Metric]; exists {
			m[exp.Metric] = append(m[exp.Metric], exp)
		} else {
			m[exp.Metric] = []*model.EExpression{exp}
		}
	}

	g.EExpressionMap.ReInit(m)
}

func syncEFilter() {
	m := make(map[string]string)

	//M map[string][]*model.EExpression
	eeMap := g.EExpressionMap.Get()
	for _, ees := range eeMap {
		for _, eexp := range ees {
			m[eexp.Metric] = eexp.Metric
		}
	}

	g.EFilterMap.ReInit(m)
}
