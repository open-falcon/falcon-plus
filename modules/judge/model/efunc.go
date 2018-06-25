package model

import (
	"math"

	"github.com/open-falcon/falcon-plus/common/model"
)

type Function interface {
	Compute(L *SafeELinkedList) (vs []*model.EHistoryData, leftValue interface{}, isTriggered bool, isEnough bool)
}

type AllFunction struct {
	Function
	Key        string
	Limit      uint64
	Operator   string
	RightValue interface{}
}

func (this AllFunction) Compute(L *SafeELinkedList) (vs []*model.EHistoryData, leftValue interface{}, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	var i uint64
	for i = 0; i < this.Limit; i++ {
		value, ok := vs[i].Filters[this.Key]
		if !ok {
			break
		}
		isTriggered = checkIsTriggered(value, this.Operator, this.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Filters[this.Key]
	return
}

type CountFunction struct {
	Function
	Key        string
	Ago        uint64
	Hits       uint64
	Operator   string
	RightValue interface{}
	Now        int64
}

func (this CountFunction) Compute(L *SafeELinkedList) (vs []*model.EHistoryData, leftValue interface{}, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryDataByTime(this.Ago, this.Hits, uint64(this.Now))
	if !isEnough {
		return
	}

	isTriggered = true
	var i uint64
	for i = 0; i < this.Hits; i++ {
		value, ok := vs[i].Filters[this.Key]
		if !ok {
			break
		}
		isTriggered = checkIsTriggered(value, this.Operator, this.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Filters[this.Key]
	return
}

func checkIsTriggered(leftValueI interface{}, operator string, rightValueI interface{}) (isTriggered bool) {
	switch rightValueI.(type) {
	case string:
		{
			leftValue := leftValueI.(string)
			rightValue := rightValueI.(string)
			switch operator {
			case "=", "==":
				isTriggered = leftValue == rightValue
			case "!=":
				isTriggered = leftValue != rightValue
			}

		}
	case float64:
		{
			leftValue := leftValueI.(float64)
			rightValue := rightValueI.(float64)
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
		}
	}

	return
}
