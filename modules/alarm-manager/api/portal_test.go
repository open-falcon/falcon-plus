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

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
)

func init() {
	config.InitApi(viper.GetViper())
}

func TestPortalAPI(t *testing.T) {
	Convey("Get action from api failed", t, func() {
		r := CurlAction(1)
		So(r.ID, ShouldEqual, 1)
	})
}
