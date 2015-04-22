package http

import (
	"encoding/json"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"net/http"
)

func configPushRoutes() {
	http.HandleFunc("/v1/push", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		m := map[string]interface{}{
			"err": 0,
			"msg": "",
		}

		if req.ContentLength == 0 {
			m["msg"] = "request.ContentLength == 0: do nothing"
			w.Write(buildJson(m))
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*model.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			m["err"] = 1
			m["msg"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			w.Write(buildJson(m))
			return
		}

		g.SendToTransfer(metrics)

		w.Write(buildJson(m))
	})
}

func buildJson(m map[string]interface{}) []byte {
	bs, err := json.Marshal(m)
	if err != nil {
		m["err"] = 2
		m["msg"] = err.Error()
		bs, _ = json.Marshal(m)
		return bs
	}
	return bs
}
