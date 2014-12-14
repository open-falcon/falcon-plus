package http

import (
	"github.com/open-falcon/agent/g"
	ff "github.com/toolkits/file"
	"net/http"
	"path/filepath"
	"strings"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			if !ff.IsExist(filepath.Join(g.Root, "/public", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
	})

}
