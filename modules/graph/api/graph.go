// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	pfc "github.com/niean/goperfcounter"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

type Graph int

func (this *Graph) GetRrd(key string, rrdfile *g.File) (err error) {
	var (
		md5    string
		dsType string
		step   int
	)
	if md5, dsType, step, err = g.SplitRrdCacheKey(key); err != nil {
		return err
	} else {
		rrdfile.Filename = g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step)
	}

	items := store.GraphItems.PopAll(key)
	if len(items) > 0 {
		rrdtool.FlushFile(rrdfile.Filename, md5, items)
	}

	rrdfile.Body, err = rrdtool.ReadFile(rrdfile.Filename, md5)
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

		endpoint := items[i].Endpoint
		if !g.IsValidString(endpoint) {
			if cfg.Debug {
				log.Printf("invalid endpoint: %s", endpoint)
			}
			pfc.Meter("invalidEnpoint", 1)
			continue
		}

		counter := cutils.Counter(items[i].Metric, items[i].Tags)
		if !g.IsValidString(counter) {
			if cfg.Debug {
				log.Printf("invalid counter: %s/%s", endpoint, counter)
			}
			pfc.Meter("invalidCounter", 1)
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
		rrdDatas    []*cmodel.RRDData
		rrdDataSize int
		sample      int
		err         error
	)

	// statistics
	proc.GraphQueryCnt.Incr()

	// form empty response
	resp.Values = []*cmodel.RRDData{}
	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter)
	if !exists {
		return nil
	}
	resp.DsType = dsType
	resp.Step = step

	now := time.Now().Unix() - time.Now().Unix()%int64(step)
	start := param.Start - param.Start%int64(step)
	end := param.End - param.End%int64(step) + int64(step)
	if end > now {
		end = now
	}
	if end-start-int64(step) < 1 {
		return fmt.Errorf("start ts and end ts is invalid")
	}

	// compute and check cache
	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	ckey := g.FormRrdCacheKey(md5, dsType, step)
	items, flag := store.GraphItems.FetchAll(ckey)
	cache := checkCacheItem(items, start, end, step, dsType)
	cacheSize := len(cache)

	//fetch rrd data
	cfg := g.Config()
	if cfg.Migrate.Enabled && flag&g.GRAPH_F_MISS != 0 {
		node, _ := rrdtool.Consistent.Get(param.Endpoint + "/" + param.Counter)
		done := make(chan error, 1)
		res := &cmodel.GraphAccurateQueryResponse{}
		rrdtool.Net_task_ch[node] <- &rrdtool.Net_task_t{
			Method: rrdtool.NET_TASK_M_QUERY,
			Done:   done,
			Args:   param,
			Reply:  res,
		}
		<-done
		rrdDatas = res.Values
		sample, rrdDataSize = getSampleAndSize(rrdDatas, step)
	} else {
		rrdDatas, err = getRrdData(param, start, end, now, step, dsType)
		if err != nil {
			if cfg.Debug {
				log.Printf("request %s:%s:%d", resp.Endpoint, resp.Counter, len(resp.Values))
			}
			proc.GraphQueryItemCnt.IncrBy(int64(len(resp.Values)))
			return err
		}
		sample, rrdDataSize = getSampleAndSize(rrdDatas, step)
	}

	if cacheSize == 0 {
		resp.Values = rrdDatas
		if cfg.Debug {
			log.Printf("request %s:%s:%d", resp.Endpoint, resp.Counter, len(resp.Values))
		}
		proc.GraphQueryItemCnt.IncrBy(int64(len(resp.Values)))
		return nil
	}

	//complement and consolidate cache
	lastTs := cache[0].Timestamp
	rrdDataIdx := rrdDataSize - 1
	for ; rrdDataIdx >= 0; rrdDataIdx-- {
		if rrdDatas[rrdDataIdx].Timestamp < cache[0].Timestamp {
			lastTs = rrdDatas[rrdDataIdx].Timestamp
			break
		}
	}
	fullCache := make([]*cmodel.RRDData, 0)
	for ts := lastTs + int64(step); ts < cache[0].Timestamp; ts += int64(step) {
		fullCache = append(fullCache, &cmodel.RRDData{Timestamp: ts, Value: cmodel.JsonFloat(math.NaN())})
	}
	fullCache = append(fullCache, cache...)
	for ts := cache[cacheSize-1].Timestamp + int64(step); ts <= end; ts = ts + int64(step) {
		fullCache = append(fullCache, &cmodel.RRDData{Timestamp: ts, Value: cmodel.JsonFloat(math.NaN())})
	}
	xff := 0.5
	sampleCache := consolidate(fullCache, xff, param.ConsolFun, sample)

	//combine rrd data with conslidated cache
	result := make([]*cmodel.RRDData, 0)
	result = append(result, rrdDatas[:rrdDataIdx+1]...)
	result = append(result, sampleCache...)

	resp.Values = result
	return nil
}

func getSampleAndSize(rrdDatas []*cmodel.RRDData, step int) (int, int) {
	var rrdDataSize, sample int
	rrdDataSize = len(rrdDatas)
	if rrdDataSize >= 2 {
		sample = int(rrdDatas[1].Timestamp-rrdDatas[0].Timestamp) / step
	} else {
		sample = 1
	}
	return sample, rrdDataSize
}

