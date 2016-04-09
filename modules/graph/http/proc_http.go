package http

import (
	"net/http"
	"strings"

	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/proc"
	"github.com/open-falcon/graph/store"
)

func configProcRoutes() {
	// counter
	http.HandleFunc("/counter/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// TO BE DISCARDed
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// items.history
	http.HandleFunc("/history/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/history/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		RenderDataJson(w, store.GetAllItems(cutils.Checksum(endpoint, metric, tags)))
	})

	http.HandleFunc("/v2/history", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0) {
			RenderDataJson(w, "bad args")
			return
		}
		endpoint := r.Form["e"][0]
		metric := r.Form["m"][0]

		tags := make(map[string]string)
		if len(r.Form["t"]) > 0 {
			tagstr := r.Form["t"][0]
			tagVals := strings.Split(tagstr, ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}

		RenderDataJson(w, store.GetAllItems(cutils.Checksum(endpoint, metric, tags)))
	})

	// items.last
	http.HandleFunc("/last/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/last/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		RenderDataJson(w, store.GetLastItem(cutils.Checksum(endpoint, metric, tags)))
	})

	http.HandleFunc("/v2/last", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0) {
			RenderDataJson(w, "bad args")
			return
		}
		endpoint := r.Form["e"][0]
		metric := r.Form["m"][0]

		tags := make(map[string]string)
		if len(r.Form["t"]) > 0 {
			tagstr := r.Form["t"][0]
			tagVals := strings.Split(tagstr, ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}

		RenderDataJson(w, store.GetLastItem(cutils.Checksum(endpoint, metric, tags)))
	})

}
