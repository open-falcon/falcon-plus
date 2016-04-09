package http

import (
	"fmt"
	"github.com/open-falcon/agent/funcs"
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
