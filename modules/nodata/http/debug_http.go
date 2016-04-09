package http

import (
	"net/http"
	"time"

	"github.com/open-falcon/nodata/collector"
	"github.com/open-falcon/nodata/config"
	"github.com/open-falcon/nodata/sender"
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

	http.HandleFunc("/debug/sender/calcGauss", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Unix()
		cnt := sender.CalcGaussOnce()
		end := time.Now().Unix()

		ret := make(map[string]int, 0)
		ret["cnt"] = cnt
		ret["time"] = int(end - start)
		RenderDataJson(w, ret)
	})
}
