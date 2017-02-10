package dashboard_graph

import (
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/dashboard"
	"sort"
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

func DashboardGraphCreate(c *gin.Context) {
}

func DashboardScreenCreate(c *gin.Context) {
}
