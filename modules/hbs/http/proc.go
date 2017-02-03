package http

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
	"net/http"
)

func configProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.Agents.Keys())
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Host, len(cache.MonitoredHosts.Get()))
		for k, v := range cache.MonitoredHosts.Get() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Strategy, len(cache.Strategies.GetMap()))
		for k, v := range cache.Strategies.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Template, len(cache.TemplateCache.GetMap()))
		for k, v := range cache.TemplateCache.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		RenderDataJson(w, cache.GetPlugins(hostname))
	})

}
