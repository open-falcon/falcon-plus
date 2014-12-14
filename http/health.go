package http

import (
	"github.com/open-falcon/agent/g"
	"net/http"
)

func initHealthRoutes() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(g.VERSION))
}
