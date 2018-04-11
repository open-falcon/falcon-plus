package string_matcher

import (
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize      = 1024000               // 102.4w
	DefaultSendTaskSleepInterval = time.Millisecond * 50 // 50ms
)

type HistoryMgr struct {
	dbConnPool   map[string]*sqlx.DB
	dbConcurrent int
	Queue        *nlist.SafeListLimited
}

type HistoryProducer struct {
	Mgr *HistoryMgr
}

type HistoryConsumer struct {
	Mgr *HistoryMgr
}

var (
	Producer *HistoryProducer
	Consumer *HistoryConsumer
)

func NewHistoryMgr() *HistoryMgr {
	mgr := new(HistoryMgr)

	cfg := g.Config()
	dbConnPool := make(map[string]*sqlx.DB)

	for name, dsn := range cfg.StringMatcher.DSN {
		conn, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			panic(err)
		}
		dbConnPool[name] = conn
	}

	mgr.dbConnPool = dbConnPool

	dbConcurrent := cfg.StringMatcher.MaxConns
	if dbConcurrent < 1 {
		dbConcurrent = 1
	}
	mgr.dbConcurrent = dbConcurrent
	mgr.Queue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)

	return mgr
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func (mgr *HistoryMgr) Append(items []interface{}) (err error) {
	conn, ok := mgr.dbConnPool["history"]
	if !ok {
		return errors.New("get Sql Connection failed")
	}

	err = conn.Ping()
	if err != nil {
		return err
	}

	tx, err := conn.Beginx()
	if err != nil {
		return err
	}

	for _, item := range items {
		meta := item.(*model.JudgeItem)
		s := "INSERT INTO history (endpoint, metric, value, counter_type, tags, Timestamp) values (?, ?, ?, ?, ?, ?); "
		tags := utils.SortedTags(meta.Tags)
		_, err = tx.Exec(s,
			meta.Endpoint,
			meta.Metric,
			truncateString(meta.ValueRaw, 1024),
			meta.JudgeType,
			tags,
			meta.Timestamp)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func InitStringMatcher() {
	mgr := NewHistoryMgr()
	Producer = new(HistoryProducer)
	Producer.Mgr = mgr
	Consumer = new(HistoryConsumer)
	Consumer.Mgr = mgr
}

func (c *HistoryConsumer) BatchDeleteHistory(before int64) error {
	conn, ok := c.Mgr.dbConnPool["history"]
	if !ok {
		return errors.New("get Sql Connection failed")
	}

	err := conn.Ping()
	if err != nil {
		return err
	}

	s := "DELETE FROM history where Timestamp < ?"
	_, err = conn.Exec(s, before)
	if err != nil {
		return err
	}
	return nil

}

func (c *HistoryConsumer) Start(batch, retry int) {
	sema := nsema.NewSemaphore(c.Mgr.dbConcurrent)

	for {
		items := c.Mgr.Queue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var err error

			for i := 0; i < retry; i++ {
				err = c.Mgr.Append(itemList)
				if err != nil {
					log.Println("SqlDbInsert failed", err)
				} else {
					//proc.SendToSqlDbCnt.IncrBy(int64(len(itemList)))
					break
				}

				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				//proc.SendToSqlDbFailCnt.IncrBy(int64(len(itemList)))
				return
			}
		}(items)
	}
}

func (p *HistoryProducer) Append(item *model.JudgeItem) bool {
	return p.Mgr.Queue.PushFront(item)
}
