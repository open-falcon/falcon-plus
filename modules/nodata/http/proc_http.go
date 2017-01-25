package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/nodata/collector"
	"github.com/open-falcon/falcon-plus/modules/nodata/config"
	"github.com/open-falcon/falcon-plus/modules/nodata/config/service"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
	"github.com/open-falcon/falcon-plus/modules/nodata/judge"
)

func configProcHttpRoutes() {
	// counters
	http.HandleFunc("/proc/counters", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetAllCounters())
	})
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetAllCounters())
	})

	// judge.status, /proc/status/$endpoint/$metric/$tags-pairs
	http.HandleFunc("/proc/status/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/proc/status/"):]
		RenderDataJson(w, judge.GetNodataStatus(urlParam))
	})

	// collector.last.item, /proc/collect/$endpoint/$metric/$tags-pairs
	http.HandleFunc("/proc/collect/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/proc/collect/"):]
		item, _ := collector.GetFirstItem(urlParam)
		RenderDataJson(w, item.String())
	})

	// config.mockcfg
	http.HandleFunc("/proc/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, service.GetMockCfgFromDB())
	})
	// config.mockcfg /proc/config/$endpoint/$metric/$tags-pairs
	http.HandleFunc("/proc/config/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/proc/config/"):]
		cfg, _ := config.GetNdConfig(urlParam)
		RenderDataJson(w, cfg)
	})

	// config.hostgroup, /group/$grpname
	http.HandleFunc("/proc/group/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/proc/group/"):]
		RenderDataJson(w, service.GetHostsFromGroup(urlParam))
	})
}
