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
	"net/http"
	"runtime"

	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/toolkits/nux"
)

func configCpuRoutes() {
	http.HandleFunc("/proc/cpu/num", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, runtime.NumCPU())
	})

	http.HandleFunc("/proc/cpu/mhz", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.CpuMHz()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/page/cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		cpuUsages, _, prepared := funcs.CpuUsagesSummary()

		if !prepared {
			RenderMsgJson(w, "not prepared")
			return
		}

		item := [10]string{
			fmt.Sprintf("%.1f%%", cpuUsages[0]),
			fmt.Sprintf("%.1f%%", cpuUsages[1]),
			fmt.Sprintf("%.1f%%", cpuUsages[2]),
			fmt.Sprintf("%.1f%%", cpuUsages[3]),
			fmt.Sprintf("%.1f%%", cpuUsages[4]),
			fmt.Sprintf("%.1f%%", cpuUsages[5]),
			fmt.Sprintf("%.1f%%", cpuUsages[6]),
			fmt.Sprintf("%.1f%%", cpuUsages[7]),
			fmt.Sprintf("%.1f%%", cpuUsages[8]),
			fmt.Sprintf("%.1f%%", cpuUsages[9]),
		}

		RenderDataJson(w, [][10]string{item})
	})

	http.HandleFunc("/proc/cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		cpuUsages, _, prepared := funcs.CpuUsagesSummary()

		if !prepared {
			RenderMsgJson(w, "not prepared")
			return
		}

		RenderDataJson(w, map[string]interface{}{
			"idle":    cpuUsages[0],
			"busy":    cpuUsages[1],
			"user":    cpuUsages[2],
			"nice":    cpuUsages[3],
			"system":  cpuUsages[4],
			"iowait":  cpuUsages[5],
			"irq":     cpuUsages[6],
			"softirq": cpuUsages[7],
			"steal":   cpuUsages[8],
			"guest":   cpuUsages[9],
		})
	})
}
