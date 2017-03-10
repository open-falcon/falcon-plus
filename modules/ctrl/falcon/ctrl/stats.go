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
package ctrl

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

const (
	ST_UPSTREAM_RECONNECT = iota
	ST_UPSTREAM_DIAL
	ST_UPSTREAM_DIAL_ERR
	ST_UPSTREAM_UPDATE
	ST_UPSTREAM_UPDATE_ITEM
	ST_UPSTREAM_UPDATE_ERR
	ST_ARRAY_SIZE
)

var (
	statName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"ST_UPSTREAM_RECONNECT",
		"ST_UPSTREAM_DIAL",
		"ST_UPSTREAM_DIAL_ERR",
		"ST_UPSTREAM_UPDATE",
		"ST_UPSTREAM_UPDATE_ITEM",
		"ST_UPSTREAM_UPDATE_ERR",
	}
)

var (
	statCnt [ST_ARRAY_SIZE]uint64
)

func statHandle() (ret string) {
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n", statName[i],
			atomic.LoadUint64(&statCnt[i]))
	}
	return ret
}

func statInc(idx, n int) {
	atomic.AddUint64(&statCnt[idx], uint64(n))
}

func statSet(idx, n int) {
	atomic.StoreUint64(&statCnt[idx], uint64(n))
}

func statGet(idx int) uint64 {
	return atomic.LoadUint64(&statCnt[idx])
}

func (p *Ctrl) statStart() {
	if p.Conf.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statHandle())
				case _, ok := <-p.running:
					if !ok {
						return
					}
				}
			}
		}()
	}
}

func (p *Ctrl) statStop() {
}
