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

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/sys"
	"log"
	"strconv"
	"strings"
)

func DuMetrics() (L []*model.MetricValue) {
	paths := g.DuPaths()
	for _, path := range paths {
		out, err := sys.CmdOutNoLn("du", "-bs", path)
		if err != nil {
			log.Println("du -bs", path, "fail", err)
			continue
		}

		arr := strings.Fields(out)
		if len(arr) == 1 {
			continue
		}

		size, err := strconv.ParseUint(arr[0], 10, 64)
		if err != nil {
			log.Println("cannot parse du -bs", path, "output")
			continue
		}

		L = append(L, GaugeValue(g.DU_BS, size, "path="+path))
	}

	return
}
