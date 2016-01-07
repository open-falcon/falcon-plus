package api

import (
	"fmt"
	"math"
	"time"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/index"
	"github.com/open-falcon/graph/proc"
	"github.com/open-falcon/graph/rrdtool"
	"github.com/open-falcon/graph/store"
)

type Graph int

func (this *Graph) GetRrd(key string, rrdfile *g.File) (err error) {
	if md5, dsType, step, err := g.SplitRrdCacheKey(key); err != nil {
		return err
	} else {
		rrdfile.Filename = g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step)
	}

	items := store.GraphItems.PopAll(key)
	if len(items) > 0 {
		rrdtool.FlushFile(rrdfile.Filename, items)
	}

	rrdfile.Body, err = rrdtool.ReadFile(rrdfile.Filename)
	return
}

func (this *Graph) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (this *Graph) Send(items []*cmodel.GraphItem, resp *cmodel.SimpleRpcResponse) error {
	go handleItems(items)
	return nil
}

// 供外部调用、处理接收到的数据 的接口
func HandleItems(items []*cmodel.GraphItem) error {
	handleItems(items)
	return nil
}

func handleItems(items []*cmodel.GraphItem) {
	if items == nil {
		return
	}

	count := len(items)
	if count == 0 {
		return
	}

	cfg := g.Config()

	for i := 0; i < count; i++ {
		if items[i] == nil {
			continue
		}
		dsType := items[i].DsType
		step := items[i].Step
		checksum := items[i].Checksum()
		key := g.FormRrdCacheKey(checksum, dsType, step)

		//statistics
		proc.GraphRpcRecvCnt.Incr()

		// To Graph
		first := store.GraphItems.First(key)
		if first != nil && items[i].Timestamp <= first.Timestamp {
			continue
		}
		store.GraphItems.PushFront(key, items[i], checksum, cfg)

		// To Index
		index.ReceiveItem(items[i], checksum)

		// To History
		store.AddItem(checksum, items[i])
	}
}

func (this *Graph) Query(param cmodel.GraphQueryParam, resp *cmodel.GraphQueryResponse) error {
	var (
		datas      []*cmodel.RRDData
		datas_size int
		retry_nb   int
	)

	// statistics
	proc.GraphQueryCnt.Incr()

	cfg := g.Config()

	// form empty response
	resp.Values = []*cmodel.RRDData{}
	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter) // complete dsType and step
	if !exists {
		return nil
	}
	resp.DsType = dsType
	resp.Step = step

	start_ts := param.Start - param.Start%int64(step)
	end_ts := param.End - param.End%int64(step) + int64(step)
	if end_ts-start_ts-int64(step) < 1 {
		return nil
	}

	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	key := g.FormRrdCacheKey(md5, dsType, step)
	filename := g.RrdFileName(cfg.RRD.Storage, md5, dsType, step)

	// read cached items
	items, flag := store.GraphItems.FetchAll(key)
	items_size := len(items)

