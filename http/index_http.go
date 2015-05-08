package http

import (
	"fmt"
	"github.com/open-falcon/graph/index"
	"net/http"
	"strconv"
	"strings"
)

func configIndexRoutes() {
	// 触发索引全量更新, 同步操作
	http.HandleFunc("/index/updateAll", func(w http.ResponseWriter, r *http.Request) {
		go index.UpdateIndexAllByDefaultStep()
		RenderDataJson(w, "ok")
	})

	// 获取索引全量更新的并行数
	http.HandleFunc("/index/updateAll/concurrent", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, index.GetConcurrentOfUpdateIndexAll())
	})

	// 更新一条索引数据,用于手动建立索引 endpoint metric step dstype tags
	http.HandleFunc("/index/update/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/index/update/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		if !(argsLen == 4 || argsLen == 5) {
			RenderDataJson(w, "bad args")
			return
		}
		endpoint := args[0]
		metric := args[1]
		step, _ := strconv.ParseInt(args[2], 10, 32)
		dstype := args[3]
		tags := make(map[string]string)
		if argsLen == 5 {
			tagVals := strings.Split(args[4], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		err := index.UpdateIndexOne(endpoint, metric, tags, dstype, int(step))
		if err != nil {
			RenderDataJson(w, fmt.Sprintf("%v", err))
			return
		}

		RenderDataJson(w, "ok")
	})
}
