package model

import (
	"encoding/json"
	"fmt"
	"log"
)

type Condition struct {
	Func       string
	Metric     string
	Parameter  string
	Limit      int
	Operator   string
	RightValue float64
}

type EExpResponse struct {
	EExps []*EExp `json:"eexpressions"`
}

type EExp struct {
	ID         int
	Func       string
	Metric     string // join(sorted(conditionMetrics), ",")
	Filters    map[string]string
	Conditions []Condition
	Priority   int
	MaxStep    int
	Note       string
}

func (c *Condition) String() string {
	return fmt.Sprintf("func:%s metric:%s parameter:%s operator:%s rightValue:%s", c.Func, c.Metric, c.Parameter, c.Operator, c.RightValue)
}

func (c *Condition) Hit(m *EMetric) bool {
	v, ok := (*m).Values[c.Metric]
	if !ok {
		log.Println("metric matched condition not found")
		return false
	}

	switch c.Operator {
	case "<":
		{
			return v < c.RightValue
		}
	case "<=":
		{
			return v <= c.RightValue
		}
	case "==":
		{
			return v == c.RightValue
		}
	case ">":
		{
			return v <= c.RightValue
		}
	case ">=":
		{
			return v <= c.RightValue
		}
	case "<>":
		{
			return v != c.RightValue
		}
	}

	return true
}

func (ee *EExp) String() string {
	outF, _ := json.Marshal(ee.Filters)
	outC, _ := json.Marshal(ee.Conditions)
	return fmt.Sprintf("func:%s filters:%s conditions:%s", ee.Func, outF, outC)
}

func (ee *EExp) Hit(m *EMetric) bool {
	for k, v := range ee.Filters {
		if k == "metric" {
			continue
		}
		vGot, ok := m.Filters[k]
		if !ok {
			log.Println("filter not matched", k)
			return false
		}
		if vGot != v {
			log.Println("filter value not matched")
			return false
		}
	}

	for _, cond := range ee.Conditions {
		if !cond.Hit(m) {
			//log.Println("condition not hit", cond, m)
			return false
		}
	}

	return true
}

func (ee *EExp) HitFilters(m *map[string]interface{}) bool {
	for k, v := range ee.Filters {
		vGot, ok := (*m)[k].(string)
		if !ok || v != vGot {
			return false
		}
	}
	return true
}
