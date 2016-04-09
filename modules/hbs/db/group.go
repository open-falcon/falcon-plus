package db

import (
	"log"
)

func QueryHostGroups() (map[int][]int, error) {
	m := make(map[int][]int)

	sql := "select grp_id, host_id from grp_host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var gid, hid int
		err = rows.Scan(&gid, &hid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, exists := m[hid]; exists {
			m[hid] = append(m[hid], gid)
		} else {
			m[hid] = []int{gid}
		}
	}

	return m, nil
}
