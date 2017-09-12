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

package db

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"log"
)

func ReadClusterMonitorItems() (M map[string]*g.Cluster, err error) {
	M = make(map[string]*g.Cluster)
	sql := "SELECT `id`, `grp_id`, `numerator`, `denominator`, `endpoint`, `metric`, `tags`, `ds_type`, `step`, `last_update` FROM `cluster`"

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
		err = rows.Scan(&c.Id, &c.GroupId, &c.Numerator, &c.Denominator, &c.Endpoint, &c.Metric, &c.Tags, &c.DsType, &c.Step, &c.LastUpdate)
		if err != nil {
			log.Println("[E]", err)
			continue
		}

		M[fmt.Sprintf("%d%v", c.Id, c.LastUpdate)] = &c
	}

	return M, err
}
