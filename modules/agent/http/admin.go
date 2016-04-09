package http

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/file"
	"net/http"
	"os"
	"time"
)

func configAdminRoutes() {
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if g.IsTrustable(r.RemoteAddr) {
			w.Write([]byte("exiting..."))
			go func() {
				time.Sleep(time.Second)
				os.Exit(0)
			}()
		} else {
			w.Write([]byte("no privilege"))
		}
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if g.IsTrustable(r.RemoteAddr) {
			g.ParseConfig(g.ConfigFile)
			RenderDataJson(w, g.Config())
		} else {
			w.Write([]byte("no privilege"))
		}
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, file.SelfDir())
	})

	http.HandleFunc("/ips", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.TrustableIps())
	})
}
