package http

import (
	"encoding/json"
	"errors"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/query/graph"
	"github.com/toolkits/logger"
	"net/http"
	"strconv"
	"time"
)

type GraphHistoryParam struct {
	Start            int                    `json:"start"`
	End              int                    `json:"end"`
	CF               string                 `json:"cf"`
	EndpointCounters []model.GraphInfoParam `json:"endpoint_counters"`
}

func configGraphRoutes() {

	// method:post
	http.HandleFunc("/graph/history", func(w http.ResponseWriter, r *http.Request) {
		var body GraphHistoryParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body.EndpointCounters) == 0 {
			StdRender(w, "", errors.New("empty_payload"))
			return
		}

		data := []*model.GraphQueryResponse{}
		for _, ec := range body.EndpointCounters {
			result, err := graph.QueryOne(int64(body.Start), int64(body.End), body.CF, ec.Endpoint, ec.Counter)
			if err != nil {
				logger.Error("query one fail: %v", err)
			}
			data = append(data, result)
		}

		StdRender(w, data, nil)
	})

	// method:get
	http.HandleFunc("/graph/history/one", func(w http.ResponseWriter, r *http.Request) {
		start := r.FormValue("start")
		end := r.FormValue("end")
		cf := r.FormValue("cf")
		endpoint := r.FormValue("endpoint")
		counter := r.FormValue("counter")

		if endpoint == "" || counter == "" {
			StdRender(w, "", errors.New("empty_endpoint_counter"))
			return
		}

		if cf != "AVERAGE" && cf != "MAX" && cf != "MIN" {
			StdRender(w, "", errors.New("invalid_cf"))
			return
		}

		now := time.Now()
		start_i64, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			start_i64 = now.Unix() - 3600
		}
		end_i64, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			end_i64 = now.Unix()
		}

		result, err := graph.QueryOne(start_i64, end_i64, cf, endpoint, counter)
		logger.Trace("query one result: %v, err: %v", result, err)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		StdRender(w, result, nil)
	})

	// get, info
	http.HandleFunc("/graph/info/one", func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.FormValue("endpoint")
		counter := r.FormValue("counter")

		if endpoint == "" || counter == "" {
			StdRender(w, "", errors.New("empty_endpoint_counter"))
			return
		}

		result, err := graph.Info(endpoint, counter)
		logger.Trace("graph.info result: %v, err: %v", result, err)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		StdRender(w, result, nil)
	})

	// post, info
	http.HandleFunc("/graph/info", func(w http.ResponseWriter, r *http.Request) {
		var body []*model.GraphInfoParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body) == 0 {
			StdRender(w, "", errors.New("empty_payload"))
			return
		}

		data := []*model.GraphFullyInfo{}
		for _, param := range body {
			info, err := graph.Info(param.Endpoint, param.Counter)
			if err != nil {
				logger.Trace("graph.info fail, resp: %v, err: %v", info, err)
			} else {
				logger.Trace("graph.info result: %v, err: %v", info, err)
			}
			data = append(data, info)
		}

		StdRender(w, data, nil)
	})

}
