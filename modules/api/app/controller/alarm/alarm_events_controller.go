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

package alarm

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
	"strings"
)

type APIGetAlarmListsInputs struct {
	StartTime     int64  `json:"startTime" form:"startTime"`
	EndTime       int64  `json:"endTime" form:"endTime"`
	Priority      int    `json:"priority" form:"priority"`
	Status        string `json:"status" form:"status"`
	ProcessStatus string `json:"process_status" form:"process_status"`
	Metrics       string `json:"metrics" form:"metrics"`
	//id
	EventId string `json:"event_id" form:"event_id"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
	//endpoints strategy template
	Endpoints  []string `json:"endpoints" form:"endpoints"`
	StrategyId int      `json:"strategy_id" form:"strategy_id"`
	TemplateId int      `json:"template_id" form:"template_id"`
}

func (input APIGetAlarmListsInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventId == "" && input.Endpoints == nil && input.StrategyId == 0 && input.TemplateId == 0 {
			return errors.New("startTime, endTime, event_id, endpoints, strategy_id or template_id, You have to at least pick one on the request.")
		}
	}
	return nil
}

func (s APIGetAlarmListsInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if s.StartTime != 0 {
		filterDB = filterDB.Where("timestamp >= FROM_UNIXTIME(?)", s.StartTime)
	}
	if s.EndTime != 0 {
		filterDB = filterDB.Where("timestamp <= FROM_UNIXTIME(?)", s.EndTime)
	}
	if s.Priority != -1 {
		filterDB = filterDB.Where("priority = ?", s.Priority)
	}
	if s.Status != "" {
		statusTmp := strings.Split(s.Status, ",")
		filterDB = filterDB.Where("status in (?)", statusTmp)
	}
	if s.ProcessStatus != "" {
		pstatusTmp := strings.Split(s.ProcessStatus, ",")
		filterDB = filterDB.Where("process_status in (?)", pstatusTmp)
	}
	if s.Metrics != "" {
		filterDB = filterDB.Where("metric regexp ?", s.Metrics)
	}
	if s.EventId != "" {
		filterDB = filterDB.Where("id = ?", s.EventId)
	}
	if s.Endpoints != nil && len(s.Endpoints) != 0 {
		filterDB = filterDB.Where("endpoint in (?)", s.Endpoints)
	}
	if s.StrategyId != 0 {
		filterDB = filterDB.Where("strategy_id = ?", s.StrategyId)
	}
	if s.TemplateId != 0 {
		filterDB = filterDB.Where("template_id = ?", s.TemplateId)
	}
	return filterDB
}

func AlarmLists(c *gin.Context) {
	var inputs APIGetAlarmListsInputs
	//set default
	inputs.Page = -1
	inputs.Limit = -1
	inputs.Priority = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//for get correct table name
	f := alm.EventCases{}
	alarmDB := inputs.collectDBFilters(db.Alarm, f.TableName(), nil)
	cevens := []alm.EventCases{}
	//if no specific, will give return first 2000 records
	if inputs.Page == -1 && inputs.Limit == -1{
		inputs.Limit = 2000
		alarmDB = alarmDB.Order("timestamp DESC").Limit(inputs.Limit)
	} else if inputs.Limit == -1 {
		// set page but not set limit
		h.JSONR(c, badstatus, errors.New("You set page but skip limit params, please check your input"))
		return
	} else {
		// set limit but not set page
		if inputs.Page == -1 {
			// limit invalid
			if inputs.Limit <= 0 {
				h.JSONR(c, badstatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
			// set default page
			inputs.Page = 1
		} else {
			// set page and limit
			// page or limit invalid
			if inputs.Page <= 0 || inputs.Limit <= 0 {
				h.JSONR(c, badstatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
		}
		//set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		step := (inputs.Page -1) * inputs.Limit
		alarmDB = alarmDB.Order("timestamp DESC").Offset(step).Limit(inputs.Limit)
	}
	alarmDB.Find(&cevens)
	h.JSONR(c, cevens)
}

type APIEventsGetInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	Status    int   `json:"status" form:"status" binding:"gte=-1,lte=1"`
	//event_caseId
	EventId string `json:"event_id" form:"event_id" binding:"required"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (s APIEventsGetInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if s.StartTime != 0 {
		filterDB = filterDB.Where("timestamp >= FROM_UNIXTIME(?)", s.StartTime)
	}
	if s.EndTime != 0 {
		filterDB = filterDB.Where("timestamp <= FROM_UNIXTIME(?)", s.EndTime)
	}
	if s.EventId != "" {
		filterDB = filterDB.Where("event_caseId = ?", s.EventId)
	}
	if s.Status == 0 || s.Status == 1 {
		filterDB = filterDB.Where("status = ?", s.Status)
	}
	return filterDB
}

func EventsGet(c *gin.Context) {
	var inputs APIEventsGetInputs
	inputs.Status = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//for get correct table name
	f := alm.Events{}
	eventDB := inputs.collectDBFilters(db.Alarm, f.TableName(), []string{"id", "step", "event_caseId", "cond", "status", "timestamp"})
	evens := []alm.Events{}
	if inputs.Limit <= 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	step := (inputs.Page -1) * inputs.Limit
	eventDB.Order("timestamp DESC").Offset(step).Limit(inputs.Limit).Scan(&evens)
	h.JSONR(c, evens)
}
