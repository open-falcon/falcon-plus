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
