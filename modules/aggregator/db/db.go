package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/aggregator/g"
	"log"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxIdleConns(g.Config().Database.Idle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
