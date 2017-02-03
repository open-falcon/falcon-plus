package http

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"github.com/open-falcon/falcon-plus/modules/judge/store"
	"net/http"
	"strings"
)

func configInfoRoutes() {
	// e.g. /strategy/lg-dinp-docker01.bj/cpu.idle
	http.HandleFunc("/strategy/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/strategy/"):]
		m := g.StrategyMap.Get()
		RenderDataJson(w, m[urlParam])
	})

	// e.g. /expression/net.port.listen/port=22
	http.HandleFunc("/expression/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/expression/"):]
		m := g.ExpressionMap.Get()
		RenderDataJson(w, m[urlParam])
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		sum := 0
		arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
		for i := 0; i < 16; i++ {
			for j := 0; j < 16; j++ {
				sum += store.HistoryBigMap[arr[i]+arr[j]].Len()
			}
		}

		out := fmt.Sprintf("total: %d\n", sum)
		w.Write([]byte(out))
	})

	http.HandleFunc("/history/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/history/"):]
		pk := utils.Md5(urlParam)
		L, exists := store.HistoryBigMap[pk[0:2]].Get(pk)
		if !exists || L.Len() == 0 {
			w.Write([]byte("not found\n"))
			return
		}

		arr := []string{}

		datas, _ := L.HistoryData(g.Config().Remain - 1)
		for i := 0; i < len(datas); i++ {
			if datas[i] == nil {
				continue
			}

			str := fmt.Sprintf(
				"%d %s %v\n",
				datas[i].Timestamp,
				utils.UnixTsFormat(datas[i].Timestamp),
				datas[i].Value,
			)
			arr = append(arr, str)
		}

		w.Write([]byte(strings.Join(arr, "")))
	})

}
