package rrdtool

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
	"github.com/open-falcon/rrdlite"
	"github.com/toolkits/file"
)

var Counter uint64

func Start() {
	// check data dir
	dataDir := g.Config().RRD.Storage
	if err := file.EnsureDirRW(dataDir); err != nil {
		log.Fatalln("rrdtool.Start error, bad data dir", dataDir+",", err)
	}
	log.Println("rrdtool.Start, ok")
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

func create(filename string, item *cmodel.GraphItem) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(item.Step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", item.DsType, item.Heartbeat, item.Min, item.Max)

	// 设置各种归档策略
	// 1分钟一个点存 12小时
	c.RRA("AVERAGE", 0.5, 1, 720)

	// 5m一个点存2d
	c.RRA("AVERAGE", 0.5, 5, 576)
	c.RRA("MAX", 0.5, 5, 576)
	c.RRA("MIN", 0.5, 5, 576)

	// 20m一个点存7d
	c.RRA("AVERAGE", 0.5, 20, 504)
	c.RRA("MAX", 0.5, 20, 504)
	c.RRA("MIN", 0.5, 20, 504)

	// 3小时一个点存3个月
	c.RRA("AVERAGE", 0.5, 180, 766)
	c.RRA("MAX", 0.5, 180, 766)
	c.RRA("MIN", 0.5, 180, 766)

	// 1天一个点存5year
	c.RRA("AVERAGE", 0.5, 720, 730)
	c.RRA("MAX", 0.5, 720, 730)
	c.RRA("MIN", 0.5, 720, 730)

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
	var debug_checksum string
	var debug bool

	storageDir := g.Config().RRD.Storage
	if g.Config().Debug {
		debug = true
		debug_checksum = g.Config().DebugChecksum
	} else {
		debug = false
	}

	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}

	for _, checksum := range keys {

		items := store.GraphItems.PopAll(checksum)
		size := len(items)
		if size == 0 {
			continue
		}

		first := items[0]
		filename := fmt.Sprintf("%s/%s/%s_%s_%d.rrd", storageDir, checksum[0:2], checksum, first.DsType, first.Step)
		if debug && debug_checksum == checksum {
			for _, item := range items {
				log.Printf(
					"2-flush:%d:%s:%lf",
					item.Timestamp,
					time.Unix(item.Timestamp, 0).Format("2006-01-02 15:04:05"),
					item.Value,
				)
			}
		}

		err := Flush(filename, items)
		if err != nil && debug && debug_checksum == checksum {
			log.Println("flush fail:", err, "filename:", filename)
		}
		Counter += 1
	}
	if debug {
		log.Println("flushrrd counter:", Counter)
	}
}
