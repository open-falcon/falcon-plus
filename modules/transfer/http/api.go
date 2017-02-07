package http

import (
	"encoding/json"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	prpc "github.com/open-falcon/falcon-plus/modules/transfer/receiver/rpc"
	"net/http"
)

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
	http.HandleFunc("/api/push", api_push_datapoints)
}