func getRrdData(param cmodel.GraphQueryParam, start int64, end int64, now int64, step int, dsType string) ([]*cmodel.RRDData, error) {
	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	filename := g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step)

	datas := make([]*cmodel.RRDData, 0)
	rrdDatas := make([]*cmodel.RRDData, 0)
	rra1StartTs := now - int64(rrdtool.RRA1PointCnt*step)

	if start < rra1StartTs {
		datas, _ = rrdtool.Fetch(filename, md5, param.ConsolFun, start-int64(step), end, step)
		if len(datas) >= 2 && datas[1].Timestamp-datas[0].Timestamp <= int64(step) {
			return rrdDatas, fmt.Errorf("Fetching data from rrd fails")
		}
	} else {
		//注意:获取12个小时以内的数据时,必须使用AVERAGE方法
		datas, _ = rrdtool.Fetch(filename, md5, "AVERAGE", start-int64(step), end, step)
	}

	for _, val := range datas {
		if val.Timestamp >= start && val.Timestamp <= end {
			rrdDatas = append(rrdDatas, val)
		}
	}

	return rrdDatas, nil
}

func checkCacheItem(items []*cmodel.GraphItem, start int64, end int64, step int, dsType string) []*cmodel.RRDData {
	cache := make([]*cmodel.RRDData, 0)
	if len(items) == 0 {
		return cache
	}

	var val cmodel.JsonFloat
	ts := items[0].Timestamp
	itemsSize := len(items)
	itemEndTs := items[itemsSize-1].Timestamp
	itemIdx := 0
	if dsType == g.DERIVE || dsType == g.COUNTER {
		for ts <= itemEndTs-int64(step) {
			val = cmodel.JsonFloat(math.NaN())
			if itemIdx < itemsSize-1 && ts == items[itemIdx].Timestamp {
				if ts+int64(step) == items[itemIdx+1].Timestamp {
					dt := items[itemIdx+1].Timestamp - ts
					val = cmodel.JsonFloat(items[itemIdx+1].Value-items[itemIdx].Value) / cmodel.JsonFloat(dt)
				}
				if val < 0 {
					val = cmodel.JsonFloat(math.NaN())
				}
				itemIdx++
			}
			if ts+int64(step) >= start && ts+int64(step) <= end {
				cache = append(cache, &cmodel.RRDData{Timestamp: ts + int64(step), Value: val})
			}
			ts = ts + int64(step)
		}
	} else if dsType == g.GAUGE {
		for ts <= itemEndTs {
			if itemIdx < itemsSize && ts == items[itemIdx].Timestamp {
				val = cmodel.JsonFloat(items[itemIdx].Value)
				itemIdx++
			} else {
				val = cmodel.JsonFloat(math.NaN())
			}

			if ts >= start && ts <= end {
				cache = append(cache, &cmodel.RRDData{Timestamp: ts, Value: val})
			}
			ts = ts + int64(step)
		}
	}
	return cache
}

func consolidate(fullCache []*cmodel.RRDData, xff float64, cf string, sample int) []*cmodel.RRDData {
	var val cmodel.JsonFloat
	result := make([]*cmodel.RRDData, 0)
	for i := 0; i < len(fullCache); i = i + sample {
		res := make([]*cmodel.RRDData, 0)
		num := 0
		for j := i; j < i+sample && j < len(fullCache); j++ {
			if !isNaN(fullCache[i].Value) {
				num++
			}
			res = append(res, fullCache[j])
		}
		if len(res) == sample {
			ts := fullCache[i+sample-1].Timestamp
			if float64(num) >= xff*float64(sample) {
				val = math2(res, cf)
				result = append(result, &cmodel.RRDData{Timestamp: ts, Value: val})
			} else {
				result = append(result, &cmodel.RRDData{Timestamp: ts, Value: cmodel.JsonFloat(math.NaN())})
			}
		}
	}
	return result
}

func math2(res []*cmodel.RRDData, cf string) cmodel.JsonFloat {
	upperCf := strings.ToUpper(cf)
	var result cmodel.JsonFloat
	switch upperCf {
	case "AVERAGE":
		result = Average(res)
	case "MAX":
		result = Max(res)
	case "MIN":
		result = Min(res)
	}
	return result
}

func Average(res []*cmodel.RRDData) cmodel.JsonFloat {
	var sum cmodel.JsonFloat
	var num int64
	for _, v := range res {
		if !isNaN(v.Value) {
			sum = sum + v.Value
			num++
		}
	}
	return sum / cmodel.JsonFloat(num)
}

func Min(res []*cmodel.RRDData) cmodel.JsonFloat {
	min := cmodel.JsonFloat(math.NaN())
	for _, v := range res {
		if isNaN(min) && !isNaN(v.Value) {
			min = v.Value
		}
		if !isNaN(min) && !isNaN(v.Value) && min > v.Value {
			min = v.Value
		}
	}
	return min
}

func Max(res []*cmodel.RRDData) cmodel.JsonFloat {
	max := cmodel.JsonFloat(math.NaN())
	for _, v := range res {
		if isNaN(max) && !isNaN(v.Value) {
			max = v.Value
		}
		if !isNaN(max) && !isNaN(v.Value) && max < v.Value {
			max = v.Value
		}
	}
	return max
}

func isNaN(v cmodel.JsonFloat) bool {
	return math.IsNaN(float64(v))
}

//从内存索引、MySQL中删除counter，并从磁盘上删除对应rrd文件
func (this *Graph) Delete(params []*cmodel.GraphDeleteParam, resp *cmodel.GraphDeleteResp) error {
	resp = &cmodel.GraphDeleteResp{}
	for _, param := range params {
		err, tags := cutils.SplitTagsString(param.Tags)
		if err != nil {
			log.Error("invalid tags:", param.Tags, "error:", err)
			continue
		}

		var item *cmodel.GraphItem = &cmodel.GraphItem{
			Endpoint: param.Endpoint,
			Metric:   param.Metric,
			Tags:     tags,
			DsType:   param.DsType,
			Step:     param.Step,
		}
		index.RemoveItem(item)
	}

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
