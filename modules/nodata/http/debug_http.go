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
	"net/http"
	"time"

	"github.com/open-falcon/falcon-plus/modules/nodata/collector"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/sender"
)

func configDebugHttpRoutes() {
	http.HandleFunc("/debug/collector/collect", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := collector.CollectDataOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		RenderDataJson(w, ret)
	})

	http.HandleFunc("/debug/config/sync", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := config.SyncNdConfigOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		RenderDataJson(w, ret)
	})

	http.HandleFunc("/debug/sender/send", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := sender.SendMockOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		RenderDataJson(w, ret)
	})
}
