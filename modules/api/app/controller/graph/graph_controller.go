package graph

import (
	"fmt"
	"strconv"
	"strings"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	g "github.com/open-falcon/falcon-plus/modules/api/graph"
)

func EndpointRegexpQuery(c *gin.Context) {
	q := c.DefaultQuery("q", "")
	label := c.DefaultQuery("tags", "")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}

	if q == "" && label == "" {
		h.JSONR(c, http.StatusBadRequest, "q and labels are all missing")
		return
	}

	labels := []string{}
	if label != "" {
		labels = strings.Split(label, ",")
	}
	qs := []string{}
	if q != "" {
		qs = strings.Split(q, " ")
	}

	var endpoint []m.Endpoint
	if len(qs) > 0 && len(labels) > 0 {
		var endpoint_id []int
		dt := db.Graph.Table("endpoint_counter").Select("distinct endpoint_id").Where("counter like ?", "%"+strings.TrimSpace(labels[0])+"%")
		for _, term := range labels[1:] {
			dt = dt.Where("counter like ?", "%"+strings.TrimSpace(term)+"%")
		}
		dt = dt.Limit(500).Pluck("distinct endpoint_id", &endpoint_id)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}

		dt = db.Graph.Table("endpoint").Select("endpoint, id").Where("id in (?)", endpoint_id).Where("endpoint regexp ?", strings.TrimSpace(qs[0]))
		for _, term := range qs[1:] {
			dt = dt.Where("endpoint regexp ?", strings.TrimSpace(term))
		}
		dt = dt.Limit(limit).Scan(&endpoint)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
	} else if len(qs) > 0 {
		dt := db.Graph.Table("endpoint").Select("endpoint, id").Where("endpoint regexp ?", strings.TrimSpace(qs[0]))
		for _, term := range qs[1:] {
			dt = dt.Where("endpoint regexp ?", strings.TrimSpace(term))
		}
		dt = dt.Limit(limit).Scan(&endpoint)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
	} else if len(labels) > 0 {
		var endpoint_id []int
		dt := db.Graph.Table("endpoint_counter").Select("distinct endpoint_id").Where("counter like ?", "%"+strings.TrimSpace(labels[0])+"%")
		for _, term := range labels[1:] {
			dt = dt.Where("counter like ?", "%"+strings.TrimSpace(term)+"%")
		}
		dt = dt.Limit(500).Pluck("distinct endpoint_id", &endpoint_id)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}

		dt = db.Graph.Table("endpoint").Select("endpoint, id").Where("id in (?)", endpoint_id)
		dt = dt.Limit(limit).Scan(&endpoint)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
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
		dt = dt.Limit(limit).Scan(&counters)
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

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cmodel.GraphQueryResponse, err error) {
	qparm := g.GenQParam(hostname, counter, consolFun, startTime, endTime, step)
	log.Debugf("qparm: %v", qparm)
	resp, err = g.QueryOne(qparm)
	if err != nil {
		log.Debugf("query graph got error: %s", err.Error())
	}
	return
}
