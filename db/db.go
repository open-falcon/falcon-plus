package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/task/g"
	"log"
)

var DB *sql.DB

func Start() {
	var err error
	DB, err = GetDbConn()
	if err != nil {
		log.Fatalln("db:Init, get db conn fail", err)
	} else {
		log.Println("db:Init, ok")
	}
}

func GetDbConn() (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", g.Config().DB.Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(g.Config().DB.MaxIdle)

	err = conn.Ping()
	if err != nil {
		conn.Close()
	}

	return conn, err
}
