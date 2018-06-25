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
	var resp model.EExpResponse
	err := g.HbsClient.Call("Judge.GetEExps", model.NullRpcRequest{}, &resp)
	if err != nil {
		log.Println("[ERROR] Judge.GetEExps:", err)
		return
	}

	rebuildEExpMap(&resp)
}

func rebuildEExpMap(resp *model.EExpResponse) {
	m := make(map[string][]model.EExp)
	for _, exp := range resp.EExps {
		if _, exists := m[exp.Key]; exists {
			m[exp.Key] = append(m[exp.Key], exp)
		} else {
			m[exp.Key] = []model.EExp{exp}
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
			m[eexp.Key] = eexp.Key
		}
	}

	g.EFilterMap.ReInit(m)
}
