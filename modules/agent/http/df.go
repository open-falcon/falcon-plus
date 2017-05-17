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
	"fmt"
	"github.com/toolkits/core"
	"github.com/toolkits/nux"
	"net/http"
)

func configDfRoutes() {
	http.HandleFunc("/page/df", func(w http.ResponseWriter, r *http.Request) {
		mountPoints, err := nux.ListMountPoint()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		var ret [][]interface{} = make([][]interface{}, 0)
		for idx := range mountPoints {
			var du *nux.DeviceUsage
			du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
			if err == nil {
				ret = append(ret,
					[]interface{}{
						du.FsSpec,
						core.ReadableSize(float64(du.BlocksAll)),
						core.ReadableSize(float64(du.BlocksUsed)),
						core.ReadableSize(float64(du.BlocksFree)),
						fmt.Sprintf("%.1f%%", du.BlocksUsedPercent),
						du.FsFile,
						core.ReadableSize(float64(du.InodesAll)),
						core.ReadableSize(float64(du.InodesUsed)),
						core.ReadableSize(float64(du.InodesFree)),
						fmt.Sprintf("%.1f%%", du.InodesUsedPercent),
						du.FsVfstype,
					})
			}
		}

		RenderDataJson(w, ret)
	})
}
