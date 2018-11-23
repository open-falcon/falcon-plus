// Copyright 2018 Xiaomi, Inc.
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

import "github.com/jinzhu/gorm"

type EventCases struct {
	gorm.Model
	EventCaseID string `json:"eventcase_id"`
	Endpoint    string `json:"endpoint"`
	Counter     string `json:"counter"`
	Func        string `json:"func"`
	Cond        string `json:"cond"`
	Note        string `json:"note"`
	MaxStep     int    `json:"step"`
	CurrentStep int    `json:"current_step"`
	Priority    int    `json:"priority"`
	Status      string `json:"status"`
	//记录第一次(CURRENT_STEP=1) PROBLEM时间，OK后再PROBLEM被更新
	EventTs string `json:"event_ts"`
	//记录每次被更新的时间
	UpdateAt       string `json:"update_at"`
	ClosedAt       string `json:"closed_at"`
	ClosedNote     string `json:"closed_note"`
	UserModified   int64  `json:"user_modified"`
	TplCreator     string `json:"tpl_creator"`
	ActionConfigID int64  `json:"action_config_id"`
	ExpressionID   int    `json:"expression_id"`
	StrategyID     int    `json:"strategy_id"`
	TemplateID     int    `json:"template_id"`
	ProcessNote    int64  `json:"process_note"`
	ProcessStatus  string `json:"process_status"`
}

func (EventCases) TableName() string {
	return "event_cases"
}

type Events struct {
	gorm.Model
	EventCaseID string `json:"eventcase_id"`
	CurrentStep int    `json:"current_step"`
	Cond        string `json:"cond"`
	Status      string `json:"status"`
	EventTs     string `json:"event_ts"`
}

func (Events) TableName() string {
	return "events"
}

type EventReceiver struct {
	gorm.Model
	EventID        uint   `json:"event_id"`
	ActionConfigID int64  `json:"action_config_id"`
	Uic            string `json:"uic"`
	User           string `json:"user"`
	EventTs        string `json:"event_ts"`
}

func (EventReceiver) TableName() string {
	return "event_receiver"
}

type EventApiInputs struct {
	UserName  string   `json:"username"`
	Uic       []string `json:"uic"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	Priority  []int    `json:"priority"`
	Status    string   `json:"status"`
	Endpoint  string   `json:"endpoint"`
	Counter   string   `json:"counter"`
	//是否包含故障，包含则传递fault=have
	HaveFault bool `json:"have_fault"`
	//event type, 未恢复报警、历史报警
	NowStatus bool `json:"now_event_status"`
	Limit     int  `json:"limit"`
	Offset    int  `json:"offset"`
}

type EventFaultInfos struct {
	Event  EventInfo    `json:"event"`
	Faults []FaultInfos `json:"faults"`
}

//告警事件信息
type EventInfo struct {
	ID           uint   `json:"event_id"`
	EventCaseID  string `json:"eventcase_id"`
	Endpoint     string `json:"endpoint"`
	Counter      string `json:"counter"`
	Func         string `json:"func"`
	Cond         string `json:"cond"`
	Note         string `json:"note"`
	MaxStep      int    `json:"max_step"`
	CurrentStep  int    `json:"current_step"`
	Priority     int    `json:"priority"`
	Status       string `json:"status"`
	EventTs      string `json:"event_ts"`
	TplCreator   string `json:"template_creator"`
	ExpressionID int    `json:"expression_id"`
	StrategyID   int    `json:"strategy_id"`
	TemplateID   int    `json:"template_id"`
}

type FaultInfos struct {
	FaultID   uint   `json:"fault_id"`
	UpdatedAt string `json:"updated_at"`
	Title     string `json:"title"`
}

//返回故障数量
type EventCount struct {
	Count int `json:"count"`
}
