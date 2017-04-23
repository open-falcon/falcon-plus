package index

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	nsema "github.com/toolkits/concurrent/semaphore"
	ntime "github.com/toolkits/time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	proc "github.com/open-falcon/falcon-plus/modules/graph/proc"
)

const (
	DefaultUpdateStepInSec     = 2 * 24 * 3600 //更新步长,一定不能大于删除步长. 两天内的数据,都可以用来建立索引
	ConcurrentOfUpdateIndexAll = 1
)

var (
	semaIndexUpdateAllTask = nsema.NewSemaphore(ConcurrentOfUpdateIndexAll) //全量同步任务 并发控制器
	semaIndexUpdateAll     = nsema.NewSemaphore(4)                          // 索引全量更新时的mysql操作并发控制
)

// 更新一条监控数据对应的索引. 用于手动添加索引,一般情况下不会使用
func UpdateIndexOne(endpoint string, metric string, tags map[string]string, dstype string, step int) error {
	itemDemo := &cmodel.GraphItem{
		Endpoint: endpoint,
		Metric:   metric,
		Tags:     tags,
		DsType:   dstype,
		Step:     step,
	}
	md5 := itemDemo.Checksum()
	uuid := itemDemo.UUID()

	cached := IndexedItemCache.Get(md5)
	if cached == nil {
		return fmt.Errorf("not found")
	}

	icitem := cached.(*IndexCacheItem)
	if icitem.UUID != uuid {
		return fmt.Errorf("bad type or step")
	}
	gitem := icitem.Item

	dbConn, err := g.GetDbConn("UpdateIndexIncrTask")
	if err != nil {
		log.Println("[ERROR] make dbConn fail", err)
		return err
	}

	return updateIndexFromOneItem(gitem, dbConn)
}

// 索引全量更新的当前并行数
func GetConcurrentOfUpdateIndexAll() int {
	return ConcurrentOfUpdateIndexAll - semaIndexUpdateAllTask.AvailablePermits()
}

// 索引的全量更新
func UpdateIndexAllByDefaultStep() {
	UpdateIndexAll(DefaultUpdateStepInSec)
}
func UpdateIndexAll(updateStepInSec int64) {
	// 减少任务积压,但高并发时可能无效(AvailablePermits不是线程安全的)
	if semaIndexUpdateAllTask.AvailablePermits() <= 0 {
		log.Println("updateIndexAll, concurrent not avaiable")
		return
	}

	semaIndexUpdateAllTask.Acquire()
	defer semaIndexUpdateAllTask.Release()

	startTs := time.Now().Unix()
	cnt := updateIndexAll(updateStepInSec)
	endTs := time.Now().Unix()
	log.Printf("UpdateIndexAll, lastStartTs %s, updateStepInSec %d, lastTimeConsumingInSec %d\n",
		ntime.FormatTs(startTs), updateStepInSec, endTs-startTs)

	// statistics
	proc.IndexUpdateAllCnt.SetCnt(int64(cnt))
	proc.IndexUpdateAll.Incr()
	proc.IndexUpdateAll.PutOther("lastStartTs", ntime.FormatTs(startTs))
	proc.IndexUpdateAll.PutOther("updateStepInSec", updateStepInSec)
	proc.IndexUpdateAll.PutOther("lastTimeConsumingInSec", endTs-startTs)
	proc.IndexUpdateAll.PutOther("updateCnt", cnt)
}

func updateIndexAll(updateStepInSec int64) int {
	var ret int = 0
	if IndexedItemCache == nil || IndexedItemCache.Size() <= 0 {
		return ret
	}

	dbConn, err := g.GetDbConn("UpdateIndexIncrTask")
	if err != nil {
		log.Println("[ERROR] make dbConn fail", err)
		return ret
	}

	// lastTs for update index
	ts := time.Now().Unix()
	lastTs := ts - updateStepInSec

	keys := IndexedItemCache.Keys()
	for _, key := range keys {
		icitem := IndexedItemCache.Get(key)
		if icitem == nil {
			continue
		}

		gitem := icitem.(*IndexCacheItem).Item
		if gitem.Timestamp < lastTs { //缓存中的数据太旧了,不能用于索引的全量更新
			IndexedItemCache.Remove(key) //在这里做个删除,有点恶心
			continue
		}
		// 并发写mysql
		semaIndexUpdateAll.Acquire()
		go func(gitem *cmodel.GraphItem, dbConn *sql.DB) {
			defer semaIndexUpdateAll.Release()
			err := updateIndexFromOneItem(gitem, dbConn)
			if err != nil {
				proc.IndexUpdateAllErrorCnt.Incr()
			}
		}(gitem, dbConn)

		ret++
	}

	return ret
}

// 根据item,更新db存储. 不用本地缓存 优化db访问频率.
func updateIndexFromOneItem(item *cmodel.GraphItem, conn *sql.DB) error {
	if item == nil {
		return nil
	}

	endpoint := item.Endpoint
	ts := item.Timestamp
	var endpointId int64 = -1
	sqlDuplicateString := " ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)" //第一个字符是空格

	// endpoint表
	{
		sqlStr := "INSERT INTO endpoint(endpoint, ts, t_create) VALUES (?, ?, now())" + sqlDuplicateString
		ret, err := conn.Exec(sqlStr, endpoint, ts)
		if err != nil {
			log.Println(err)
			return err
		}

		endpointId, err = ret.LastInsertId()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// tag_endpoint表
	{
		sqlStr := "INSERT INTO tag_endpoint(tag, endpoint_id, ts, t_create) VALUES (?, ?, ?, now())" + sqlDuplicateString
		for tagKey, tagVal := range item.Tags {
			tag := fmt.Sprintf("%s=%s", tagKey, tagVal)

			ret, err := conn.Exec(sqlStr, tag, endpointId, ts)
			if err != nil {
				log.Println(err)
				return err
			}

			_, err = ret.LastInsertId()
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	// endpoint_counter表
	{
		counter := item.Metric
		if len(item.Tags) > 0 {
			counter = fmt.Sprintf("%s/%s", counter, cutils.SortedTags(item.Tags))
		}

		sqlStr := "INSERT INTO endpoint_counter(endpoint_id,counter,step,type,ts,t_create) VALUES (?,?,?,?,?,now())" + sqlDuplicateString
		ret, err := conn.Exec(sqlStr, endpointId, counter, item.Step, item.DsType, ts)
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = ret.LastInsertId()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
