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

package store

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"math"
	"strconv"
	"strings"
)

type Function interface {
	Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool)
}

type MaxFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this MaxFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	max := vs[0].Value
	for i := 1; i < this.Limit; i++ {
		if max < vs[i].Value {
			max = vs[i].Value
		}
	}

	leftValue = max
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type MinFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this MinFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	min := vs[0].Value
	for i := 1; i < this.Limit; i++ {
		if min > vs[i].Value {
			min = vs[i].Value
		}
	}

	leftValue = min
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type AllFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this AllFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	for i := 0; i < this.Limit; i++ {
		isTriggered = checkIsTriggered(vs[i].Value, this.Operator, this.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Value
	return
}

type LookupFunction struct {
	Function
	Num        int
	Limit      int
	Operator   string
	RightValue float64
}

func (this LookupFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	leftValue = vs[0].Value

	for n, i := 0, 0; i < this.Limit; i++ {
		if checkIsTriggered(vs[i].Value, this.Operator, this.RightValue) {
			n++
			if n == this.Num {
				isTriggered = true
				return
			}
		}
	}

	return
}

type SumFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this SumFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < this.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type AvgFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this AvgFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < this.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum / float64(this.Limit)
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type DiffFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

// 只要有一个点的diff触发阈值，就报警
func (this DiffFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	// 此处this.Limit要+1，因为通常说diff(#3)，是当前点与历史的3个点相比较
	// 然而最新点已经在linkedlist的第一个位置，所以……
	vs, isEnough = L.HistoryData(this.Limit + 1)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	first := vs[0].Value

	isTriggered = false
	for i := 1; i < this.Limit+1; i++ {
		// diff是当前值减去历史值
		leftValue = first - vs[i].Value
		isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
		if isTriggered {
			break
		}
	}

	return
}

type StdDeviationFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

/*
	离群点检测函数，更多请参考3-sigma算法：https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule
	stddev(#10) = 3 //取最新 **10** 个点的数据分别计算得到他们的标准差和均值，分别计为 σ 和 μ，其中当前值计为 X，那么当 X 落在区间 [μ-3σ, μ+3σ] 之外时则报警。
*/

func (this StdDeviationFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	leftValue = vs[0].Value

	var datas []float64
	for _, i := range vs {
		datas = append(datas, i.Value)
	}

	isTriggered = false

	std := utils.ComputeStdDeviation(datas)
	mean := utils.ComputeMean(datas)

	upperBound := mean + this.RightValue*std
	lowerBound := mean - this.RightValue*std

	if leftValue < lowerBound || leftValue > upperBound {
		isTriggered = true
	}

	return
}

// pdiff(#3)
type PDiffFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this PDiffFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit + 1)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	first := vs[0].Value

	isTriggered = false
	for i := 1; i < this.Limit+1; i++ {
		if vs[i].Value == 0 {
			continue
		}

		leftValue = (first - vs[i].Value) / vs[i].Value * 100.0
		isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
		if isTriggered {
			break
		}
	}

	return
}


//            3
//          ____
//   3    /
//  ____/
//
//
//kpdiff(#3,3)  告警上面类似的数据模型，解决某个指标长期稳定，当指标变化并且是持续的情况下告警。原来的diff 跟pdiff是只要一个点变化就告警。
//kpdiff(#3,3)  左边的3是指稳定数据走势的数据点，右边的3是指持续变化的数据点

type KPDiffFunction struct {
	Function
	Num		   int
	Limit      int
	Operator   string
	RightValue float64
}

func (this KPDiffFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit + this.Num)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}
	for i := 0; i < this.Num; i++ {
		isTriggered = false
		if vs[i].Value == 0 {
			break
		}

		//kpdiff(#3,3) 全部右边的点都对全部左边的点相减，得到3*3个差，再将3*3个差值分别除以对应左边的点，得到3*3个商值，全部商值满足阈值则报警
		for j := 0; j < this.Limit; j++ {
			leftValue = (vs[j].Value - vs[this.Limit + this.Num -1 -i].Value) / vs[j].Value * 100.0
			isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
			if isTriggered == false {
				break
			}
		}
		if isTriggered == false {
			return
		}
	}
	return
}

type KDiffFunction struct {
	Function
	Num        int
	Limit      int
	Operator   string
	RightValue float64
}

func (this KDiffFunction) Compute(L *SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit + this.Num)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}
	for i := 0; i < this.Num; i++ {
		isTriggered = false
		if vs[i].Value == 0 {
			break
		}

		for j := 0; j < this.Limit; j++ {
			leftValue = vs[j].Value - vs[this.Limit + this.Num -1 -i].Value
			isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
			if isTriggered == false {
				break
			}
		}
		if isTriggered == false {
			return
		}
	}
	return
}

func atois(s string) (ret []int, err error) {
	a := strings.Split(s, ",")
	ret = make([]int, len(a))
	for i, v := range a {
		ret[i], err = strconv.Atoi(v)
		if err != nil {
			return
		}
	}
	return
}

// @str: e.g. all(#3) sum(#3) avg(#10) diff(#10) stddev(#10)
func ParseFuncFromString(str string, operator string, rightValue float64) (fn Function, err error) {
	if str == "" {
		return nil, fmt.Errorf("func can not be null!")
	}
	idx := strings.Index(str, "#")
	args, err := atois(str[idx+1 : len(str)-1])
	if err != nil {
		return nil, err
	}

	switch str[:idx-1] {
	case "max":
		fn = &MaxFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "min":
		fn = &MinFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "all":
		fn = &AllFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "sum":
		fn = &SumFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "avg":
		fn = &AvgFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "diff":
		fn = &DiffFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "pdiff":
		fn = &PDiffFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "lookup":
		fn = &LookupFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	case "stddev":
		fn = &StdDeviationFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "kdiff":
		fn = &KDiffFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	case "kpdiff":
		fn = &KPDiffFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	default:
		err = fmt.Errorf("not_supported_method")
	}

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
