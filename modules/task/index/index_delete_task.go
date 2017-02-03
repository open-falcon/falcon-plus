package index

import (
	"log"
	"time"

	Mdb "github.com/open-falcon/falcon-plus/common/db"
	cron "github.com/toolkits/cron"
	ntime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/task/proc"
)

const (
	indexDeleteCronSpec = "0 0 2 ? * 6" // 每周6晚上22:00执行一次
	deteleStepInSec     = 7 * 24 * 3600 // 索引的最大生存周期, sec
)

var (
	indexDeleteCron = cron.New()
)

// 启动 索引全量更新 定时任务
func StartIndexDeleteTask() {
	indexDeleteCron.AddFuncCC(indexDeleteCronSpec, func() { DeleteIndex() }, 1)
	indexDeleteCron.Start()
}

// 索引的全量更新
func DeleteIndex() {
	startTs := time.Now().Unix()
	deleteIndex()
	endTs := time.Now().Unix()
	log.Printf("deleteIndex, start %s, ts %ds", ntime.FormatTs(startTs), endTs-startTs)

	// statistics
	proc.IndexDeleteCnt.Incr()
}

// 先select 得到可能被删除的index的信息, 然后以相同的条件delete. select和delete不是原子操作,可能有一些不一致,但不影响正确性
func deleteIndex() error {
	dbConn, err := GetDbConn()
	if err != nil {
		log.Println("[ERROR] get dbConn fail", err)
		return err
	}
	defer dbConn.Close()

	ts := time.Now().Unix()
	lastTs := ts - deteleStepInSec
	log.Printf("deleteIndex, lastTs %d\n", lastTs)

	// reinit statistics
	proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", 0)

	// endpoint表
	{
		// select
		rows, err := dbConn.Query("SELECT id, endpoint FROM endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		cnt := 0
		for rows.Next() {
			item := &Mdb.GraphEndpoint{}
			err := rows.Scan(&item.Id, &item.Endpoint)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("will delete endpoint:", item)
			cnt++
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			return err
		}

		// delete
		_, err = dbConn.Exec("DELETE FROM endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("delete endpoint, done, cnt %d\n", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", cnt)
	}

	// tag_endpoint表
	{
		// select
		rows, err := dbConn.Query("SELECT id, tag, endpoint_id FROM tag_endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		cnt := 0
		for rows.Next() {
			item := &Mdb.GraphTagEndpoint{}
			err := rows.Scan(&item.Id, &item.Tag, &item.EndpointId)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("will delete tag_endpoint:", item)
			cnt++
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			return err
		}

		// delete
		_, err = dbConn.Exec("DELETE FROM tag_endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("delete tag_endpoint, done, cnt %d\n", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", cnt)
	}
	// endpoint_counter表
	{
		// select
		rows, err := dbConn.Query("SELECT id, endpoint_id, counter FROM endpoint_counter WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		cnt := 0
		for rows.Next() {
			item := &Mdb.GraphEndpointCounter{}
			err := rows.Scan(&item.Id, &item.EndpointId, &item.Counter)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("will delete endpoint_counter:", item)
			cnt++
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			return err
		}

		// delete
		_, err = dbConn.Exec("DELETE FROM endpoint_counter WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("delete endpoint_counter, done, cnt %d\n", cnt)

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", cnt)
	}

	return nil
}
