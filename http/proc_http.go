package http

import (
	"github.com/open-falcon/task/proc"
	"net/http"
)

func configProcHttpRoutes() {
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})
}
