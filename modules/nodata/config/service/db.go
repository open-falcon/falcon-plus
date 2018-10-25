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

package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/Sirupsen/logrus"
	"sync"

	"github.com/open-falcon/falcon-plus/modules/nodata/g"
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
