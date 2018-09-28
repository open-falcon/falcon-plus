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
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/sys"
)

var timeout = 30

func DuMetrics() (L []*model.MetricValue) {
	paths := g.DuPaths()
	result := make(chan *model.MetricValue, len(paths))
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			var err error
			defer func() {
				if err != nil {
					log.Println(err)
					result <- GaugeValue(g.DU_BS, -1, "path="+path)
				}
				wg.Done()
			}()
			//tips:osx  does not support -b.
			cmd := exec.Command("du", "-bs", path)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Start()
			if err != nil {
				return

			}
			err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Second)
			if isTimeout {
				err = errors.New(fmt.Sprintf("exec cmd : du -bs %s timeout", path))
				return
			}

			errStr := stderr.String()
			if errStr != "" {
				err = errors.New(errStr)
				return
			}

			if err != nil {
				err = errors.New(fmt.Sprintf("du -bs %s failed: %s", path, err.Error()))
				return
			}

			arr := strings.Fields(stdout.String())
			if len(arr) < 2 {
				err = errors.New(fmt.Sprintf("du -bs %s failed: %s", path, "return fields < 2"))
				return
			}

			size, err := strconv.ParseUint(arr[0], 10, 64)
			if err != nil {
				err = errors.New(fmt.Sprintf("cannot parse du -bs %s output", path))
				return
			}
			result <- GaugeValue(g.DU_BS, size, "path="+path)
		}(path)
	}
	wg.Wait()

	resultLen := len(result)
	for i := 0; i < resultLen; i++ {
		L = append(L, <-result)
	}
	return
}
