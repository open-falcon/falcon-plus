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
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
	"log"
)

func UdpMetrics() []*model.MetricValue {
	udp, err := nux.Snmp("Udp")
	if err != nil {
		log.Println("read snmp fail", err)
		return []*model.MetricValue{}
	}

	count := len(udp)
	ret := make([]*model.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.Udp."+key, val)
		i++
	}

	return ret
}
