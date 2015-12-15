package rrdtool

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/rrdlite"
	"github.com/toolkits/file"

	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
)

var Counter uint64

func Start() {
	// check data dir
	dataDir := g.Config().RRD.Storage
	if err := file.EnsureDirRW(dataDir); err != nil {
		log.Fatalln("rrdtool.Start error, bad data dir", dataDir+",", err)
	}

	// sync disk
	go syncDisk()
	log.Println("rrdtool.Start ok")
}

// RRD Files' Lock
type RRDLocker struct {
	sync.Mutex
	M map[string]*sync.Mutex
}

func (t *RRDLocker) GetLock(key string) *sync.Mutex {
	t.Lock()
	defer t.Unlock()

	if lock, exists := t.M[key]; !exists {
		t.M[key] = new(sync.Mutex)
		return t.M[key]
	} else {
		return lock
	}
}

var (
	L *RRDLocker = &RRDLocker{
		M: make(map[string]*sync.Mutex),
	}
)

// RRA.Point.Size
const (
	RRA1PointCnt   = 720 // 1m一个点存12h
	RRA5PointCnt   = 576 // 5m一个点存2d
	RRA20PointCnt  = 504 // 20m一个点存7d
	RRA180PointCnt = 766 // 3h一个点存3month
	RRA720PointCnt = 730 // 12h一个点存1year
)

func create(filename string, item *cmodel.GraphItem) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(item.Step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", item.DsType, item.Heartbeat, item.Min, item.Max)

	// 设置各种归档策略
	// 1分钟一个点存 12小时
	c.RRA("AVERAGE", 0.5, 1, RRA1PointCnt)

	// 5m一个点存2d
	c.RRA("AVERAGE", 0.5, 5, RRA5PointCnt)
	c.RRA("MAX", 0.5, 5, RRA5PointCnt)
	c.RRA("MIN", 0.5, 5, RRA5PointCnt)

	// 20m一个点存7d
	c.RRA("AVERAGE", 0.5, 20, RRA20PointCnt)
	c.RRA("MAX", 0.5, 20, RRA20PointCnt)
	c.RRA("MIN", 0.5, 20, RRA20PointCnt)

	// 3小时一个点存3个月
	c.RRA("AVERAGE", 0.5, 180, RRA180PointCnt)
	c.RRA("MAX", 0.5, 180, RRA180PointCnt)
	c.RRA("MIN", 0.5, 180, RRA180PointCnt)

	// 12小时一个点存1year
	c.RRA("AVERAGE", 0.5, 720, RRA720PointCnt)
	c.RRA("MAX", 0.5, 720, RRA720PointCnt)
	c.RRA("MIN", 0.5, 720, RRA720PointCnt)

	return c.Create(true)
}

func update(filename string, items []*cmodel.GraphItem) error {
	u := rrdlite.NewUpdater(filename)

	for _, item := range items {
		v := math.Abs(item.Value)
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if item.DsType == "DERIVE" || item.DsType == "COUNTER" {
			u.Cache(item.Timestamp, int(item.Value))
		} else {
			u.Cache(item.Timestamp, item.Value)
		}
	}

	return u.Update()
}

// flush to disk from memory
// 最新的数据在列表的最后面
// TODO fix me, filename fmt from item[0], it's hard to keep consistent
func Flush(filename string, items []*cmodel.GraphItem) error {
	if items == nil || len(items) == 0 {
		return errors.New("empty items")
	}

	lock := L.GetLock(filename)
	lock.Lock()
	defer lock.Unlock()

	if !g.IsRrdFileExist(filename) {
		baseDir := file.Dir(filename)

		err := file.InsureDir(baseDir)
		if err != nil {
			return err
		}

		err = create(filename, items[0])
		if err != nil {
			return err
		}
	}

	return update(filename, items)
}

func Fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

	lock := L.GetLock(filename)
	lock.Lock()
	defer lock.Unlock()

	fetchRes, err := rrdlite.Fetch(filename, cf, start_t, end_t, step_t)
	if err != nil {
		return []*cmodel.RRDData{}, err
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	ret := make([]*cmodel.RRDData, size)

	start_ts := fetchRes.Start.Unix()
	step_s := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := start_ts + int64(i+1)*int64(step_s)
		d := &cmodel.RRDData{
			Timestamp: ts,
			Value:     cmodel.JsonFloat(val),
		}
		ret[i] = d
	}

	return ret, nil
}

func FlushAll() {
	n := store.GraphItems.Size / 10
	for i := 0; i < store.GraphItems.Size; i++ {
		FlushRRD(i)
		if i%n == 0 {
			log.Println("flush hash idx:", i, "size", store.GraphItems.Size, "counter", Counter)
		}
	}
	log.Println("flush hash done. counter", Counter)
}

func FlushRRD(idx int) {
	cfg := g.Config()

	storageDir := cfg.RRD.Storage

	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}

	for _, ckey := range keys {
		// get md5, dstype, step
		checksum, dsType, step, err := g.SplitRrdCacheKey(ckey)
		if err != nil {
			continue
		}

		items, flag := store.GraphItems.PopAll(ckey)
		size := len(items)
		if size == 0 {
			continue
		}

		if flag&store.GRAPH_F_MISS != 0 {
			//remote
			store.GraphItems.SetFlag(ckey, flag|store.GRAPH_F_SENDING)
			defer store.GraphItems.SetFlag(ckey, flag)
			node, _ := store.GraphItems.Consistent.Get(items[0].PrimaryKey())
			client := store.GraphItems.Client[node]
			resp := &cmodel.SimpleRpcResponse{}
			err := store.Jsonrpc_call(client, "Graph.Send", items, resp,
				time.Duration(cfg.CallTimeout)*time.Millisecond)
			if err != nil {
				store.GraphItems.PushAll(ckey, items)
			}
		} else {
			//local
			filename := g.RrdFileName(storageDir, checksum, dsType, step)
			Flush(filename, items)
			Counter += 1
		}
	}
}
