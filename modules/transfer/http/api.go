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
