package index

import (
	cron "github.com/niean/cron"
	Mdb "github.com/open-falcon/model/db"
	"github.com/open-falcon/task/db"
	"github.com/open-falcon/task/proc"
	TSemaphore "github.com/toolkits/concurrent/semaphore"
	"log"
	"time"
)

const (
	indexDeleteCronSpec = "0 20 * * * ?" //"0 0 3 * * ?" //索引垃圾清理的cron周期描述
	deteleStepInSec     = 2 * 24 * 3600  // 索引的最大生存周期, sec
)

var (
	semaIndexDelete = TSemaphore.NewSemaphore(1)
	indexDeleteCron = cron.New()
)

// 启动 索引全量更新 定时任务
func StartIndexDeleteTask() {
	indexDeleteCron.AddFunc(indexDeleteCronSpec, func() {
		DeleteIndex()
	})
	indexDeleteCron.Start()
}

func StopIndexDeleteTask() {
	indexDeleteCron.Stop()
}

// 索引的全量更新
func DeleteIndex() {
	// 阻止多个并发的访问,高并发时可能无效
	if semaIndexDelete.AvailablePermits() <= 0 {
		log.Printf("deleteIndex, concurrent not avaiable")
		return
	}

	semaIndexDelete.Acquire()
	defer semaIndexDelete.Release()

	startTs := time.Now().Unix()
	deleteIndex()
	endTs := time.Now().Unix()
	log.Printf("deleteIndex, startTs %s, time-consuming %d sec\n", proc.FmtUnixTs(startTs), endTs-startTs)

	// statistics
	proc.IndexDelete.Incr()
	proc.IndexDelete.PutOther("lastStartTs", proc.FmtUnixTs(startTs))
	proc.IndexDelete.PutOther("lastTimeConsumingInSec", endTs-startTs)
}

// 先select 得到可能被删除的index的信息, 然后以相同的条件delete. select和delete不是原子操作,可能有一些不一致,但不影响正确性
func deleteIndex() error {
	dbConn, err := db.GetDbConn()
	if err != nil {
		log.Println("[ERROR] get dbConn fail", err)
		return err
	}
	defer dbConn.Close()

	ts := time.Now().Unix()
	lastTs := ts - deteleStepInSec
	log.Printf("deleteIndex, lastTs %d\n", lastTs)

	// reinit statistics
	// TODO 侵入性有点强阿, 改下这里
	proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", 0)
	proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", 0)

	// endpoint表
	{
		// select
		endpoints := make([]*Mdb.GraphEndpoint, 0)
		rows, err := dbConn.Query("SELECT id, endpoint FROM endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		for rows.Next() {
			item := &Mdb.GraphEndpoint{}
			err := rows.Scan(&item.Id, &item.Endpoint)
			if err != nil {
				log.Println(err)
				return err
			}
			endpoints = append(endpoints, item)
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

		for _, item := range endpoints {
			log.Println("delete endpoint:", item)
		}
		log.Printf("delete endpoint, delete cnt %d\n", len(endpoints))

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpoint", len(endpoints))
	}

	// tag_endpoint表
	{
		// select
		tagEndpoints := make([]*Mdb.GraphTagEndpoint, 0)
		rows, err := dbConn.Query("SELECT id, tag, endpoint_id FROM tag_endpoint WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		for rows.Next() {
			item := &Mdb.GraphTagEndpoint{}
			err := rows.Scan(&item.Id, &item.Tag, &item.EndpointId)
			if err != nil {
				log.Println(err)
				return err
			}
			tagEndpoints = append(tagEndpoints, item)
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

		for _, item := range tagEndpoints {
			log.Println("delete tag_endpoint:", item)
		}
		log.Printf("delete tag_endpoint, delete cnt %d\n", len(tagEndpoints))

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntTagEndpoint", len(tagEndpoints))
	}
	// endpoint_counter表
	{
		// select
		endpointCounters := make([]*Mdb.GraphEndpointCounter, 0)
		rows, err := dbConn.Query("SELECT id, endpoint_id, counter FROM endpoint_counter WHERE ts < ?", lastTs)
		if err != nil {
			log.Println(err)
			return err
		}

		for rows.Next() {
			item := &Mdb.GraphEndpointCounter{}
			err := rows.Scan(&item.Id, &item.EndpointId, &item.Counter)
			if err != nil {
				log.Println(err)
				return err
			}
			endpointCounters = append(endpointCounters, item)
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

		for _, item := range endpointCounters {
			log.Println("delete endpoint_counter:", item)
		}
		log.Printf("delete endpoint_counter, delete cnt %d\n", len(endpointCounters))

		// statistics
		proc.IndexDeleteCnt.PutOther("deleteCntEndpointCounter", len(endpointCounters))
	}

	return nil
}
