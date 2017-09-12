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
	"fmt"
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
	"strings"
	"time"
)

type APIGetNotesOfAlarmInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	//id
	EventId string `json:"event_id" form:"event_id"`
	Status  string `json:"status" form:"status"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (input APIGetNotesOfAlarmInputs) checkInputsContain() error {
	if input.StartTime == 0 && input.EndTime == 0 {
		if input.EventId == "" {
			return errors.New("startTime, endTime OR event_id, You have to at least pick one on the request.")
		}
	}
	return nil
}

func (s APIGetNotesOfAlarmInputs) collectFilters() string {
	tmp := []string{}
	if s.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
	}
	if s.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
	}
	if s.Status != "" {
		tmp = append(tmp, fmt.Sprintf("status = '%s'", s.Status))
	}
	if s.EventId != "" {
		tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", s.EventId))
	}
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

type APIGetNotesOfAlarmOuput struct {
	EventCaseId string     `json:"event_caseId"`
	Note        string     `json:"note"`
	CaseId      string     `json:"case_id"`
	Status      string     `json:"status"`
	Timestamp   *time.Time `json:"timestamp"`
	UserName    string     `json:"user"`
}

func GetNotesOfAlarm(c *gin.Context) {
	var inputs APIGetNotesOfAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "binding input got error: "+err.Error())
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.EventNote{}
	notes := []alm.EventNote{}
	if inputs.Limit == 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	perparedSql := fmt.Sprintf(
		"select id, event_caseId, note, case_id, status, timestamp, user_id from %s %s order by timestamp DESC limit %d,%d",
		f.TableName(),
		filterCollector,
		inputs.Page,
		inputs.Limit,
	)
	db.Alarm.Raw(perparedSql).Scan(&notes)
	output := []APIGetNotesOfAlarmOuput{}
	for _, n := range notes {
		output = append(output, APIGetNotesOfAlarmOuput{
			EventCaseId: n.EventCaseId,
			Note:        n.Note,
			CaseId:      n.CaseId,
			Status:      n.Status,
			Timestamp:   n.Timestamp,
			UserName:    n.GetUserName(),
		})
	}
	h.JSONR(c, output)
}

type APIAddNotesToAlarmInputs struct {
	EventId string `json:"event_id" form:"event_id" binding:"required"`
	Note    string `json:"note" form:"note" binding:"required"`
	Status  string `json:"status" form:"status" binding:"required"`
	CaseId  string `json:"case_id" form:"case_id"`
}

func (s APIAddNotesToAlarmInputs) CheckingFormating() error {
	switch s.Status {
	case "in progress":
		return nil
	case "unresolved":
		return nil
	case "resolved":
		return nil
	case "ignored":
		return nil
	case "comment":
		return nil
	default:
		return errors.New(`params status: only accepect ["in progress", "unresolved", "resolved", "ignored", "comment"]`)
	}
}

func AddNotesToAlarm(c *gin.Context) {
	var inputs APIAddNotesToAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckingFormating(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	Anote := alm.EventNote{
		UserId:      user.ID,
		Note:        inputs.Note,
		Status:      inputs.Status,
		EventCaseId: inputs.EventId,
		CaseId:      inputs.CaseId,
		//time will update on database self
	}
	dt := db.Alarm.Begin()
	if err := dt.Save(&Anote); err.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err.Error)
		return
	}
	if inputs.Status != "comment" {
		ecase := alm.EventCases{
			ProcessNote:   Anote.ID,
			ProcessStatus: Anote.Status,
		}
		if db := dt.Table(ecase.TableName()).Where("id = ?", Anote.EventCaseId).Update(&ecase); db.Error != nil {
			dt.Rollback()
			h.JSONR(c, badstatus, "update got error during update event_cases:"+db.Error.Error())
			return
		}
	}
	dt.Commit()
	h.JSONR(c, map[string]string{
		"id":      inputs.EventId,
		"message": fmt.Sprintf("add note to %s successfuled", inputs.EventId),
	})
	return
}
