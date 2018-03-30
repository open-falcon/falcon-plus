// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SyncStrategies() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncStrategies()
		syncExpression()
		syncFilter()
		rebuildStrMatcherMap()
		rebuildStrMatcherExpMap()
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

func syncFilter() {
	m := make(map[string]string)

	//M map[string][]model.Strategy
	strategyMap := g.StrategyMap.Get()
	for _, strategies := range strategyMap {
		for _, strategy := range strategies {
			m[strategy.Metric] = strategy.Metric
		}
	}

	//M map[string][]*model.Expression
	expressionMap := g.ExpressionMap.Get()
	for _, expressions := range expressionMap {
		for _, expression := range expressions {
			m[expression.Metric] = expression.Metric
		}
	}

	g.FilterMap.ReInit(m)
}

func parsePatternFromFunc(s string) (pattern string) {
	NOT_FOUND := -1
	idxMatchBracket := strings.Index(s, "match(")
	idxComma := strings.LastIndex(s, ",")
	if idxMatchBracket != NOT_FOUND && idxComma != NOT_FOUND {
		pattern = s[len("match("):idxComma]
	}
	return pattern
}

func rebuildStrMatcherMap() {
	m := make(map[string]map[string]*regexp.Regexp)

	strategyMap := g.StrategyMap.Get()
	for endpointSlashMetric, strategies := range strategyMap {
		parts := strings.Split(endpointSlashMetric, "/")
		if len(parts) < 1 {
			continue
		}
		endpoint := parts[0]

		for _, strategy := range strategies {
			if strategy.Metric != "str.match" {
				continue
			}

			if strategy.Func == "" {
				log.Println(`WARN: strategy.Func are empty`, strategy)
				continue
			}

			pattern := parsePatternFromFunc(strategy.Func)
			if pattern == "" {
				log.Println(`WARN: pattern is empty or parse pattern failed`, strategy.Func)
				continue
			}

			// auto append prefix to ignore case
			re, err := regexp.Compile(`(?i)` + pattern)
			if err != nil {
				log.Println(`WARN: compiling pattern failed`, pattern)
				continue
			}

			if _, ok := m[endpoint]; !ok {
				subM := make(map[string]*regexp.Regexp)
				m[endpoint] = subM

			}
			m[endpoint][pattern] = re
		}
	}

	g.StrMatcherMap.ReInit(m)
}

func rebuildStrMatcherExpMap() {
	m := make(map[string]map[string]*regexp.Regexp)

	exps := g.ExpressionMap.Get()

	for metricSlashTag, exps := range exps {
		parts := strings.Split(metricSlashTag, "/")
		if len(parts) < 2 {
			log.Println("WARN: parse metric from g.ExpressionMap failed", metricSlashTag)
			continue
		}
		metric := parts[0]
		if metric != "str.match" {
			continue
		}

		tag := parts[1]

		for _, exp := range exps {
			if exp.Func == "" {
				log.Println(`WARN: expression.Func are empty`, exp)
				continue
			}

			pattern := parsePatternFromFunc(exp.Func)
			if pattern == "" {
				log.Println(`WARN: pattern is empty or parse pattern failed`, exp.Func)
				continue
			}

			// auto append prefix to ignore case
			re, err := regexp.Compile(`(?i)` + pattern)
			if err != nil {
				log.Println(`WARN: compiling pattern failed`, pattern)
				continue
			}

			if _, ok := m[tag]; !ok {
				subM := make(map[string]*regexp.Regexp)
				m[tag] = subM
			}
			m[tag][pattern] = re
		}
	}

	g.StrMatcherExpMap.ReInit(m)
}
