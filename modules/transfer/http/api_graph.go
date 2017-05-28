package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	proc "github.com/open-falcon/falcon-plus/modules/transfer/proc"
	graph "github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

type GraphHistoryParam struct {
	Start            int                     `json:"start"`
	End              int                     `json:"end"`
	CF               string                  `json:"cf"`
	EndpointCounters []cmodel.GraphInfoParam `json:"endpoint_counters"`
}

func configApiGraphRoutes() {
	http.HandleFunc("/graph/history", apiGraphHistory)  // post
	http.HandleFunc("/graph/info", apiGraphInfo)        // post
	http.HandleFunc("/graph/last", apiGraphLast)        // post
	http.HandleFunc("/graph/last/raw", apiGraphLastRaw) // post
}

func apiGraphHistory(w http.ResponseWriter, r *http.Request) {
	// statistics
	proc.HistoryRequestCnt.Incr()

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

	data := []*cmodel.GraphQueryResponse{}
	for _, ec := range body.EndpointCounters {
		request := cmodel.GraphQueryParam{
			Start:     int64(body.Start),
			End:       int64(body.End),
			ConsolFun: body.CF,
			Endpoint:  ec.Endpoint,
			Counter:   ec.Counter,
		}
		result, err := graph.QueryOne(request)
		if err != nil {
			log.Printf("graph.queryOne fail, %v", err)
		}
		if result == nil {
			continue
		}
		data = append(data, result)
	}

	// statistics
	proc.HistoryResponseCounterCnt.IncrBy(int64(len(data)))
	for _, item := range data {
		proc.HistoryResponseItemCnt.IncrBy(int64(len(item.Values)))
	}

	StdRender(w, data, nil)
}

func apiGraphInfo(w http.ResponseWriter, r *http.Request) {
	// statistics
	proc.InfoRequestCnt.Incr()

	var body []*cmodel.GraphInfoParam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		StdRender(w, "", err)
		return
	}

	if len(body) == 0 {
		StdRender(w, "", errors.New("empty"))
		return
	}

	data := []*cmodel.GraphFullyInfo{}
	for _, param := range body {
		if param == nil {
			continue
		}
		info, err := graph.Info(*param)
		if err != nil {
			log.Printf("graph.info fail, resp: %v, err: %v", info, err)
		}
		if info == nil {
			continue
		}
		data = append(data, info)
	}

	StdRender(w, data, nil)
}

func apiGraphLast(w http.ResponseWriter, r *http.Request) {
	// statistics
	proc.LastRequestCnt.Incr()

	var body []*cmodel.GraphLastParam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		StdRender(w, "", err)
		return
	}

	if len(body) == 0 {
		StdRender(w, "", errors.New("empty"))
		return
	}

	data := []*cmodel.GraphLastResp{}
	for _, param := range body {
		if param == nil {
			continue
		}
		last, err := graph.Last(*param)
		if err != nil {
			log.Printf("graph.last fail, resp: %v, err: %v", last, err)
		}
		if last == nil {
			continue
		}
		data = append(data, last)
	}

	// statistics
	proc.LastRequestItemCnt.IncrBy(int64(len(data)))

	StdRender(w, data, nil)
}

func apiGraphLastRaw(w http.ResponseWriter, r *http.Request) {
	// statistics
	proc.LastRawRequestCnt.Incr()

	var body []*cmodel.GraphLastParam
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		StdRender(w, "", err)
		return
	}

	if len(body) == 0 {
		StdRender(w, "", errors.New("empty"))
		return
	}

	data := []*cmodel.GraphLastResp{}
	for _, param := range body {
		if param == nil {
			continue
		}
		last, err := graph.LastRaw(*param)
		if err != nil {
			log.Printf("graph.last.raw fail, resp: %v, err: %v", last, err)
		}
		if last == nil {
			continue
		}
		data = append(data, last)
	}
	// statistics
	proc.LastRawRequestItemCnt.IncrBy(int64(len(data)))
	StdRender(w, data, nil)
}
