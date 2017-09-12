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
	"github.com/open-falcon/falcon-plus/modules/aggregator/db"
	"net/http"
)

func configProcRoutes() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		items, err := db.ReadClusterMonitorItems()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		for _, v := range items {
			w.Write([]byte(v.String()))
			w.Write([]byte("\n"))
		}
	})
}
