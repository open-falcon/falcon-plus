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
	"github.com/toolkits/sys"
	"io/ioutil"
	"net/http"
)

func configRunRoutes() {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Http.Backdoor {
			w.Write([]byte("/run disabled"))
			return
		}

		if g.IsTrustable(r.RemoteAddr) {
			if r.ContentLength == 0 {
				http.Error(w, "body is blank", http.StatusBadRequest)
				return
			}

			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			body := string(bs)
			out, err := sys.CmdOutBytes("sh", "-c", body)
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
