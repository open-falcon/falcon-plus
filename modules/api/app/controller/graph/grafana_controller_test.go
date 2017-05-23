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

package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGrafanaController(t *testing.T) {
	Convey("test cutEndpointCounterHelp", t, func() {
		hosts, _ := cutEndpointCounterHelp("1.1.1.1")
		So(len(hosts), ShouldEqual, 0)
		hosts, _ = cutEndpointCounterHelp("{1.1.1.1}")
		So(len(hosts), ShouldEqual, 0)
		hosts, counter := cutEndpointCounterHelp("{1.1.1.1}#cpu.+")
		So(len(hosts), ShouldEqual, 1)
		So(counter, ShouldEqual, "cpu.+")
		hosts, counter = cutEndpointCounterHelp("{1.1.1.1,2.2.2.2,3.3.3.3}#cpu#idle.+")
		So(hosts[0], ShouldEqual, "1.1.1.1")
		So(hosts[2], ShouldEqual, "3.3.3.3")
		So(counter, ShouldEqual, "cpu\\.idle.+")
		hosts, counter = cutEndpointCounterHelp("1.1.1.1#net#if#bin.+")
		So(hosts[0], ShouldEqual, "1.1.1.1")
		So(counter, ShouldEqual, "net\\.if\\.bin.+")
	})

	Convey("test expandableChecking", t, func() {
		expsub, needexp := expandableChecking("cpu.idle", "cpu.+")
		So(expsub, ShouldEqual, "idle")
		So(needexp, ShouldEqual, false)
		expsub, needexp = expandableChecking("cpu.idle", "cpu")
		So(expsub, ShouldEqual, "idle")
		So(needexp, ShouldEqual, false)
		expsub, needexp = expandableChecking("net.if.out.bits/iface=eth_all", "net\\.if.+")
		So(expsub, ShouldEqual, "out")
		So(needexp, ShouldEqual, true)
		expsub, needexp = expandableChecking("net.if.out.bits/iface=eth_all", "net\\.if\\.out")
		So(expsub, ShouldEqual, "bits/iface=eth_all")
		So(needexp, ShouldEqual, false)
	})

}
