package http

import (
	"github.com/falcon/task/g"
	"github.com/falcon/task/proc"
	"net/http"
)

func configProcHttpRoutes() {
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	http.HandleFunc("/statistics/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.Config())
	})
}
