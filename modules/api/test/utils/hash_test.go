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

package test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestHash(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	Convey("Test Hash method", t, func() {
		val := utils.HashIt("test2")
		So(val, ShouldEqual, "c0fc7c3e09f7efc71567b453ec5b9cd2")
	})
}
