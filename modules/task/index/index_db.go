package index

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"github.com/open-falcon/falcon-plus/modules/task/g"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = GetDbConn()
	if err != nil {
		log.Fatalln("index:InitDB error,", err)
	} else {
		log.Println("index:InitDB ok")
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
