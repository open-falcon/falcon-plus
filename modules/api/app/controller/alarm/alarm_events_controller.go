package alarm

import (
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	alm "github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
	"strings"
	"fmt"
)

// {
// 	startTime:,
// 	endTime:,
// 	priority: 0,
// 	status: "OK,PROBLEM",
// 	process_status: "unresolved",
// 	metrics: "cpu.idle",
// 	id: "",
// 	limit: 1000,
// 	page: 1,
// }
type APIGetAlarmListsInputs struct {
	StartTime int64 `json:"startTime" form:"startTime" binding:"required"`
	EndTime int64 `json:"endTime" form:"endTime" binding:"required"`
	Priority int `json:"priority" form:"priority"`
	Status string `json:"status" form:"status"`
	ProcessStatus string `json:"process_status" form:"process_status"`
	Metrics string `json:"metrics" form:"metrics"`
	//id
	Id string `json:"id" form:"id"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page"`
}

func (s APIGetAlarmListsInputs) collectFilters() string{
	tmp := []string{}
	if s.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
	}
	if s.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
	}
	if s.Priority != -1 {
		tmp = append(tmp, fmt.Sprintf("priority = %d", s.Priority))
	}
	if s.Status != "" {
		status := ""
		statusTmp := strings.Split(s.Status, ",")
		for indx, n := range statusTmp {
			if indx == 0 {
				status = fmt.Sprintf(" status = '%s' ", n)
			}else{
				status = fmt.Sprintf(" %s OR status = '%s' ",status, n)
			}
		}
		tmp = append(tmp, status)
	}
	if s.ProcessStatus != ""{
		pstatus := ""
		pstatusTmp := strings.Split(s.ProcessStatus, ",")
		for indx, n := range pstatusTmp {
			if indx == 0 {
				pstatus = fmt.Sprintf(" process_status = '%s' ", n)
			}else{
				pstatus = fmt.Sprintf(" %s OR process_status = '%s' ",pstatus, n)
			}
		}
		tmp = append(tmp, pstatus)
	}
	if s.Metrics != ""{
		tmp = append(tmp, fmt.Sprintf("metrics regexp '%s'", s.Metrics))
	}
	if s.Id != ""{
		tmp = append(tmp, fmt.Sprintf("id = '%s'", s.Id))
	}
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func AlarmLists(c *gin.Context) {
	var inputs APIGetAlarmListsInputs
	//set default
	inputs.Page = -1
	inputs.Priority = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	fmt.Printf("debug-inputs: %v", inputs)
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.EventCases{}
	cevens := []alm.EventCases{}
	perparedSql := ""
	//if no specific, will give return first 2000 records
	if inputs.Page == -1 {
		if inputs.Limit >= 2000 || inputs.Limit == 0 {
			inputs.Limit = 2000
		}
		perparedSql = fmt.Sprintf("select * from %s %s limit %d", f.TableName(), filterCollector, inputs.Limit)
	}else {
		//set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		perparedSql = fmt.Sprintf("select * from %s %s limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
	}
	db.Alarm.Raw(perparedSql).Find(&cevens)
  h.JSONR(c, cevens)
}


type APIEventsGetInputs struct {
	StartTime int64 `json:"startTime" form:"startTime" binding:"required"`
	EndTime int64 `json:"endTime" form:"endTime" binding:"required"`
	Status int `json:"status" form:"status"`
	//id
	Id string `json:"id" form:"id"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (s APIEventsGetInputs) collectFilters() string{
	tmp := []string{}
	if s.StartTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
	}
	if s.EndTime != 0 {
		tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
	}
	if s.Id != ""{
		tmp = append(tmp, fmt.Sprintf("id = '%s'", s.Id))
	}
	tmp = append(tmp, fmt.Sprintf("status = %d", s.Status))
	filterStrTmp := strings.Join(tmp, " AND ")
	if filterStrTmp != "" {
		filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
	}
	return filterStrTmp
}

func EventsGet(c *gin.Context){
	var inputs APIEventsGetInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	filterCollector := inputs.collectFilters()
	//for get correct table name
	f := alm.Events{}
	evens := []alm.Events{}
	if inputs.Limit == 0 || inputs.Limit >= 50 {
		inputs.Limit = 50
	}
	perparedSql := fmt.Sprintf("select * from %s %s limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
	db.Alarm.Raw(perparedSql).Find(&evens)
  h.JSONR(c, evens)
}
