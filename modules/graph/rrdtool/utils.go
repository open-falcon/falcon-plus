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

package rrdtool

import (
	"strconv"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

// getIndex
// 输入: md5
// 输出: slice的index
func getIndex(md5 string) (index int) {
	batchNum := g.Config().IOWorkerNum
	firstBytesSize := g.Config().FirstBytesSize

	if batchNum <= 1 || len(md5) < firstBytesSize || firstBytesSize == 0 {
		return 0
	}

	m, err := strconv.ParseInt(md5[0:firstBytesSize], 16, 64)
	if err != nil {
		return 0
	}

	return int(m) % int(batchNum)
}
