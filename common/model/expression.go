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

type Expression struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	ActionId   int               `json:"actionId"`
	CreateUser string            `json:"createUser"`
	Pause      int               `json:"pause"`
}

func (this *Expression) String() string {
	return fmt.Sprintf(
		"<Id:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d %s ActionId:%d CreateUser:%s Pause:%d>",
		this.Id,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
		this.ActionId,
		this.CreateUser,
		this.Pause,
	)
}

type ExpressionResponse struct {
	Expressions []*Expression `json:"expressions"`
}
