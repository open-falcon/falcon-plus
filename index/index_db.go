package index

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/task/g"
	"log"
)

var DB *sql.DB

func StartDB() {
	var err error
	DB, err = GetDbConn()
	if err != nil {
		log.Fatalln("db:Init, get db conn fail", err)
	} else {
		log.Println("db:Init, ok")
	}
}

func GetDbConn() (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", g.Config().Index.Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(g.Config().Index.MaxIdle)

	err = conn.Ping()
	if err != nil {
		conn.Close()
	}

	return conn, err
}
