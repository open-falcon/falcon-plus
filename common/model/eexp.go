package model

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
)

type Filter struct {
	Func       string      `json:"func"`
	Key        string      `json:"key"`
	Ago        uint64      `json:"ago"`   // for func count
	Hits       uint64      `json:"hits"`  // for func count
	Limit      uint64      `json:"limit"` // for func all
	Operator   string      `json:"operator"`
	RightValue interface{} `json:"rightValue"`
}

type EExpResponse struct {
	EExps []EExp `json:"eexps"`
}

type EExp struct {
	ID       int               `json:"id"`
	Key      string            `json:"key"` // join(sorted(conditionKeys), ",")
	Filters  map[string]Filter `json:"filters"`
	Priority int               `json:"priority"`
	MaxStep  int               `json:"maxStep"`
	Note     string            `json:"note"`
}

func (c *Filter) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}

func (ee *EExp) String() string {
	out, _ := json.Marshal(ee)
	return string(out)
}

func opResultFloat64(leftValue float64, operator string, rightValue float64) (isTriggered bool) {
	switch operator {
	case "=", "==":
		isTriggered = math.Abs(leftValue-rightValue) < 0.0001
	case "!=":
		isTriggered = math.Abs(leftValue-rightValue) > 0.0001
	case "<":
		isTriggered = leftValue < rightValue
	case "<=":
		isTriggered = leftValue <= rightValue
	case ">":
		isTriggered = leftValue > rightValue
	case ">=":
		isTriggered = leftValue >= rightValue
	}
	return
}

func opResultString(leftValue string, operator string, rightValue string) (isTriggered bool) {
	switch operator {
	case "=", "==":
		isTriggered = leftValue == rightValue
	case "!=":
		isTriggered = leftValue != rightValue
	}
	return
}

func (ee *EExp) HitFilters(m *map[string]interface{}) bool {
	var err error
	for k, filter := range ee.Filters {
		valueI, ok := (*m)[k]
		if !ok {
			return false
		}

		switch filter.RightValue.(type) {
		case float64:
			{

				var leftValue float64
				leftValueStr, ok := valueI.(string)
				if ok {
					leftValue, err = strconv.ParseFloat(leftValueStr, 64)
					if err != nil {
						log.Println("strconv.ParseFloat failed", err)
						return false
					}
				} else {
					leftValue, ok = valueI.(float64)
					if !ok {
						log.Println("parse in float64 failed")
						return false
					}
				}

				rightValue := filter.RightValue.(float64)

				if !opResultFloat64(leftValue, filter.Operator, rightValue) {
					return false
				}

			}
		case string:
			{
				leftValue, ok := valueI.(string)
				if !ok {
					log.Println("parse in string failed")
					return false
				}

				rightValue := filter.RightValue.(string)

				if !opResultString(leftValue, filter.Operator, rightValue) {
					return false
				}

			}
		}

	}
	return true
}
