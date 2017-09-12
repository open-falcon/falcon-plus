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

package index

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"time"

	nsema "github.com/toolkits/concurrent/semaphore"
	ntime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	proc "github.com/open-falcon/falcon-plus/modules/graph/proc"
)

const (
	IndexUpdateIncrTaskSleepInterval = time.Duration(1) * time.Second // 增量更新间隔时间, 默认30s
)

var (
	semaUpdateIndexIncr = nsema.NewSemaphore(2) // 索引增量更新时操作mysql的并发控制
)

// 启动索引的 异步、增量更新 任务, 每隔一定时间，刷新cache中的数据到数据库中
func StartIndexUpdateIncrTask() {
	for {
		time.Sleep(IndexUpdateIncrTaskSleepInterval)
		startTs := time.Now().Unix()
		cnt := updateIndexIncr()
		endTs := time.Now().Unix()
		// statistics
		proc.IndexUpdateIncrCnt.SetCnt(int64(cnt))
		proc.IndexUpdateIncr.Incr()
		proc.IndexUpdateIncr.PutOther("lastStartTs", ntime.FormatTs(startTs))
		proc.IndexUpdateIncr.PutOther("lastTimeConsumingInSec", endTs-startTs)
	}
}

// 进行一次增量更新
func updateIndexIncr() int {
	ret := 0
	if unIndexedItemCache == nil || unIndexedItemCache.Size() <= 0 {
		return ret
	}

	dbConn, err := g.GetDbConn("UpdateIndexIncrTask")
	if err != nil {
		log.Error("[ERROR] get dbConn fail", err)
		return ret
	}

	keys := unIndexedItemCache.Keys()
	for _, key := range keys {
		icitem := unIndexedItemCache.Get(key)
		unIndexedItemCache.Remove(key)
		if icitem != nil {
			// 并发更新mysql
			semaUpdateIndexIncr.Acquire()
			go func(key string, icitem *IndexCacheItem, dbConn *sql.DB) {
				defer semaUpdateIndexIncr.Release()
				err := updateIndexFromOneItem(icitem.Item, dbConn)
				if err != nil {
					proc.IndexUpdateIncrErrorCnt.Incr()
				} else {
					IndexedItemCache.Put(key, icitem)
				}
			}(key, icitem.(*IndexCacheItem), dbConn)
			ret++
		}
	}

	return ret
}
