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

package event

import (
	"time"
)

type EventCases struct {
	// uniuq
	Id       string `json:"id" orm:"pk"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Func     string `json:"func"`
	//leftValue + operator + rightValue
	Cond          string    `json:"cond"`
	Note          string    `json:"note"`
	MaxStep       int       `json:"max_step"`
	CurrentStep   int       `json:"current_step"`
	Priority      int       `json:"priority"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"start_at"`
	UpdateAt      time.Time `json:"update_at"`
	ProcessNote   int       `json:"process_note"`
	ProcessStatus string    `json:"process_status"`
	TplCreator    string    `json:"tpl_creator"`
	ExpressionId  int       `json:"expression_id"`
	StrategyId    int       `json:"strategy_id"`
	TemplateId    int       `json:"template_id"`
	Events        []*Events `json:"evevnts" orm:"reverse(many)"`
}

type Events struct {
	Id          int         `json:"id" orm:"pk"`
	Step        int         `json:"step"`
	Cond        string      `json:"cond"`
	Status      int         `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	EventCaseId *EventCases `json:"event_caseId" orm:"rel(fk)"`
}
