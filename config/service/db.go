package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"

	"github.com/open-falcon/nodata/g"
)

const (
	dbBaseConnName = "db.base"
)

var (
	dbLock    = sync.RWMutex{}
	dbConnMap = make(map[string]*sql.DB)
)

func InitDB() {
	_, err := GetDbConn(dbBaseConnName)
	if err != nil {
		log.Fatalln("config.InitDB error", err)
		return // never go here
	}

	log.Println("config.InitDB ok")
}

func GetBaseConn() (c *sql.DB, e error) {
	return GetDbConn(dbBaseConnName)
}

func GetDbConn(connName string) (c *sql.DB, e error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	var err error
	var dbConn *sql.DB
	dbConn = dbConnMap[connName]
	if dbConn == nil {
		dbConn, err = makeDbConn()
		if err != nil {
			closeDbConn(dbConn)
			return nil, err
		}
		dbConnMap[connName] = dbConn
	}

	err = dbConn.Ping()
	if err != nil {
		closeDbConn(dbConn)
		delete(dbConnMap, connName)
		return nil, err
	}

	return dbConn, err
}

// internal
func makeDbConn() (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", g.Config().Config.Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(int(g.Config().Config.MaxIdle))
	err = conn.Ping()

	return conn, err
}

func closeDbConn(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}