retry:
	if cfg.Migrate.Enabled && flag&g.GRAPH_F_MISS != 0 {
		/* Code is ugly, but sort of works, chan maybe better */
		/* https://github.com/tjgq/broadcast/blob/master/broadcast.go  */
		if flag&g.GRAPH_F_FETCHING != 0 && retry_nb < 20 {
			time.Sleep(time.Millisecond * 100)
			retry_nb++
			goto retry
		}
		//update items befor get from remote
		node, _ := rrdtool.Consistent.Get(param.Endpoint + "/" + param.Counter)
		done := make(chan error, 1)
		if items_size > 0 {
			//rrdtool.Task_ch
			rrdtool.Task_ch[node] <- &rrdtool.Task_t{
				Method: "Graph.Send",
				Done:   done,
				Key:    key,
			}
			<-done
		}
		rrdtool.Task_ch[node] <- &rrdtool.Task_t{
			Method: "Graph.Query",
			Done:   done,
			Args:   param,
			Reply:  resp,
		}
		err := <-done
		// danger!!! query call query!!!
		return err
	}

	// read data from rrd file
	datas, _ = rrdtool.Fetch(filename, param.ConsolFun, start_ts, end_ts, step)
	datas_size = len(datas)

	nowTs := time.Now().Unix()
	lastUpTs := nowTs - nowTs%int64(step)
	rra1StartTs := lastUpTs - int64(rrdtool.RRA1PointCnt*step)
	// consolidated, do not merge
	if start_ts < rra1StartTs {
		resp.Values = datas
		goto _RETURN_OK
	}

	// no cached items, do not merge
	if items_size < 1 {
		resp.Values = datas
		goto _RETURN_OK
	}

	// merge
	{
		// fmt cached items
		var val cmodel.JsonFloat
		cache := make([]*cmodel.RRDData, 0)

		ts := items[0].Timestamp
		itemEndTs := items[items_size-1].Timestamp
		itemIdx := 0
		if dsType == g.DERIVE || dsType == g.COUNTER {
			for ts < itemEndTs {
				if itemIdx < items_size-1 && ts == items[itemIdx].Timestamp &&
					ts == items[itemIdx+1].Timestamp-int64(step) {
					val = cmodel.JsonFloat(items[itemIdx+1].Value-items[itemIdx].Value) / cmodel.JsonFloat(step)
					if val < 0 {
						val = cmodel.JsonFloat(math.NaN())
					}
					itemIdx++
				} else {
					// missing
					val = cmodel.JsonFloat(math.NaN())
				}

				if ts >= start_ts && ts <= end_ts {
					cache = append(cache, &cmodel.RRDData{Timestamp: ts, Value: val})
				}
				ts = ts + int64(step)
			}
		} else if dsType == g.GAUGE {
			for ts <= itemEndTs {
				if itemIdx < items_size && ts == items[itemIdx].Timestamp {
					val = cmodel.JsonFloat(items[itemIdx].Value)
					itemIdx++
				} else {
					// missing
					val = cmodel.JsonFloat(math.NaN())
				}

				if ts >= start_ts && ts <= end_ts {
					cache = append(cache, &cmodel.RRDData{Timestamp: ts, Value: val})
				}
				ts = ts + int64(step)
			}
		}
		cache_size := len(cache)

		// do merging
		merged := make([]*cmodel.RRDData, 0)
		if datas_size > 0 {
			for _, val := range datas {
				if val.Timestamp >= start_ts && val.Timestamp <= end_ts {
					merged = append(merged, val) //rrdtool返回的数据,时间戳是连续的、不会有跳点的情况
				}
			}
		}

		if cache_size > 0 {
			rrdDataSize := len(merged)
			lastTs := cache[0].Timestamp

			// find junction
			rrdDataIdx := 0
			for rrdDataIdx = rrdDataSize - 1; rrdDataIdx >= 0; rrdDataIdx-- {
				if merged[rrdDataIdx].Timestamp < cache[0].Timestamp {
					lastTs = merged[rrdDataIdx].Timestamp
					break
				}
			}

			// fix missing
			for ts := lastTs + int64(step); ts < cache[0].Timestamp; ts += int64(step) {
				merged = append(merged, &cmodel.RRDData{Timestamp: ts, Value: cmodel.JsonFloat(math.NaN())})
			}

			// merge cached items to result
			rrdDataIdx += 1
			for cacheIdx := 0; cacheIdx < cache_size; cacheIdx++ {
				if rrdDataIdx < rrdDataSize {
					if !math.IsNaN(float64(cache[cacheIdx].Value)) {
						merged[rrdDataIdx] = cache[cacheIdx]
					}
				} else {
					merged = append(merged, cache[cacheIdx])
				}
				rrdDataIdx++
			}
		}
		mergedSize := len(merged)

		// fmt result
		ret_size := int((end_ts - start_ts) / int64(step))
		ret := make([]*cmodel.RRDData, ret_size, ret_size)
		mergedIdx := 0
		ts = start_ts
		for i := 0; i < ret_size; i++ {
			if mergedIdx < mergedSize && ts == merged[mergedIdx].Timestamp {
				ret[i] = merged[mergedIdx]
				mergedIdx++
			} else {
				ret[i] = &cmodel.RRDData{Timestamp: ts, Value: cmodel.JsonFloat(math.NaN())}
			}
			ts += int64(step)
		}
		resp.Values = ret
	}

_RETURN_OK:
	// statistics
	proc.GraphQueryItemCnt.IncrBy(int64(len(resp.Values)))
	return nil
}

func (this *Graph) Info(param cmodel.GraphInfoParam, resp *cmodel.GraphInfoResp) error {
	// statistics
	proc.GraphInfoCnt.Incr()

	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter)
	if !exists {
		return nil
	}

	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	filename := fmt.Sprintf("%s/%s/%s_%s_%d.rrd", g.Config().RRD.Storage, md5[0:2], md5, dsType, step)

	resp.ConsolFun = dsType
	resp.Step = step
	resp.Filename = filename

	return nil
}

func (this *Graph) Last(param cmodel.GraphLastParam, resp *cmodel.GraphLastResp) error {
	// statistics
	proc.GraphLastCnt.Incr()

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = GetLast(param.Endpoint, param.Counter)

	return nil
}

func (this *Graph) LastRaw(param cmodel.GraphLastParam, resp *cmodel.GraphLastResp) error {
	// statistics
	proc.GraphLastRawCnt.Incr()

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = GetLastRaw(param.Endpoint, param.Counter)

	return nil
}

// 非法值: ts=0,value无意义
func GetLast(endpoint, counter string) *cmodel.RRDData {
	dsType, step, exists := index.GetTypeAndStep(endpoint, counter)
	if !exists {
		return cmodel.NewRRDData(0, 0.0)
	}

	if dsType == g.GAUGE {
		return GetLastRaw(endpoint, counter)
	}

	if dsType == g.COUNTER || dsType == g.DERIVE {
		md5 := cutils.Md5(endpoint + "/" + counter)
		items := store.GetAllItems(md5)
		if len(items) < 2 {
			return cmodel.NewRRDData(0, 0.0)
		}

		f0 := items[0]
		f1 := items[1]
		delta_ts := f0.Timestamp - f1.Timestamp
		delta_v := f0.Value - f1.Value
		if delta_ts != int64(step) || delta_ts <= 0 {
			return cmodel.NewRRDData(0, 0.0)
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return cmodel.NewRRDData(f0.Timestamp, delta_v/float64(delta_ts))
	}

	return cmodel.NewRRDData(0, 0.0)
}

// 非法值: ts=0,value无意义
func GetLastRaw(endpoint, counter string) *cmodel.RRDData {
	md5 := cutils.Md5(endpoint + "/" + counter)
	item := store.GetLastItem(md5)
	return cmodel.NewRRDData(item.Timestamp, item.Value)
}
