package http

import (
	"fmt"
	"github.com/open-falcon/query/g"
	"github.com/toolkits/file"
	"net/http"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.VERSION))
	})

	http.HandleFunc("/versiongit", func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf("%s %s", g.VERSION, g.COMMIT)
		w.Write([]byte(s))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, file.SelfDir())
	})
}
