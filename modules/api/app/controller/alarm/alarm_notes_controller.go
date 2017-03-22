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
	StartTime int64 `json:"startTime" form:"startTime" binding:"required"`
	EndTime   int64 `json:"endTime" form:"endTime" binding:"required"`
	//id
	Id     string `json:"id" form:"id"`
	Status string `json:"status" form:"status"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
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
	if s.Id != "" {
		tmp = append(tmp, fmt.Sprintf("id = '%s'", s.Id))
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
		h.JSONR(c, badstatus, err)
		return
	}
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.EventNote{}
	notes := []alm.EventCases{}
	if inputs.Limit == 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	perparedSql := fmt.Sprintf("select * from %s %s limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
	db.Alarm.Raw(perparedSql).Find(&notes)
	h.JSONR(c, notes)
}

type APIAddNotesToAlarmInputs struct {
	EvnetCaseId     string `json:"evnet_caseId" form:"id" binding:"required"`
	Note string `json:"note" form:"note" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
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
	default:
		return errors.New(`params status: only accepect ["in progress", "unresolved", "resolved", "ignored"]`)
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
	}
	dt := db.Alarm.Begin()
	if err := dt.Save(&Anote); err.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err.Error)
		return
	}
	if inputs.Status == "in progress" || inputs.Status == "resolved" {
		ecase := alm.EventCases{
			ProcessNote: Anote.ID,
			ProcessStatus: Anote.Status,
		}
		if err := dt.Update(&ecase); err != nil {
			dt.Rollback()
			h.JSONR(c, badstatus, err.Error)
			return
		}
	}
	 h.JSONR(c, fmt.Sprintf("add note to %s successfuled", inputs.EvnetCaseId))
	 return
}
