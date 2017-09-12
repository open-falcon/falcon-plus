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
	"log"
)

func QueryPlugins() (map[int][]string, error) {
	m := make(map[int][]string)

	sql := "select grp_id, dir from plugin_dir"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id  int
			dir string
		)

		err = rows.Scan(&id, &dir)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, exists := m[id]; exists {
			m[id] = append(m[id], dir)
		} else {
			m[id] = []string{dir}
		}
	}

	return m, nil
}
