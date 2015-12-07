package db

import (
	"fmt"
	"github.com/open-falcon/aggregator/g"
	"log"
)

func ReadClusterMonitorItems() (M map[string]*g.Cluster, err error) {
	M = make(map[string]*g.Cluster)
	sql := "SELECT `id`, `node`, `numerator`, `denominator`, `endpoint`, `metric`, `tags`, `ds_type`, `step`, `last_update` FROM `cluster`"

	cfg := g.Config()
	ids := cfg.Database.Ids
	if len(ids) != 2 {
		log.Fatalln("ids configuration error")
	}

	if ids[0] != -1 && ids[1] != -1 {
		sql = fmt.Sprintf("%s WHERE `id` >= %d and `id` <= %d", sql, ids[0], ids[1])
	} else {
		if ids[0] != -1 {
			sql = fmt.Sprintf("%s WHERE `id` >= %d", sql, ids[0])
		}

		if ids[1] != -1 {
			sql = fmt.Sprintf("%s WHERE `id` <= %d", sql, ids[1])
		}
	}

	if cfg.Debug {
		log.Println(sql)
	}

	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("[E]", err)
		return M, err
	}

	defer rows.Close()
	for rows.Next() {
		var c g.Cluster
		err = rows.Scan(&c.Id, &c.Node, &c.Numerator, &c.Denominator, &c.Endpoint, &c.Metric, &c.Tags, &c.DsType, &c.Step, &c.LastUpdate)
		if err != nil {
			log.Println("[E]", err)
			continue
		}

		M[fmt.Sprintf("%d%v", c.Id, c.LastUpdate)] = &c
	}

	return M, err
}
