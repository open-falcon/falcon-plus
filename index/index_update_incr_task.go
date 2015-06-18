package index

import (
	"database/sql"
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/g"
	proc "github.com/open-falcon/graph/proc"
	nsema "github.com/toolkits/concurrent/semaphore"
	ntime "github.com/toolkits/time"
	"log"
	"time"
)

const (
	IndexUpdateIncrTaskSleepInterval = time.Duration(1) * time.Second // 增量更新间隔时间, 默认30s
)

var (
	semaUpdateIndexIncr = nsema.NewSemaphore(2) // 索引增量更新时操作mysql的并发控制
)

// 启动索引的 异步、增量更新 任务
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
		log.Println("[ERROR] get dbConn fail", err)
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
				err := maybeUpdateIndexFromOneItem(icitem.Item, dbConn)
				if err != nil {
					proc.IndexUpdateIncrErrorCnt.Incr()
				} else {
					indexedItemCache.Put(key, icitem)
				}
			}(key, icitem.(*IndexCacheItem), dbConn)
			ret++
		}
	}

	return ret
}

//
func maybeUpdateIndexFromOneItem(item *cmodel.GraphItem, conn *sql.DB) error {
	if item == nil {
		return nil
	}

	endpoint := item.Endpoint
	ts := item.Timestamp
	var endpointId int64 = -1
	sqlDuplicateString := " ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)" //第一个字符是空格

	// endpoint表
	{
		err := conn.QueryRow("SELECT id FROM endpoint WHERE endpoint = ?", endpoint).Scan(&endpointId)
		if err != nil && err != sql.ErrNoRows {
			log.Println(endpoint, err)
			return err
		}
		proc.IndexUpdateIncrDbEndpointSelectCnt.Incr()

		if err == sql.ErrNoRows || endpointId <= 0 { // 数据库中也没有, insert
			sqlStr := "INSERT INTO endpoint (endpoint, ts, t_create) VALUES (?, ?, now())" + sqlDuplicateString
			ret, err := conn.Exec(sqlStr, endpoint, ts)
			if err != nil {
				log.Println(err)
				return err
			}
			proc.IndexUpdateIncrDbEndpointInsertCnt.Incr()

			endpointId, err = ret.LastInsertId()
			if err != nil {
				log.Println(err)
				return err
			}
		} else { // do not update
		}
		// 更新缓存
		//dbEndpointCache.Set(endpoint, endpointId, 0)
	}

	// tag_endpoint表
	{
		for tagKey, tagVal := range item.Tags {
			tag := fmt.Sprintf("%s=%s", tagKey, tagVal)

			var tagEndpointId int64 = -1
			err := conn.QueryRow("SELECT id FROM tag_endpoint WHERE tag = ? and endpoint_id = ?",
				tag, endpointId).Scan(&tagEndpointId)
			if err != nil && err != sql.ErrNoRows {
				log.Println(tag, endpointId, err)
				return err
			}
			proc.IndexUpdateIncrDbTagEndpointSelectCnt.Incr()

			if err == sql.ErrNoRows || tagEndpointId <= 0 {
				sqlStr := "INSERT INTO tag_endpoint (tag, endpoint_id, ts, t_create) VALUES (?, ?, ?, now())" + sqlDuplicateString
				ret, err := conn.Exec(sqlStr, tag, endpointId, ts)
				if err != nil {
					log.Println(err)
					return err
				}
				proc.IndexUpdateIncrDbTagEndpointInsertCnt.Incr()

				tagEndpointId, err = ret.LastInsertId()
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}

	// endpoint_counter表
	{
		counter := item.Metric
		if len(item.Tags) > 0 {
			counter = fmt.Sprintf("%s/%s", counter, cutils.SortedTags(item.Tags))
		}

		var endpointCounterId int64 = -1
		var step int = 0
		var dstype string = "nil"

		err := conn.QueryRow("SELECT id,step,type FROM endpoint_counter WHERE endpoint_id = ? and counter = ?",
			endpointId, counter).Scan(&endpointCounterId, &step, &dstype)
		if err != nil && err != sql.ErrNoRows {
			log.Println(counter, endpointId, err)
			return err
		}
		proc.IndexUpdateIncrDbEndpointCounterSelectCnt.Incr()

		if err == sql.ErrNoRows || endpointCounterId <= 0 {
			sqlStr := "INSERT INTO endpoint_counter (endpoint_id,counter,step,type,ts,t_create) VALUES (?,?,?,?,?,now())" +
				" ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts), step=VALUES(step),type=VALUES(type)"
			ret, err := conn.Exec(sqlStr, endpointId, counter, item.Step, item.DsType, ts)
			if err != nil {
				log.Println(err)
				return err
			}
			proc.IndexUpdateIncrDbEndpointCounterInsertCnt.Incr()

			endpointCounterId, err = ret.LastInsertId()
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			if !(item.Step == step && item.DsType == dstype) {
				_, err := conn.Exec("UPDATE endpoint_counter SET step = ?, type = ? where id = ?",
					item.Step, item.DsType, endpointCounterId)
				proc.IndexUpdateIncrDbEndpointCounterUpdateCnt.Incr()
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}

	return nil
}
