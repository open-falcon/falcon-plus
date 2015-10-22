package service

import (
	"fmt"
	"log"
	"time"
)

// FIX ME: too many JOIN
func GetHostsFromGroup(grpName string) map[string]int {
	hosts := make(map[string]int)

	now := time.Now().Unix()
	q := fmt.Sprintf("SELECT host.id, host.hostname FROM grp_host AS gh "+
		" INNER JOIN host ON host.id=gh.host_id AND (host.maintain_begin > %d OR host.maintain_end < %d)"+
		" INNER JOIN grp ON grp.id=gh.grp_id AND grp.grp_name='%s'", now, now, grpName)

	dbConn, err := GetDbConn("nodata.host")
	if err != nil {
		log.Println("db.get_conn error, host", err)
		return hosts
	}

	rows, err := dbConn.Query(q)
	if err != nil {
		log.Println("[ERROR]", err)
		return hosts
	}

	defer rows.Close()
	for rows.Next() {
		hid := -1
		hostname := ""
		err = rows.Scan(&hid, &hostname)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}
		if hid < 0 || hostname == "" {
			continue
		}

		hosts[hostname] = hid
	}

	return hosts
}
