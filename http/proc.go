package http

import (
	"github.com/open-falcon/aggregator/db"
	"net/http"
)

func configProcRoutes() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ReadClusterMonitorItems()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		for _, v := range items {
			w.Write([]byte(v.String()))
			w.Write([]byte("\n"))
		}
	})
}
