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
package plugins

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl"
)

func init() {
	ctrl.RegisterPrestart(start)
}

func start(conf *falcon.ConfCtrl) error {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-ticker.C:
				beego.Debug("demo")
			}
		}
	}()
	return nil
}
