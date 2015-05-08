package http

import (
	"fmt"
	"github.com/open-falcon/transfer/g"
	"github.com/toolkits/file"
	"net/http"
	"strings"
)

func configCommonRoutes() {
	// GET
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.Config())
	})

	// SET
	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			w.Write([]byte("not supported\n")) // 暂时不支持配置的重新加载
		} else {
			w.Write([]byte("no privilege\n"))
		}
	})

}
