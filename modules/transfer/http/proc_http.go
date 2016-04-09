package http

import (
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/transfer/proc"
	"github.com/open-falcon/transfer/sender"
	"net/http"
	"strconv"
	"strings"
)

func configProcHttpRoutes() {
	// counter
	http.HandleFunc("/counter/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// TO BE DISCARDed
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// step
	http.HandleFunc("/proc/step", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, map[string]interface{}{"min_step": sender.MinStep})
	})

	// trace
	http.HandleFunc("/trace/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/trace/"):]
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
		proc.RecvDataTrace.SetPK(cutils.PK(endpoint, metric, tags))
		RenderDataJson(w, proc.RecvDataTrace.GetAllTraced())
	})

	// filter
	http.HandleFunc("/filter/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/filter/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		opt := args[2]

		threadholdStr := args[3]
		threadhold, err := strconv.ParseFloat(threadholdStr, 64)
		if err != nil {
			RenderDataJson(w, "bad threadhold")
			return
		}

		tags := make(map[string]string)
		if argsLen > 4 {
			tagVals := strings.Split(args[4], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}

		err = proc.RecvDataFilter.SetFilter(cutils.PK(endpoint, metric, tags), opt, threadhold)
		if err != nil {
			RenderDataJson(w, err.Error())
			return
		}

		RenderDataJson(w, proc.RecvDataFilter.GetAllFiltered())
	})
}
