package store

import (
	"database/sql"
	"fmt"
	db "github.com/open-falcon/graph/db"
	"log"
	"strconv"
	"strings"
)

func LoadEndpointId(endpoint string) (id int64, exists bool) {
	id, exists = Endpoint2Ids.EndpointId(endpoint)
	if exists {
		return
	}

	// try to load from db
	err := db.DB.QueryRow("SELECT id FROM endpoint WHERE endpoint = ?", endpoint).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		// 数据库连接出问题了
		log.Println("query endpoint id fail", err)
		return
	}

	if err == sql.ErrNoRows {
		// 是真的没查到
		return
	}

	// 肯定查到了
	exists = true
	Endpoint2Ids.Set(endpoint, id)

	return
}

func LoadDsTypeAndStep(endpointId int64, counter string) (dsType string, step int, exists bool) {
	key := fmt.Sprintf("%d-%s", endpointId, counter)
	var dsType_step string
	dsType_step, exists = Counters.Get(key)
	if exists {
		// 内存中有，太棒了，不用查询数据库了
		arr := strings.Split(dsType_step, "_")
		if len(arr) != 2 {
			exists = false
			return
		}

		var err error
		step, err = strconv.Atoi(arr[1])
		if err != nil {
			exists = false
			return
		}

		dsType = arr[0]
		return
	}

	// 内存中没有，很遗憾，需要查DB了
	err := db.DB.QueryRow("SELECT type, step FROM endpoint_counter WHERE endpoint_id = ? and counter = ?", endpointId, counter).Scan(&dsType, &step)
	if err != nil && err != sql.ErrNoRows {
		// 数据库连接出问题了
		log.Println("query type and step fail", err)
		return
	}

	if err == sql.ErrNoRows {
		// 没查到
		return
	}

	exists = true
	Counters.Set(key, fmt.Sprintf("%s_%d", dsType, step))
	return
}
