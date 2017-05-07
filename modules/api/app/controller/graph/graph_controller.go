package graph

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	grh "github.com/open-falcon/falcon-plus/modules/api/graph"
	"net/http"
)

type APIEndpointRegexpQueryInputs struct {
	Q     string `json:"q" form:"q"`
	Label string `json:"tags" form:"tags"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

func EndpointRegexpQuery(c *gin.Context) {
	inputs := APIEndpointRegexpQueryInputs{
		//set default is 500
		Limit: 500,
		Page:  1,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Q == "" && inputs.Label == "" {
		h.JSONR(c, http.StatusBadRequest, "q and labels are all missing")
		return
	}

	labels := []string{}
	if inputs.Label != "" {
		labels = strings.Split(inputs.Label, ",")
	}
	qs := []string{}
	if inputs.Q != "" {
		qs = strings.Split(inputs.Q, " ")
	}

	var offset int = 0
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}

	var endpoint []m.Endpoint
	var endpoint_id []int
	var dt *gorm.DB
	if len(labels) != 0 {
		dt = db.Graph.Table("endpoint_counter").Select("distinct endpoint_id")
		for _, trem := range labels {
			dt = dt.Where(" counter like ? ", "%"+strings.TrimSpace(trem)+"%")
		}
		dt = dt.Limit(inputs.Limit).Offset(offset).Pluck("distinct endpoint_id", &endpoint_id)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
	}
	if len(qs) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id")
		if len(endpoint_id) != 0 {
			dt = dt.Where("id in (?)", endpoint_id)
		}

		for _, trem := range qs {
			dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(trem))
		}
		dt.Limit(inputs.Limit).Offset(offset).Scan(&endpoint)
	} else if len(endpoint_id) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id").
			Where("id in (?)", endpoint_id).
			Scan(&endpoint)
	}
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, e := range endpoint {
		endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint})
	}

	h.JSONR(c, endpoints)
}

func EndpointCounterRegexpQuery(c *gin.Context) {
	eid := c.DefaultQuery("eid", "")
	metricQuery := c.DefaultQuery("metricQuery", ".+")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	pageTmp := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	var offset int = 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	if eid == "" {
		h.JSONR(c, http.StatusBadRequest, "eid is missing")
	} else {
		eids := utils.ConverIntStringToList(eid)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
		}

		var counters []m.EndpointCounter
		dt := db.Graph.Table("endpoint_counter").Select("counter, step, type").Where(fmt.Sprintf("endpoint_id IN %s", eids))
		if metricQuery != "" {
			qs := strings.Split(metricQuery, " ")
			if len(qs) > 0 {
				for _, term := range qs {
					dt = dt.Where("counter regexp ?", strings.TrimSpace(term))
				}
			}
		}
		dt = dt.Limit(limit).Offset(offset).Scan(&counters)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}

		countersResp := []interface{}{}
		for _, c := range counters {
			countersResp = append(countersResp, map[string]interface{}{
				"counter": c.Counter,
				"step":    c.Step,
				"type":    c.Type,
			})
		}
		h.JSONR(c, countersResp)
	}
	return
}

type APIQueryGraphDrawData struct {
	HostNames []string `json:"hostnames" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
	ConsolFun string   `json:"consol_fun" binding:"required"`
	StartTime int64    `json:"start_time" binding:"required"`
	EndTime   int64    `json:"end_time" binding:"required"`
	Step      int      `json:"step"`
}

func QueryGraphDrawData(c *gin.Context) {
	var inputs APIQueryGraphDrawData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	respData := []*cmodel.GraphQueryResponse{}
	for _, host := range inputs.HostNames {
		for _, counter := range inputs.Counters {
			// TODO:cache step
			var step []int
			dt := db.Graph.Raw("select a.step from endpoint_counter as a, endpoint as b where b.endpoint = ? and a.endpoint_id = b.id and a.counter = ? limit 1", host, counter).Scan(&step)
			if dt.Error != nil || len(step) == 0 {
				continue
			}
			data, _ := fetchData(host, counter, inputs.ConsolFun, inputs.StartTime, inputs.EndTime, step[0])
			respData = append(respData, data)
		}
	}
	h.JSONR(c, respData)
}

