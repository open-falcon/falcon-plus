package alarm

import (
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
	"strings"
	"fmt"
	"errors"
)

type APIGetNotesOfAlarmInputs struct {
	StartTime int64 `json:"startTime" form:"startTime"`
	EndTime   int64 `json:"endTime" form:"endTime"`
	//id
	EventId     string `json:"event_id" form:"event_id"`
	Status string `json:"status" form:"status"`
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

func GetNotesOfAlarm(c *gin.Context) {
	var inputs APIGetNotesOfAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "binding input got error: " + err.Error())
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	fmt.Printf("%v", inputs)
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
	h.JSONR(c, notes)
}

type APIAddNotesToAlarmInputs struct {
	EventId     string `json:"event_id" form:"event_id" binding:"required"`
	Note string `json:"note" form:"note" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	CaseId string `json:"case_id" form:"case_id"`
}

func (s APIAddNotesToAlarmInputs) CheckingFormating() error {
	switch s.Status{
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
		UserId: user.ID,
		Note: inputs.Note,
		Status: inputs.Status,
		EventCaseId: inputs.EventId,
		CaseId: inputs.CaseId,
		//time will update on database self
	}
	fmt.Printf("%v", user)
	dt := db.Alarm.Begin()
	if err := dt.Save(&Anote); err.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err.Error)
		return
	}
	if inputs.Status != "comment" {
		ecase := alm.EventCases{
			ProcessNote: Anote.ID,
			ProcessStatus: Anote.Status,
		}
		if db := dt.Table(ecase.TableName()).Where("id = ?", Anote.EventCaseId).Update(&ecase); db.Error != nil {
			dt.Rollback()
			h.JSONR(c, badstatus, "update got error during update event_cases:" + db.Error.Error())
			return
		}
	}
	dt.Commit()
	h.JSONR(c, map[string]string{
		"id": inputs.EventId,
		"message": fmt.Sprintf("add note to %s successfuled", inputs.EventId),
	})
	return
}
