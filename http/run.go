package http

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/sys"
	"io/ioutil"
	"net/http"
)

func configRunRoutes() {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if g.InWhiteIPs(r.RemoteAddr) {
			if r.ContentLength == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("req.Body is blank"))
				return
			}

			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("read req.Body fail: " + err.Error()))
				return
			}

			body := string(bs)
			out, err := sys.CmdOutBytes("/bin/bash", "-c", body)
			if err != nil {
				w.Write([]byte("exec fail: " + err.Error()))
				return
			}

			w.Write(out)
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
