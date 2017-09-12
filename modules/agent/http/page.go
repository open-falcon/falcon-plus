// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"net/http"
	"path/filepath"
	"strings"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			if !file.IsExist(filepath.Join(g.Root, "/public", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
	})

}
