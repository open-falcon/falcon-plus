package model

import (
	"math"

	"github.com/open-falcon/falcon-plus/common/model"
)

type Function interface {
	Compute(L *SafeELinkedList) (vs []*model.EHistoryData, leftValue float64, isTriggered bool, isEnough bool)
}

type AllFunction struct {
	Function
	Metric     string
	Limit      int
	Operator   string
	RightValue float64
}

func (this AllFunction) Compute(L *SafeELinkedList) (vs []*model.EHistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	for i := 0; i < this.Limit; i++ {
		value, ok := vs[i].Values[this.Metric]
		if !ok {
			break
		}
		isTriggered = checkIsTriggered(value, this.Operator, this.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Values[this.Metric]
	return
}

func checkIsTriggered(leftValue float64, operator string, rightValue float64) (isTriggered bool) {
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
