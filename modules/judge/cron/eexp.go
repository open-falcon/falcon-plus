package cron

import (
	"log"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SyncEExps() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncEExps()
		syncEFilter()
		time.Sleep(duration)
	}
}

func syncEExps() {
	var eeResp model.EExpResponse
	err := g.HbsClient.Call("Judge.GetEExps", model.NullRpcRequest{}, &eeResp)
	if err != nil {
		log.Println("[ERROR] Judge.GetEExps:", err)
		return
	}

	rebuildEExpMap(&eeResp)
}

func rebuildEExpMap(eeResp *model.EExpResponse) {
	m := make(map[string][]*model.EExp)
	for _, exp := range eeResp.EExps {
		if _, exists := m[exp.Metric]; exists {
			m[exp.Metric] = append(m[exp.Metric], exp)
		} else {
			m[exp.Metric] = []*model.EExp{exp}
		}
	}

	g.EExpMap.ReInit(m)
}

func syncEFilter() {
	m := make(map[string]string)

	//M map[string][]*model.EExp
	eeMap := g.EExpMap.Get()
	for _, ees := range eeMap {
		for _, eexp := range ees {
			m[eexp.Metric] = eexp.Metric
		}
	}

	g.EFilterMap.ReInit(m)
}
