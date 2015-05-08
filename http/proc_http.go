package http

import (
	MP "github.com/open-falcon/common/proc"
	"github.com/open-falcon/transfer/g"
	"github.com/open-falcon/transfer/proc"
	"net/http"
	"strings"
)

func configProcHttpRoutes() {
	// TOP
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	http.HandleFunc("/statistics/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.Config())
	})

	// 向每个节点的发送计数
	http.HandleFunc("/statistics/sendCntPerNode", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, []map[string]*MP.SCounterQps{proc.SendToJudgeCntPerNode, proc.SendToGraphCntPerNode,
			proc.SendToGraphMigratingCntPerNode})
	})

	// 向每个节点发送的丢弃计数
	http.HandleFunc("/statistics/sendDropCntPerNode", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, []map[string]*MP.SCounterQps{proc.SendToJudgeDropCntPerNode, proc.SendToGraphDropCntPerNode,
			proc.SendToGraphMigratingDropCntPerNode})
	})

	// 发送缓存
	http.HandleFunc("/statistics/sendCacheSize", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, []*MP.SCounterBase{proc.JudgeQueuesCnt, proc.GraphQueuesCnt,
			proc.GraphMigratingQueuesCnt})
	})
	http.HandleFunc("/statistics/sendCacheSizePerNode", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, []map[string]*MP.SCounterBase{proc.JudgeQueuesCntPerNode, proc.GraphQueuesCntPerNode,
			proc.GraphMigratingQueuesCntPerNode})
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
		proc.RecvDataTrace.SetTraceConfig(endpoint, metric, tags)
		RenderDataJson(w, proc.RecvDataTrace.FilterAll())
	})
}
