package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-falcon/falcon-plus/modules/graph/index"
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

	// index.cached
	http.HandleFunc("/index/cache/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/index/cache/"):]
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

		item, err := index.GetIndexedItemCache(endpoint, metric, tags, dstype, int(step))
		if err != nil {
			RenderDataJson(w, fmt.Sprintf("%v", err))
			return
		}

		RenderDataJson(w, item)
	})

	http.HandleFunc("/v2/index/cache", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0 && len(r.Form["step"]) > 0 && len(r.Form["type"]) > 0) {
			RenderDataJson(w, "bad args")
			return
		}
		endpoint := r.Form["e"][0]
		metric := r.Form["m"][0]
		step, _ := strconv.ParseInt(r.Form["step"][0], 10, 32)
		dstype := r.Form["type"][0]

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

		item, err := index.GetIndexedItemCache(endpoint, metric, tags, dstype, int(step))
		if err != nil {
			RenderDataJson(w, fmt.Sprintf("%v", err))
			return
		}

		RenderDataJson(w, item)
	})

}
