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

package g

import (
	"runtime"
)

// changelog:
// 0.0.1 init project
// 0.0.2 make mock item.Ts one step after now(), rm sending log, add flood in proc
// 0.0.3 mv common to falcon.common, simplify nodata's codes, mv cfgcenter to nodata
// 0.0.4 fix bug of nil response on collecting from query
// 0.0.5 collect items concurrently from query
// 0.0.6 clear send buffer when blocking
// 0.0.7 use gauss distribution to get threshold, sync judge and sender, fix bug of collector's cache
// 0.0.8 simplify project

const (
	VERSION = "0.0.8"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
