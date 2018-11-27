// Copyright 2018 Xiaomi, Inc.
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

package api

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	config.InitApi(viper.GetViper())
}

func TestUicAPI(t *testing.T) {
	Convey("Get team users from api failed", t, func() {
		r := CurlUic("plus-dev")
		for _, x := range r {
			log.Debugf("%#v", x)
		}
		So(len(r), ShouldEqual, 1)
	})
}
