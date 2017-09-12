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
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/toolkits/nux"
	"net/http"
	"runtime"
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
		if !funcs.CpuPrepared() {
			RenderMsgJson(w, "not prepared")
			return
		}

		idle := funcs.CpuIdle()
		busy := 100.0 - idle

		item := [10]string{
			fmt.Sprintf("%.1f%%", idle),
			fmt.Sprintf("%.1f%%", busy),
			fmt.Sprintf("%.1f%%", funcs.CpuUser()),
			fmt.Sprintf("%.1f%%", funcs.CpuNice()),
			fmt.Sprintf("%.1f%%", funcs.CpuSystem()),
			fmt.Sprintf("%.1f%%", funcs.CpuIowait()),
			fmt.Sprintf("%.1f%%", funcs.CpuIrq()),
			fmt.Sprintf("%.1f%%", funcs.CpuSoftIrq()),
			fmt.Sprintf("%.1f%%", funcs.CpuSteal()),
			fmt.Sprintf("%.1f%%", funcs.CpuGuest()),
		}

		RenderDataJson(w, [][10]string{item})
	})

	http.HandleFunc("/proc/cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		if !funcs.CpuPrepared() {
			RenderMsgJson(w, "not prepared")
			return
		}

		idle := funcs.CpuIdle()
		busy := 100.0 - idle

		RenderDataJson(w, map[string]interface{}{
			"idle":    idle,
			"busy":    busy,
			"user":    funcs.CpuUser(),
			"nice":    funcs.CpuNice(),
			"system":  funcs.CpuSystem(),
			"iowait":  funcs.CpuIowait(),
			"irq":     funcs.CpuIrq(),
			"softirq": funcs.CpuSoftIrq(),
			"steal":   funcs.CpuSteal(),
			"guest":   funcs.CpuGuest(),
		})
	})
}
