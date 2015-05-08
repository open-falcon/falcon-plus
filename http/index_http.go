package http

import (
	"github.com/open-falcon/task/index"
	"net/http"
)

func configIndexHttpRoutes() {
	http.HandleFunc("/index/delete", func(w http.ResponseWriter, r *http.Request) {
		index.DeleteIndex()
		RenderDataJson(w, "ok")
	})
	http.HandleFunc("/index/updateAll", func(w http.ResponseWriter, r *http.Request) {
		index.UpdateAllIndex()
		RenderDataJson(w, "ok")
	})
}
