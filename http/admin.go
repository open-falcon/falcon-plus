package http

import (
	"github.com/open-falcon/agent/g"
	"net/http"
	"os"
	"time"
)

func configAdminRoutes() {
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if g.InWhiteIPs(r.RemoteAddr) {
			w.Write([]byte("exiting..."))
			go func() {
				time.Sleep(time.Second)
				os.Exit(0)
			}()
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
