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

package rrdtool

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"math"
	"sync/atomic"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/rrdlite"
	"github.com/toolkits/file"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

var (
	disk_counter uint64
	net_counter  uint64
)

type fetch_t struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*cmodel.RRDData
}

type flushfile_t struct {
	filename string
	items    []*cmodel.GraphItem
}

type readfile_t struct {
	filename string
	data     []byte
}

func Start() {
	cfg := g.Config()
	var err error
	// check data dir
	if err = file.EnsureDirRW(cfg.RRD.Storage); err != nil {
		log.Fatal("rrdtool.Start error, bad data dir "+cfg.RRD.Storage+",", err)
	}

	migrate_start(cfg)
	log.Info("rrdtool migrateWorker started")

	go syncDisk()
	log.Info("rrdtool syncDiskWorker started")

	go ioWorker()
	log.Info("rrdtool ioWorker started")
}

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
	c.RRA("AVERAGE", 0, 1, RRA1PointCnt)

	// 5m一个点存2d
	c.RRA("AVERAGE", 0, 5, RRA5PointCnt)
	c.RRA("MAX", 0, 5, RRA5PointCnt)
	c.RRA("MIN", 0, 5, RRA5PointCnt)

	// 20m一个点存7d
	c.RRA("AVERAGE", 0, 20, RRA20PointCnt)
	c.RRA("MAX", 0, 20, RRA20PointCnt)
	c.RRA("MIN", 0, 20, RRA20PointCnt)

	// 3小时一个点存3个月
	c.RRA("AVERAGE", 0, 180, RRA180PointCnt)
	c.RRA("MAX", 0, 180, RRA180PointCnt)
	c.RRA("MIN", 0, 180, RRA180PointCnt)

	// 12小时一个点存1year
	c.RRA("AVERAGE", 0, 720, RRA720PointCnt)
	c.RRA("MAX", 0, 720, RRA720PointCnt)
	c.RRA("MIN", 0, 720, RRA720PointCnt)

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
func flushrrd(filename string, items []*cmodel.GraphItem) error {
	if items == nil || len(items) == 0 {
		return errors.New("empty items")
	}

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

func ReadFile(filename, md5 string) ([]byte, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_READ,
		args:   &readfile_t{filename: filename},
		done:   done,
	}

	io_task_chans[getIndex(md5)] <- task
	err := <-done
	return task.args.(*readfile_t).data, err
}

func CommitFile(filename, md5 string, items []*cmodel.GraphItem) error {
	done := make(chan error, 1)
	io_task_chans[getIndex(md5)] <- &io_task_t{
		method: IO_TASK_M_FLUSH,
		args: &flushfile_t{
			filename: filename,
			items:    items,
		},
		done: done,
	}
	atomic.AddUint64(&disk_counter, 1)
	return <-done
}

func Fetch(filename string, md5 string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_FETCH,
		args: &fetch_t{
			filename: filename,
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	io_task_chans[getIndex(md5)] <- task
	err := <-done
	return task.args.(*fetch_t).data, err
}

func fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	start_t := time.Unix(start, 0)
	end_t := time.Unix(end, 0)
	step_t := time.Duration(step) * time.Second

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

func CommitByKey(key string) error {
	md5, dsType, step, err := g.SplitRrdCacheKey(key)
	if err != nil {
		return err
	}
	filename := g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step)

	items := store.GraphItems.PopAll(key)
	if len(items) == 0 {
		return nil
	}
	return CommitFile(filename, md5, items)
}

func PullByKey(key string) {
	done := make(chan error)

	item := store.GraphItems.First(key)
	if item == nil {
		return
	}
	node, err := Consistent.Get(item.PrimaryKey())
	if err != nil {
		return
	}
	Net_task_ch[node] <- &Net_task_t{
		Method: NET_TASK_M_PULL,
		Key:    key,
		Done:   done,
	}
	// net_task slow, shouldn't block syncDisk() or CommitBeforeQuit()
	// warning: recev sigout when migrating, maybe lost memory data
	go func() {
		err := <-done
		if err != nil {
			log.Errorf("get %s %s from remote err[%s]", key, item.UUID(), err)
			return
		}
		atomic.AddUint64(&net_counter, 1)
	}()
}

func SendByKey(key string) {
	done := make(chan error)

	item := store.GraphItems.First(key)
	if item == nil {
		return
	}
	node, err := Consistent.Get(item.PrimaryKey())
	if err != nil {
		return
	}
	Net_task_ch[node] <- &Net_task_t{
		Method: NET_TASK_M_SEND,
		Key:    key,
		Done:   done,
	}

	go func() {
		err := <-done
		if err != nil {
			log.Errorf("transmit %s %s to remote err[%s]", key, item.UUID(), err)
		} else {
			log.Debugf("transmit %s %s to remote succ", key, item.UUID())
		}
	}()
}

func CommitBeforeQuit() {
	n := store.GraphItems.Size / 10
	for i := 0; i < store.GraphItems.Size; i++ {
		commitByIdxBeforeQuit(i)
		if i%n == 0 {
			log.Infof("flush rrd before quit, hash idx:%03d size:%03d disk:%08d net:%08d",
				i, store.GraphItems.Size, disk_counter, net_counter)
		}
	}
	log.Infof("flush done (disk:%08d net:%08d)", disk_counter, net_counter)
}

func commitByIdxBeforeQuit(idx int) {
	begin := time.Now()

	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}

	is_migrate := g.Config().Migrate.Enabled
	for _, key := range keys {
		flag, _ := store.GraphItems.GetFlag(key)

		if is_migrate && flag&g.GRAPH_F_MISS != 0 {
			filename, _ := getFilenameByKey(key)
			if !g.IsRrdFileExist(filename) {
				//transmit cache data to remote graph
				SendByKey(key)
			} else {
				CommitByKey(key)
			}
		} else {
			CommitByKey(key)
		}

		//check if there is backlog
		if time.Since(begin) > time.Millisecond*g.FLUSH_DISK_STEP {
			log.Warnf("commit rrd too slow, check the backlog of idx %d", idx)
		}
	}
}

func commitByIdx(idx int) {
	begin := time.Now()
	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}

	is_migrate := g.Config().Migrate.Enabled
	for _, key := range keys {
		flag, _ := store.GraphItems.GetFlag(key)
		if is_migrate {
			if flag&g.GRAPH_F_MISS == 0 && shouldFlush(key) {
				CommitByKey(key)
			}

			if flag&g.GRAPH_F_MISS != 0 {
				filename, _ := getFilenameByKey(key)
				if !g.IsRrdFileExist(filename) {
					PullByKey(key)
				} else {
					CommitByKey(key)
				}
			}
		} else if shouldFlush(key) {
			CommitByKey(key)
		}

		//check if there is backlog
		if time.Since(begin) > time.Millisecond*g.FLUSH_DISK_STEP {
			log.Warnf("commit rrd too slow, check the backlog of idx %d", idx)
		}
	}
}

func shouldFlush(key string) bool {
	if store.GraphItems.ItemCnt(key) >= g.FLUSH_MIN_COUNT {
		return true
	}

	deadline := time.Now().Unix() - int64(g.FLUSH_MAX_WAIT)
	back := store.GraphItems.Back(key)
	if back != nil && back.Timestamp <= deadline {
		return true
	}

	return false
}

func getFilenameByKey(key string) (string, error) {
	md5, dsType, step, err := g.SplitRrdCacheKey(key)
	if err != nil {
		return "", err
	}
	return g.RrdFileName(g.Config().RRD.Storage, md5, dsType, step), nil
}
