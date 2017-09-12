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
	"fmt"
	"github.com/toolkits/nux"
	"net/http"
	"runtime"
	"time"
)

func configSystemRoutes() {

	http.HandleFunc("/system/date", func(w http.ResponseWriter, req *http.Request) {
		RenderDataJson(w, time.Now().Format("2006-01-02 15:04:05"))
	})

	http.HandleFunc("/page/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := nux.SystemUptime()
		AutoRender(w, fmt.Sprintf("%d days %d hours %d minutes", days, hours, mins), err)
	})

	http.HandleFunc("/proc/system/uptime", func(w http.ResponseWriter, req *http.Request) {
		days, hours, mins, err := nux.SystemUptime()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		RenderDataJson(w, map[string]interface{}{
			"days":  days,
			"hours": hours,
			"mins":  mins,
		})
	})

	http.HandleFunc("/page/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		cpuNum := runtime.NumCPU()
		load, err := nux.LoadAvg()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		ret := [3][2]interface{}{
			{load.Avg1min, int64(load.Avg1min * 100.0 / float64(cpuNum))},
			{load.Avg5min, int64(load.Avg5min * 100.0 / float64(cpuNum))},
			{load.Avg15min, int64(load.Avg15min * 100.0 / float64(cpuNum))},
		}
		RenderDataJson(w, ret)
	})

	http.HandleFunc("/proc/system/loadavg", func(w http.ResponseWriter, req *http.Request) {
		data, err := nux.LoadAvg()
		AutoRender(w, data, err)
	})

}
