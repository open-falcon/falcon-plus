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
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
	log "github.com/Sirupsen/logrus"
)

func PortMetrics() (L []*model.MetricValue) {

	reportPorts := g.ReportPorts()
	sz := len(reportPorts)
	if sz == 0 {
		return
	}

	allTcpPorts, err := nux.TcpPorts()
	if err != nil {
		log.Println(err)
		return
	}

	allUdpPorts, err := nux.UdpPorts()
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", reportPorts[i])
		if slice.ContainsInt64(allTcpPorts, reportPorts[i]) || slice.ContainsInt64(allUdpPorts, reportPorts[i]) {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 1, tags))
		} else {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 0, tags))
		}
	}

	return
}
