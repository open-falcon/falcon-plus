package goperfcounter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func startHttp(addr string, debug bool) {
	configCommonRoutes()
	configProcRoutes()
	if len(addr) >= 9 { //x.x.x.x:x
		s := &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 30,
		}
		go func() {
			if debug {
				log.Println("[perfcounter] http server start, listening on", addr)
			}
			s.ListenAndServe()
			if debug {
				log.Println("[perfcounter] http server stop,", addr)
			}
		}()
	}
}

// routers
func configProcRoutes() {
	http.HandleFunc("/pfc/proc/metrics/json", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, rawMetrics())
	})
	http.HandleFunc("/pfc/proc/metrics/falcon", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, falconMetrics())
	})
	// url=/pfc/proc/metric/{json,falcon}
	http.HandleFunc("/pfc/proc/metrics/", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		urlParam := r.URL.Path[len("/pfc/proc/metrics/"):]
		args := strings.Split(urlParam, "/")
		argsLen := len(args)
		if argsLen != 2 {
			RenderJson(w, "")
			return
		}

		types := []string{}
		typeslice := strings.Split(args[0], ",")
		for _, t := range typeslice {
			nt := strings.TrimSpace(t)
			if nt != "" {
				types = append(types, nt)
			}
		}

		if args[1] == "json" {
			RenderJson(w, rawMetric(types))
			return
		}
		if args[1] == "falcon" {
			RenderJson(w, falconMetric(types))
			return
		}
	})

	http.HandleFunc("/pfc/proc/metrics/size", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, rawSizes())
	})

}

func configCommonRoutes() {
	http.HandleFunc("/pfc/health", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/pfc/version", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		w.Write([]byte(fmt.Sprintf("%s\n", VERSION)))
	})

	http.HandleFunc("/pfc/config", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		RenderJson(w, config())
	})

	http.HandleFunc("/pfc/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if !isLocalReq(r.RemoteAddr) {
			RenderJson(w, "no privilege")
			return
		}
		loadConfig()
		RenderJson(w, "ok")
	})
}

func isLocalReq(raddr string) bool {
	if strings.HasPrefix(raddr, "127.0.0.1") {
		return true
	}
	return false
}

// render
func RenderJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, Dto{Msg: "success", Data: data})
}

func RenderString(w http.ResponseWriter, msg string) {
	renderJson(w, map[string]string{"msg": msg})
}

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

// common http return
type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
