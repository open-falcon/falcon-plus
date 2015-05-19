package http

import (
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/api"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func configDebugRoutes() {
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < store.GraphItems.Size; i++ {
			keys := store.GraphItems.KeysByIndex(i)
			if len(keys) == 0 {
				w.Write([]byte(fmt.Sprintf("%d\n", 0)))
				return
			}

			oneHourAgo := time.Now().Unix() - 3600

			count := 0
			for _, checksum := range keys {
				item := store.GraphItems.First(checksum)
				if item == nil {
					continue
				}

				if item.Timestamp > oneHourAgo {
					count++
				}
			}

			w.Write([]byte(fmt.Sprintf("%d\n", count)))
		}
	})

	// 接收数据 endpoint metric ts step dstype value [tags]
	http.HandleFunc("/api/recv/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/api/recv/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		if !(argsLen == 6 || argsLen == 7) {
			RenderDataJson(w, "bad args")
			return
		}

		endpoint := args[0]
		metric := args[1]
		ts, _ := strconv.ParseInt(args[2], 10, 64)
		step, _ := strconv.ParseInt(args[3], 10, 32)
		dstype := args[4]
		value, _ := strconv.ParseFloat(args[5], 64)
		tags := make(map[string]string)
		if argsLen == 7 {
			tags = cutils.DictedTagstring(args[6])
		}

		item := &cmodel.MetaData{
			Endpoint:    endpoint,
			Metric:      metric,
			Timestamp:   ts,
			Step:        step,
			CounterType: dstype,
			Value:       value,
			Tags:        tags,
		}
		gitem, err := convert2GraphItem(item)
		if err != nil {
			RenderDataJson(w, err)
			return
		}

		api.HandleItems([]*cmodel.GraphItem{gitem})
		RenderDataJson(w, "ok")
	})

}

func convert2GraphItem(d *cmodel.MetaData) (*cmodel.GraphItem, error) {
	item := &cmodel.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < g.MIN_STEP {
		item.Step = g.MIN_STEP
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == g.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == g.COUNTER {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == g.DERIVE {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}
