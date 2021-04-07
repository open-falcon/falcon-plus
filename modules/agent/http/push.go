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
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

var once sync.Once

func configPushRoutes() {
	http.HandleFunc("/v1/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}
		if g.Config().AgentMemCtrl == true {
			ok := isHandReq(w)
			if !ok {
				return
			}
		}
		decoder := json.NewDecoder(req.Body)
		var metrics []*model.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "cannot decode body", http.StatusBadRequest)
			return
		}
		if len(metrics) > g.Config().Batch {
			g.SendToTransfer(metrics[:g.Config().Batch])
			http.Error(w, "post Metric too Big !!! have sent max Batch: "+strconv.Itoa(g.Config().Batch), http.StatusBadRequest)
			return
		}
		g.SendToTransfer(metrics)
		w.Write([]byte("success"))
	})
}

func isHandReq(w http.ResponseWriter) bool {
	once.Do(funcs.InitCgroup)
	memUsed, err := funcs.GetAgentMem()
	if err != nil {
		RenderMsgJson(w, err.Error())
		return false
	}
	if uint64(memUsed) > g.Config().AgentMemLimit {
		log.Printf("memory consumption has exceeded the threshold")
		http.Error(w, "memory consumption has exceeded the threshold", http.StatusBadRequest)
		return false
	}
	return true
}
