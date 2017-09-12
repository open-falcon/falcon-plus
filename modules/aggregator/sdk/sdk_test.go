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

package sdk

import (
	"testing"

	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	g.ParseConfig("../cfg.example.json")
}

func TestSDK(t *testing.T) {
	Convey("get hostnames by id", t, func() {
		r, err := HostnamesByID(1)
		t.Log(r, err)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeGreaterThanOrEqualTo, 0)
	})

	Convey("query last points", t, func() {
		r, err := QueryLastPoints([]string{"laiweiofficemac"}, []string{"agent.alive"})
		t.Log(r, err)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeGreaterThanOrEqualTo, 0)
		for _, x := range r {
			t.Log(x)
		}
	})
}
