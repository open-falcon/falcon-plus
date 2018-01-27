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

package dashboard_graph

import (
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/dashboard"
	"sort"
	"strconv"
	"strings"
	"time"
)

type APITmpGraphCreateReqData struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
}

func DashboardTmpGraphCreate(c *gin.Context) {
	var inputs APITmpGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)

	es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
	ck := cutils.Md5(es_string + ":" + cs_string)

	dt := db.Dashboard.Exec("insert ignore into `tmp_graph` (endpoints, counters, ck) values(?, ?, ?) on duplicate key update time_=?", es_string, cs_string, ck, time.Now())
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	tmp_graph := m.DashboardTmpGraph{}
	dt = db.Dashboard.Table("tmp_graph").Where("ck=?", ck).First(&tmp_graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": int(tmp_graph.ID)})
}

func DashboardTmpGraphQuery(c *gin.Context) {
	id := c.Param("id")

	tmp_graph := m.DashboardTmpGraph{}
	dt := db.Dashboard.Table("tmp_graph").Where("id = ?", id).First(&tmp_graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	es := strings.Split(tmp_graph.Endpoints, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(tmp_graph.Counters, TMP_GRAPH_FILED_DELIMITER)

	ret := map[string][]string{
		"endpoints": es,
		"counters":  cs,
	}

	h.JSONR(c, ret)
}

type APIGraphCreateReqData struct {
	ScreenId   int      `json:"screen_id" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Endpoints  []string `json:"endpoints" binding:"required"`
	Counters   []string `json:"counters" binding:"required"`
	TimeSpan   int      `json:"timespan"`
	GraphType  string   `json:"graph_type"`
	Method     string   `json:"method"`
	Position   int      `json:"position"`
	FalconTags string   `json:"falcon_tags"`
}

func DashboardGraphCreate(c *gin.Context) {
	var inputs APIGraphCreateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)
	es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)

	d := m.DashboardGraph{
		Title:     inputs.Title,
		Hosts:     es_string,
		Counters:  cs_string,
		ScreenId:  int64(inputs.ScreenId),
		TimeSpan:  inputs.TimeSpan,
		GraphType: inputs.GraphType,
		Method:    inputs.Method,
		Position:  inputs.Position,
	}
	if d.TimeSpan == 0 {
		d.TimeSpan = 3600
	}
	if d.GraphType == "" {
		d.GraphType = "h"
	}

	tx := db.Dashboard.Begin()
	dt := tx.Table("dashboard_graph").Create(&d)
	if dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var lid []int
	dt = tx.Table("dashboard_graph").Raw("select LAST_INSERT_ID() as id").Pluck("id", &lid)
	if dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	aid := lid[0]

	h.JSONR(c, map[string]int{"id": aid})

}

type APIGraphUpdateReqData struct {
	ScreenId   int      `json:"screen_id"`
	Title      string   `json:"title"`
	Endpoints  []string `json:"endpoints"`
	Counters   []string `json:"counters"`
	TimeSpan   int      `json:"timespan"`
	GraphType  string   `json:"graph_type"`
	Method     string   `json:"method"`
	Position   int      `json:"position"`
	FalconTags string   `json:"falcon_tags"`
}

func DashboardGraphUpdate(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	var inputs APIGraphUpdateReqData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	d := m.DashboardGraph{}

	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		d.Hosts = es_string
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		d.Counters = cs_string
	}
	if inputs.Title != "" {
		d.Title = inputs.Title
	}
	if inputs.ScreenId != 0 {
		d.ScreenId = int64(inputs.ScreenId)
	}
	if inputs.TimeSpan != 0 {
		d.TimeSpan = inputs.TimeSpan
	}
	if inputs.GraphType != "" {
		d.GraphType = inputs.GraphType
	}
	if inputs.Position != 0 {
		d.Position = inputs.Position
	}
	if inputs.FalconTags != "" {
		d.FalconTags = inputs.FalconTags
	}
       d.Method = inputs.Method

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Model(&graph).Where("id = ?", gid).Updates(d)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

       dt = db.Dashboard.Table("dashboard_graph").Model(&graph).Where("id = ?", gid).Update("method", d.Method)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}


	h.JSONR(c, map[string]int{"id": gid})

}

func DashboardGraphGet(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).First(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

	h.JSONR(c, map[string]interface{}{
		"graph_id":    graph.ID,
		"title":       graph.Title,
		"endpoints":   es,
		"counters":    cs,
		"screen_id":   graph.ScreenId,
		"graph_type":  graph.GraphType,
		"timespan":    graph.TimeSpan,
		"method":      graph.Method,
		"position":    graph.Position,
		"falcon_tags": graph.FalconTags,
	})

}

func DashboardGraphDelete(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).Delete(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": gid})

}

func DashboardGraphGetsByScreenID(c *gin.Context) {
	id := c.Param("screen_id")
	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}
	limit := c.DefaultQuery("limit", "500")

	graphs := []m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("screen_id = ?", sid).Limit(limit).Find(&graphs)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	ret := []map[string]interface{}{}
	for _, graph := range graphs {
		es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

		r := map[string]interface{}{
			"graph_id":    graph.ID,
			"title":       graph.Title,
			"endpoints":   es,
			"counters":    cs,
			"screen_id":   graph.ScreenId,
			"graph_type":  graph.GraphType,
			"timespan":    graph.TimeSpan,
			"method":      graph.Method,
			"position":    graph.Position,
			"falcon_tags": graph.FalconTags,
		}
		ret = append(ret, r)
	}

	h.JSONR(c, ret)
}
