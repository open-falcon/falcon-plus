package http

import (
	"encoding/json"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/proxy/g"
	prpc "github.com/open-falcon/falcon-plus/modules/proxy/receiver/rpc"
	"net/http"
)

//TODO:delete
func api_query_info(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Query + "/graph/info"
	postByJson(rw, req, url)
}

func api_query_history(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Query + "/graph/history"
	postByJson(rw, req, url)
}

func api_dashboard_endpoints(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + req.URL.RequestURI()
	getRequest(rw, url)
}

func api_dashboard_counters(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + "/api/counters"
	postByForm(rw, req, url)
}

func api_dashboard_chart(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + "/chart"
	postByForm(rw, req, url)
}

func api_push_datapoints(rw http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(rw, "blank body", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var metrics []*cmodel.MetricValue
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}

	reply := &cmodel.TransferResponse{}
	prpc.RecvMetricValues(metrics, reply, "http")

	RenderDataJson(rw, reply)
}

func configApiRoutes() {
	http.HandleFunc("/api/info", api_query_info)
	http.HandleFunc("/api/history", api_query_history)
	http.HandleFunc("/api/endpoints", api_dashboard_endpoints)
	http.HandleFunc("/api/counters", api_dashboard_counters)
	http.HandleFunc("/api/chart", api_dashboard_chart)
	http.HandleFunc("/api/push", api_push_datapoints)
}
