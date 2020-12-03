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
	"errors"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"log"
	"net/http"
	"strconv"
)

func configPushRoutes() {
	http.HandleFunc("/v1/push", cors(func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}
		if g.Config().MemoryCtrl == true {
			errMem := MemCtrl(w)
			if errMem != nil {
				return
			}
		}
		decoder := json.NewDecoder(req.Body)
		var metrics []*model.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}
		if len(metrics) >= g.Config().Batch {
			http.Error(w, "cannot post Metric too Big !!! curr count: "+strconv.Itoa(len(metrics)), http.StatusBadRequest)
		}
		g.SendToTransfer(metrics)
		w.Write([]byte("success"))
	}))
}

func MemCtrl(w http.ResponseWriter) error {
	log.Printf("memory control start")
	mem, err := funcs.AgentMemInfo()
	if err != nil {
		RenderMsgJson(w, err.Error())
		return err
	}
	memUsed := mem.VmRSS
	if memUsed > g.Config().MaxMemory {
		log.Printf("memory consumption has exceeded the threshold")
		http.Error(w, "memory consumption has exceeded the threshold", http.StatusBadRequest)
		return errors.New("memory consumption has exceeded the threshold")
	}
	return err
}

// 确认下相关请求是否有关联
func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		switch r.Method {
		case "OPTIONS":
			method := r.Header.Get("Access-Control-Request-Method")
			log.Printf("preflight request method: %v", method)
			if method == "POST" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "POST":
			if origin == r.Header.Get("Origin") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if origin == "" {
				f(w, r)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}
