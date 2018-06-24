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

package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

// 机器监控和实例监控都会产生Event，共用这么一个struct
type Event struct {
	Id          string            `json:"id"`
	Strategy    *Strategy         `json:"strategy"`
	Expression  *Expression       `json:"expression"`
	EExp *EExp      `json:"eexp"`
	Status      string            `json:"status"` // OK or PROBLEM
	Endpoint    string            `json:"endpoint"`
	LeftValue   float64           `json:"leftValue"`
	CurrentStep int               `json:"currentStep"`
	EventTime   int64             `json:"eventTime"`
	PushedTags  map[string]string `json:"pushedTags"`
}

func (this *Event) FormattedTime() string {
	return utils.UnixTsFormat(this.EventTime)
}

func (this *Event) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Status:%s, Strategy:%v, Expression:%v, EExp:%v LeftValue:%s, CurrentStep:%d, PushedTags:%v, TS:%s>",
		this.Endpoint,
		this.Status,
		this.Strategy,
		this.Expression,
		this.EExp,
		utils.ReadableFloat(this.LeftValue),
		this.CurrentStep,
		this.PushedTags,
		this.FormattedTime(),
	)
}

func (this *Event) ExpressionId() int {
	if this.Expression != nil {
		return this.Expression.Id
	}

	return 0
}

func (this *Event) EExpID() int {
	if this.Expression != nil {
		return this.EExp.ID
	}
	return 0
}

func (this *Event) StrategyId() int {
	if this.Strategy != nil {
		return this.Strategy.Id
	}

	return 0
}

func (this *Event) TplId() int {
	if this.Strategy != nil {
		return this.Strategy.Tpl.Id
	}

	return 0
}

func (this *Event) Tpl() *Template {
	if this.Strategy != nil {
		return this.Strategy.Tpl
	}

	return nil
}

func (this *Event) ActionId() int {
	if this.Expression != nil {
		return this.Expression.ActionId
	}

	if this.Strategy != nil {
		return this.Strategy.Tpl.ActionId
	}

	return -1

}

func (this *Event) Priority() int {
	if this.Strategy != nil {
		return this.Strategy.Priority
	}
	if this.Expression != nil {
		return this.Expression.Priority
	}

	if this.EExp != nil {
		return this.EExp.Priority
	}
	return -1
}

func (this *Event) Note() string {
	if this.Strategy != nil {
		return this.Strategy.Note
	}

	if this.Expression != nil {
		return this.Expression.Note

	}
	if this.EExp != nil {
		return this.EExp.Note
	}
	return ""
}

func (this *Event) Metric() string {
	if this.Strategy != nil {
		return this.Strategy.Metric
	}
	if this.Expression != nil {
		return this.Expression.Metric
	}
	if this.EExp != nil {
		return this.EExp.Metric
	}
	return ""
}

func (this *Event) RightValue() float64 {
	if this.Strategy != nil {
		return this.Strategy.RightValue
	}

	if this.Expression != nil {
		return this.Expression.RightValue
	}
	return 0.0
}

func (this *Event) Operator() string {
	if this.Strategy != nil {
		return this.Strategy.Operator
	}

	if this.Expression != nil {
		return this.Expression.Operator
	}
	return ""
}

func (this *Event) Func() string {
	if this.Strategy != nil {
		return this.Strategy.Func
	}

	if this.Expression != nil {
		return this.Expression.Func
	}

	if this.EExp != nil {
		return this.EExp.Func
	}
	return ""
}

func (this *Event) MaxStep() int {
	if this.Strategy != nil {
		return this.Strategy.MaxStep
	}

	if this.Expression != nil {
		return this.Expression.MaxStep
	}
	if this.EExp != nil {
		return this.EExp.MaxStep
	}
	return 1
}

func (this *Event) Counter() string {
	return fmt.Sprintf("%s/%s %s", this.Endpoint, this.Metric(), utils.SortedTags(this.PushedTags))
}