func QueryGraphLastPoint(c *gin.Context) {
	var inputs []cmodel.GraphLastParam
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	respData := []*cmodel.GraphLastResp{}

	for _, param := range inputs {
		one_resp, err := grh.Last(param)
		if err != nil {
			log.Warn("query last point from graph fail:", err)
		} else {
			respData = append(respData, one_resp)
		}
	}

	h.JSONR(c, respData)
}

func DeleteGraphEndpoint(c *gin.Context) {
	var inputs []string = []string{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterId int
		Counter   string
		Type      string
		Step      int
	}
	rows := []DBRows{}
	dt := db.Graph.Raw(
		`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a, endpoint_counter as b
		where b.endpoint_id = a.id
		AND a.endpoint in (?)`, inputs).Scan(&rows)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var affected_counter int64 = 0
	var affected_endpoint int64 = 0

	if len(rows) > 0 {
		var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
		for _, row := range rows {
			param := &cmodel.GraphDeleteParam{
				Endpoint: row.Endpoint,
				DsType:   row.Type,
				Step:     row.Step,
			}
			fields := strings.SplitN(row.Counter, "/", 2)
			if len(fields) == 1 {
				param.Metric = fields[0]
			} else if len(fields) == 2 {
				param.Metric = fields[0]
				param.Tags = fields[1]
			} else {
				log.Error("invalid counter", row.Counter)
				continue
			}
			params = append(params, param)
		}
		grh.Delete(params)
	}

	tx := db.Graph.Begin()

	if len(rows) > 0 {
		var cids []int = make([]int, len(rows))
		for i, row := range rows {
			cids[i] = row.CounterId
		}

		dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&m.EndpointCounter{})
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			tx.Rollback()
			return
		}
		affected_counter = dt.RowsAffected

		dt = tx.Exec(`delete from tag_endpoint where endpoint_id in 
			(select id from endpoint where endpoint in (?))`, inputs)
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			tx.Rollback()
			return
		}
	}

	dt = tx.Table("endpoint").Where("endpoint in (?)", inputs).Delete(&m.Endpoint{})
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	affected_endpoint = dt.RowsAffected
	tx.Commit()

	h.JSONR(c, map[string]int64{
		"affected_endpoint": affected_endpoint,
		"affected_counter":  affected_counter,
	})
}

type APIGraphDeleteCounterInputs struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
}

func DeleteGraphCounter(c *gin.Context) {
	var inputs APIGraphDeleteCounterInputs = APIGraphDeleteCounterInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterId int
		Counter   string
		Type      string
		Step      int
	}
	rows := []DBRows{}
	dt := db.Graph.Raw(`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a, endpoint_counter as b
		where b.endpoint_id = a.id 
		AND a.endpoint in (?)
		AND b.counter in (?)`, inputs.Endpoints, inputs.Counters).Scan(&rows)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	if len(rows) == 0 {
		h.JSONR(c, map[string]int64{
			"affected_counter": 0,
		})
		return
	}

	var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
	for _, row := range rows {
		param := &cmodel.GraphDeleteParam{
			Endpoint: row.Endpoint,
			DsType:   row.Type,
			Step:     row.Step,
		}
		fields := strings.SplitN(row.Counter, "/", 2)
		if len(fields) == 1 {
			param.Metric = fields[0]
		} else if len(fields) == 2 {
			param.Metric = fields[0]
			param.Tags = fields[1]
		} else {
			log.Error("invalid counter", row.Counter)
			continue
		}
		params = append(params, param)
	}
	grh.Delete(params)

	tx := db.Graph.Begin()
	var cids []int = make([]int, len(rows))
	for i, row := range rows {
		cids[i] = row.CounterId
	}

	dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&m.EndpointCounter{})
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	affected_counter := dt.RowsAffected
	tx.Commit()

	h.JSONR(c, map[string]int64{
		"affected_counter": affected_counter,
	})
}

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cmodel.GraphQueryResponse, err error) {
	qparm := grh.GenQParam(hostname, counter, consolFun, startTime, endTime, step)
	// log.Debugf("qparm: %v", qparm)
	resp, err = grh.QueryOne(qparm)
	if err != nil {
		log.Debugf("query graph got error: %s", err.Error())
	}
	return
}
