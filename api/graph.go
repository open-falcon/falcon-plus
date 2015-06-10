package api

import (
	"fmt"
	"log"
	"math"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/index"
	"github.com/open-falcon/graph/proc"
	"github.com/open-falcon/graph/rrdtool"
	"github.com/open-falcon/graph/store"
	//"sync/atomic"
)

//var DropCounter int64

type Graph int

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

	for i := 0; i < count; i++ {
		if items[i] == nil {
			continue
		}
		checksum := items[i].Checksum()

		//statistics
		proc.GraphRpcRecvCnt.Incr()
		proc.RecvDataTrace.Trace(checksum, items[i])
		proc.RecvDataFilter.Filter(checksum, items[i].Value, items[i])

		// To Graph
		first := store.GraphItems.First(checksum)
		if first != nil && items[i].Timestamp <= first.Timestamp {
			continue
		}
		store.GraphItems.PushFront(checksum, items[i])

		// To Index
		index.ReceiveItem(items[i], checksum)
	}
}

func (this *Graph) Query(param cmodel.GraphQueryParam, resp *cmodel.GraphQueryResponse) error {
	// statistics
	proc.GraphQueryCnt.Incr()

	resp.Values = []*cmodel.RRDData{}
	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter)
	if !exists {
		return nil
	}

	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	filename := fmt.Sprintf("%s/%s/%s_%s_%d.rrd", g.Config().RRD.Storage, md5[0:2], md5, dsType, step)
	datas, err := rrdtool.Fetch(filename, param.ConsolFun, param.Start, param.End, step)
	if err != nil {
		if store.GraphItems.LenOf(md5) <= 2 {
			return nil
		}
		// TODO not atomic, fix me
		items := store.GraphItems.PopAll(md5)
		size := len(items)
		if size > 2 {
			filename := fmt.Sprintf("%s/%s/%s_%s_%d.rrd", g.Config().RRD.Storage, md5[0:2],
				md5, items[0].DsType, items[0].Step)
			err := rrdtool.Flush(filename, items)
			if err != nil && g.Config().Debug && g.Config().DebugChecksum == md5 {
				log.Println("flush fail:", err, "filename:", filename)
			}
		} else {
			return nil
		}
	}
	items := store.GraphItems.FetchAll(md5)

	// merge
	items_size := len(items)
	datas_size := len(datas)
	if items_size > 1 && datas_size > 2 &&
		int(datas[1].Timestamp-datas[0].Timestamp) == step &&
		items[items_size-1].Timestamp > datas[0].Timestamp {

		var val cmodel.JsonFloat
		cache_size := int(items[items_size-1].Timestamp-items[0].Timestamp)/step + 1
		cache := make([]*cmodel.RRDData, cache_size, cache_size)

		//fix items
		items_idx := 0
		ts := items[0].Timestamp
		if dsType == g.DERIVE || dsType == g.COUNTER {
			for i := 0; i < cache_size; i++ {
				if items_idx < items_size-1 &&
					ts == items[items_idx].Timestamp &&
					ts != items[items_idx+1].Timestamp {
					val = cmodel.JsonFloat(items[items_idx+1].Value-items[items_idx].Value) /
						cmodel.JsonFloat(items[items_idx+1].Timestamp-items[items_idx].Timestamp)
					if val < 0 {
						val = cmodel.JsonFloat(math.NaN())
					}
					items_idx++
				} else {
					// miss
					val = cmodel.JsonFloat(math.NaN())
				}
				cache[i] = &cmodel.RRDData{
					Timestamp: ts,
					Value:     val,
				}
				ts = ts + int64(step)
			}
		} else if dsType == g.GAUGE {
			for i := 0; i < cache_size; i++ {
				if items_idx < items_size && ts == items[items_idx].Timestamp {
					val = cmodel.JsonFloat(items[items_idx].Value)
					items_idx++
				} else {
					// miss
					val = cmodel.JsonFloat(math.NaN())
				}
				cache[i] = &cmodel.RRDData{
					Timestamp: ts,
					Value:     val,
				}
				ts = ts + int64(step)
			}
		} else {
			log.Println("not support dstype")
			return nil
		}

		size := int(items[items_size-1].Timestamp-datas[0].Timestamp)/step + 1
		ret := make([]*cmodel.RRDData, size, size)
		cache_idx := 0
		ts = datas[0].Timestamp

		if g.Config().Debug && g.Config().DebugChecksum == md5 {
			log.Println("param.start", param.Start, "param.End:", param.End,
				"items:", items, "datas:", datas)
		}

		for i := 0; i < size; i++ {
			if g.Config().Debug && g.Config().DebugChecksum == md5 {
				log.Println("i", i, "size:", size, "items_idx:", items_idx, "ts:", ts)
			}
			if i < datas_size {
				if ts == cache[cache_idx].Timestamp {
					if math.IsNaN(float64(cache[cache_idx].Value)) {
						val = datas[i].Value
					} else {
						val = cache[cache_idx].Value
					}
					cache_idx++
				} else {
					val = datas[i].Value
				}
			} else {
				if cache_idx < cache_size && ts == cache[cache_idx].Timestamp {
					val = cache[cache_idx].Value
					cache_idx++
				} else {
					//miss
					val = cmodel.JsonFloat(math.NaN())
				}
			}
			ret[i] = &cmodel.RRDData{
				Timestamp: ts,
				Value:     val,
			}
			ts = ts + int64(step)
		}
		resp.Values = ret
	} else {
		resp.Values = datas
	}

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.DsType = dsType
	resp.Step = step

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
