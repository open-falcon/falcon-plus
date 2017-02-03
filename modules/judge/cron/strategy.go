package cron

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"log"
	"time"
)

func SyncStrategies() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncStrategies()
		syncExpression()
		time.Sleep(duration)
	}
}

func syncStrategies() {
	var strategiesResponse model.StrategiesResponse
	err := g.HbsClient.Call("Hbs.GetStrategies", model.NullRpcRequest{}, &strategiesResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetStrategies:", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

func rebuildStrategyMap(strategiesResponse *model.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]model.Strategy)
	for _, hs := range strategiesResponse.HostStrategies {
		hostname := hs.Hostname
		if g.Config().Debug && hostname == g.Config().DebugHost {
			log.Println(hostname, "strategies:")
			bs, _ := json.Marshal(hs.Strategies)
			fmt.Println(string(bs))
		}
		for _, strategy := range hs.Strategies {
			key := fmt.Sprintf("%s/%s", hostname, strategy.Metric)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], strategy)
			} else {
				m[key] = []model.Strategy{strategy}
			}
		}
	}

	g.StrategyMap.ReInit(m)
}

func syncExpression() {
	var expressionResponse model.ExpressionResponse
	err := g.HbsClient.Call("Hbs.GetExpressions", model.NullRpcRequest{}, &expressionResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetExpressions:", err)
		return
	}

	rebuildExpressionMap(&expressionResponse)
}

func rebuildExpressionMap(expressionResponse *model.ExpressionResponse) {
	m := make(map[string][]*model.Expression)
	for _, exp := range expressionResponse.Expressions {
		for k, v := range exp.Tags {
			key := fmt.Sprintf("%s/%s=%s", exp.Metric, k, v)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], exp)
			} else {
				m[key] = []*model.Expression{exp}
			}
		}
	}

	g.ExpressionMap.ReInit(m)
}
